package ttl_cache_test

import (
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/cache_test"
	"github.com/preston-wagner/go-cache/ttl_cache"
)

func TestWrapWithTTLCache(t *testing.T) {
	wrappedGetter, stop := ttl_cache.WrapWithTTLCache[int, string](cache_test.ItoaOnce(t), time.Second, time.Second*3)
	res := wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTTLCache wrapped getter returned wrong value")
	}
	res = wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTTLCache wrapped getter returned wrong cached value")
	}
	stop()
}

func TestWrapWithTTLCacheWithUncachedError(t *testing.T) {
	wrappedGetter, stop := ttl_cache.WrapWithTTLCacheWithUncachedError[string, int](cache_test.AtoiOnce(t), time.Second, time.Second*3)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTTLCacheWithUncachedError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTTLCacheWithUncachedError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTTLCacheWithUncachedError should return an error when the getter does")
	}
	stop()
}

func TestWrapWithTTLCacheWithError(t *testing.T) {
	wrappedGetter, stop := ttl_cache.WrapWithTTLCacheWithError[string, int](cache_test.AtoiOnce(t), time.Second, time.Second*3)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTTLCacheWithError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTTLCacheWithError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTTLCacheWithError should return an error when the getter does")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTTLCacheWithError wrapped getter returned wrong cached error")
	}
	stop()
}
