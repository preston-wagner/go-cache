package ttl_cache

import (
	"time"

	"github.com/preston-wagner/go-cache/cache_interface"
	"github.com/preston-wagner/go-cache/cache_wrapper"
)

// each of these wrappers returns both the wrapped getter function, as well as the StopReaping function of the underlying cache

func WrapWithTTLCache[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.Getter[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
) (cache_interface.Getter[KEY_TYPE, VALUE_TYPE], func()) {
	ttlCache := NewTTLCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency)
	return cache_wrapper.WrapWithCache[KEY_TYPE, VALUE_TYPE](getter, ttlCache), ttlCache.StopReaping
}

func WrapWithTTLCacheWithUncachedError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
) (cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE], func()) {
	ttlCache := NewTTLCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency)
	return cache_wrapper.WrapWithCacheWithUncachedError[KEY_TYPE, VALUE_TYPE](getter, ttlCache), ttlCache.StopReaping
}

func WrapWithTTLCacheWithError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
) (cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE], func()) {
	valueCache := NewTTLCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency)
	errCache := NewTTLCache[KEY_TYPE, error](ttl, reapFrequency)
	stopReaping := func() {
		valueCache.StopReaping()
		errCache.StopReaping()
	}
	return cache_wrapper.WrapWithCacheWithError(getter, valueCache, errCache), stopReaping
}
