package tlru_cache_test

import (
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/cache_test"
	"github.com/preston-wagner/go-cache/tlru_cache"
)

func TestWrapWithTLRUCache(t *testing.T) {
	wrappedGetter, stop := tlru_cache.WrapWithTLRUCache[int, string](
		cache_test.ItoaOnce(t),
		time.Second,
		time.Second,
		0.9,
	)
	res := wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTLRUCache wrapped getter returned wrong value")
	}
	res = wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTLRUCache wrapped getter returned wrong cached value")
	}
	stop()
}

func TestWrapWithTLRUCacheWithUncachedError(t *testing.T) {
	wrappedGetter, stop := tlru_cache.WrapWithTLRUCacheWithUncachedError[string, int](
		cache_test.AtoiOnce(t),
		time.Second,
		time.Second,
		0.9,
	)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTLRUCacheWithUncachedError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTLRUCacheWithUncachedError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTLRUCacheWithUncachedError should return an error when the getter does")
	}
	stop()
}

func TestWrapWithTLRUCacheWithError(t *testing.T) {
	wrappedGetter, stop := tlru_cache.WrapWithTLRUCacheWithError[string, int](
		cache_test.AtoiOnce(t),
		time.Second,
		time.Second,
		0.9,
	)
	res, err := wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTLRUCacheWithError wrapped getter returned wrong value")
	}
	res, err = wrappedGetter("7")
	if err != nil {
		t.Error(err)
	}
	if res != 7 {
		t.Error("WrapWithTLRUCacheWithError wrapped getter returned wrong cached value")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTLRUCacheWithError should return an error when the getter does")
	}
	_, err = wrappedGetter("A")
	if err == nil {
		t.Error("WrapWithTLRUCacheWithError wrapped getter returned wrong cached error")
	}
	stop()
}
