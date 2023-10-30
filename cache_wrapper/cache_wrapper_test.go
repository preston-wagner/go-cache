package cache_wrapper_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/cache_wrapper"
	"github.com/preston-wagner/go-cache/ttl_cache"
	"github.com/preston-wagner/unicycle/sets"
)

func itoaOnce(t *testing.T) func(int) string {
	keys := sets.Set[int]{}
	return func(key int) string {
		if keys.Has(key) {
			t.Error("wrapped function was called with the same key multiple times!")
		}
		keys.Add(key)
		return strconv.Itoa(key)
	}
}

func atoiOnce(t *testing.T) func(string) (int, error) {
	keys := sets.Set[string]{}
	return func(key string) (int, error) {
		if keys.Has(key) {
			t.Error("wrapped function was called with the same key multiple times!")
		}
		keys.Add(key)
		return strconv.Atoi(key)
	}
}

func TestWrapWithCache(t *testing.T) {
	wrappedGetter := cache_wrapper.WrapWithCache[int, string](
		itoaOnce(t),
		ttl_cache.NewTTLCache[int, string](time.Second, time.Second*3),
	)
	res := wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithCache wrapped getter returned wrong value")
	}
	res = wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithCache wrapped getter returned wrong cached value")
	}
}

func TestWrapWithCacheWithUncachedError(t *testing.T) {
	wrappedGetter := cache_wrapper.WrapWithCacheWithUncachedError[string, int](
		atoiOnce(t),
		ttl_cache.NewTTLCache[string, int](time.Second, time.Second*3),
	)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithCacheWithUncachedError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithCacheWithUncachedError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithCacheWithUncachedError should return an error when the getter does")
	}
}

func TestWrapWithCacheWithError(t *testing.T) {
	wrappedGetter := cache_wrapper.WrapWithCacheWithError[string, int](
		atoiOnce(t),
		ttl_cache.NewTTLCache[string, int](time.Second, time.Second*3),
		ttl_cache.NewTTLCache[string, error](time.Second, time.Second*3),
	)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithCacheWithError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithCacheWithError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithCacheWithError should return an error when the getter does")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithCacheWithError wrapped getter returned wrong cached error")
	}
}
