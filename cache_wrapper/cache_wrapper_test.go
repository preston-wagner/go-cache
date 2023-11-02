package cache_wrapper_test

import (
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/cache_test"
	"github.com/preston-wagner/go-cache/cache_wrapper"
	"github.com/preston-wagner/go-cache/ttl_cache"
)

func TestWrapWithCache(t *testing.T) {
	wrappedGetter := cache_wrapper.WrapWithCache[int, string](
		cache_test.ItoaOnce(t),
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
		cache_test.AtoiOnce(t),
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
		cache_test.AtoiOnce(t),
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
