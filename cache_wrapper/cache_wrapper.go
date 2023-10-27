package cache_wrapper

import "github.com/preston-wagner/go-cache/cache_interface"

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
