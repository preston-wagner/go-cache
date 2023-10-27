package ttl_cache

import (
	"time"

	"github.com/preston-wagner/go-cache/cache_interface"
	"github.com/preston-wagner/go-cache/cache_wrapper"
)

func WrapWithTTLCache[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.Getter[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
) cache_interface.Getter[KEY_TYPE, VALUE_TYPE] {
	ttlCache := NewTTLCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency)
	return cache_wrapper.WrapWithCache[KEY_TYPE, VALUE_TYPE](getter, ttlCache)
}

func WrapWithTTLCacheWithUncachedError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	ttl time.Duration,
	reapFrequency time.Duration,
) cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE] {
	ttlCache := NewTTLCache[KEY_TYPE, VALUE_TYPE](ttl, reapFrequency)
	return cache_wrapper.WrapWithCacheWithUncachedError[KEY_TYPE, VALUE_TYPE](getter, ttlCache)
}
