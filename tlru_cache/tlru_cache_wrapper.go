package tlru_cache

import (
	"time"

	"github.com/preston-wagner/go-cache/cache_interface"
	"github.com/preston-wagner/go-cache/cache_wrapper"
)

// each of these wrappers returns both the wrapped getter function, as well as the StopReaping function of the underlying cache

func WrapWithTLRUCache[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.Getter[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
	memoryLimit float64,
) (cache_interface.Getter[KEY_TYPE, VALUE_TYPE], func()) {
	tlruCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency, memoryLimit)
	return cache_wrapper.WrapWithCache[KEY_TYPE, VALUE_TYPE](getter, tlruCache), tlruCache.StopReaping
}

func WrapWithTLRUCacheWithUncachedError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
	memoryLimit float64,
) (cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE], func()) {
	tlruCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency, memoryLimit)
	return cache_wrapper.WrapWithCacheWithUncachedError[KEY_TYPE, VALUE_TYPE](getter, tlruCache), tlruCache.StopReaping
}

func WrapWithTLRUCacheWithError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
	memoryLimit float64,
) (cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE], func()) {
	valueCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency, memoryLimit)
	errCache := NewTLRUCache[KEY_TYPE, error](ttl, reapFrequency, memoryLimit)
	stopReaping := func() {
		valueCache.StopReaping()
		errCache.StopReaping()
	}
	return cache_wrapper.WrapWithCacheWithError(getter, valueCache, errCache), stopReaping
}
