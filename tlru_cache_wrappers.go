package cache

import (
	"time"

	"github.com/nuvi/unicycle"
)

func WrapWithTLRUCache[KEY_TYPE comparable, VALUE_TYPE any](getter Getter[KEY_TYPE, VALUE_TYPE], ttl time.Duration, memoryLimit float64) Getter[KEY_TYPE, VALUE_TYPE] {
	tlruCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, memoryLimit)
	return WrapWithCache[KEY_TYPE, VALUE_TYPE](getter, tlruCache)
}

func WrapWithTLRUCacheWithUncachedError[KEY_TYPE comparable, VALUE_TYPE any](getter GetterWithError[KEY_TYPE, VALUE_TYPE], ttl time.Duration, memoryLimit float64) GetterWithError[KEY_TYPE, VALUE_TYPE] {
	tlruCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, memoryLimit)
	return WrapWithCacheWithUncachedError[KEY_TYPE, VALUE_TYPE](getter, tlruCache)
}

func WrapWithTLRUCacheWithCachedError[KEY_TYPE comparable, VALUE_TYPE any](getter GetterWithError[KEY_TYPE, VALUE_TYPE], ttl time.Duration, memoryLimit float64) GetterWithError[KEY_TYPE, VALUE_TYPE] {
	valueCache := NewTLRUCache[KEY_TYPE, VALUE_TYPE](ttl, memoryLimit)
	errCache := NewTLRUCache[KEY_TYPE, error](ttl, memoryLimit)
	return func(key KEY_TYPE) (VALUE_TYPE, error) {
		value, ok := valueCache.Get(key)
		if ok {
			return value, nil
		}
		err, ok := errCache.Get(key)
		if ok {
			return unicycle.ZeroValue[VALUE_TYPE](), err
		}

		value, err = getter(key)
		if err != nil {
			errCache.Set(key, err)
		} else {
			valueCache.Set(key, value)
		}
		return value, err
	}
}
