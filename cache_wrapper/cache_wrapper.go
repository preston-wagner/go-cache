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
// note: the VALUE_CACHE_TYPE and ERROR_CACHE_TYPE look redundant (and are),
// but there's an issue with the go compiler where it erroneously claims that implementations
// of the Cache interface like TLRUCache don't match unless it's written like this,
// possibly because there's some hard limit to the number of chained inferences
func WrapWithCacheWithError[
	KEY_TYPE comparable,
	VALUE_TYPE any,
	VALUE_CACHE_TYPE cache_interface.Cache[KEY_TYPE, VALUE_TYPE],
	ERROR_CACHE_TYPE cache_interface.Cache[KEY_TYPE, error],
](
	getter cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE],
	valueCache VALUE_CACHE_TYPE,
	errCache ERROR_CACHE_TYPE,
) cache_interface.GetterWithError[KEY_TYPE, VALUE_TYPE] {
	return func(key KEY_TYPE) (VALUE_TYPE, error) {
		value, ok := valueCache.Get(key)
		if ok {
			return value, nil
		}
		err, ok := errCache.Get(key)
		if ok {
			return defaults.ZeroValue[VALUE_TYPE](), err
		}
		value, err = getter(key)
		if err == nil {
			valueCache.Set(key, value)
		} else {
			errCache.Set(key, err)
		}
		return value, err
	}
}
