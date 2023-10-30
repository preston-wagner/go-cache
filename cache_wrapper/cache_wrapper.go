package cache_wrapper

import (
	"github.com/preston-wagner/go-cache/cache_interface"
	"github.com/preston-wagner/unicycle/defaults"
)

func WrapWithCache[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.Getter[KEY_TYPE, VALUE_TYPE],
	cache cache_interface.Cache[KEY_TYPE, VALUE_TYPE],
) cache_interface.Getter[KEY_TYPE, VALUE_TYPE] {
	return func(key KEY_TYPE) VALUE_TYPE {
		value, ok := cache.Get(key)
		if ok {
			return value
		}
		value = getter(key)
		cache.Set(key, value)
		return value
	}
}

// only caches successful calls, not errors
func WrapWithCacheWithUncachedError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	cache cache_interface.Cache[KEY_TYPE, VALUE_TYPE],
) cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE] {
	return func(key KEY_TYPE) (VALUE_TYPE, error) {
		value, ok := cache.Get(key)
		if ok {
			return value, nil
		}
		value, err := getter(key)
		if err == nil {
			cache.Set(key, value)
		}
		return value, err
	}
}

// caches both successful calls and errors
func WrapWithCacheWithError[KEY_TYPE comparable, VALUE_TYPE any](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	cache cache_interface.Cache[KEY_TYPE, VALUE_TYPE],
	errCache cache_interface.Cache[KEY_TYPE, error],
) cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE] {
	return func(key KEY_TYPE) (VALUE_TYPE, error) {
		value, ok := cache.Get(key)
		if ok {
			return value, nil
		}
		err, ok := errCache.Get(key)
		if ok {
			return defaults.ZeroValue[VALUE_TYPE](), err
		}
		value, err = getter(key)
		if err == nil {
			cache.Set(key, value)
		} else {
			errCache.Set(key, err)
		}
		return value, err
	}
}
