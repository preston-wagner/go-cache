package cache

type Cache[KEY_TYPE comparable, VALUE_TYPE any] interface {
	Get(key KEY_TYPE) (VALUE_TYPE, bool)
	Set(key KEY_TYPE, value VALUE_TYPE)
	Remove(key KEY_TYPE)
}
