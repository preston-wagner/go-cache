package cache_interface

type Getter[KEY_TYPE comparable, VALUE_TYPE any] func(KEY_TYPE) VALUE_TYPE

type GetterWithError[KEY_TYPE comparable, VALUE_TYPE any] func(KEY_TYPE) (VALUE_TYPE, error)
