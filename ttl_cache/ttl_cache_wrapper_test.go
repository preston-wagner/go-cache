package ttl_cache_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/ttl_cache"
)

func TestWrapWithTTLCache(t *testing.T) {
	wrappedGetter := ttl_cache.WrapWithTTLCache[int, string](strconv.Itoa, time.Second, time.Second*3)
	res := wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTTLCache wrapped getter returned wrong value")
	}
	res = wrappedGetter(7)
	if res != "7" {
		t.Error("WrapWithTTLCache wrapped getter returned wrong cached value")
	}
}

func TestWrapWithTTLCacheWithUncachedError(t *testing.T) {
	wrappedGetter := ttl_cache.WrapWithTTLCacheWithUncachedError[string, int](strconv.Atoi, time.Second, time.Second*3)
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
}
