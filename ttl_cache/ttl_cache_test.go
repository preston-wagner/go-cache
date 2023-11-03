package ttl_cache_test

import (
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/ttl_cache"
)

func TestTTLCache(t *testing.T) {
	// test purging based on TTL expiration
	ttlCache := ttl_cache.NewTTLCache[int, string](time.Second, time.Second)

	ttlCache.Set(1, "1")
	if value, ok := ttlCache.Get(1); value != "1" || !ok {
		t.Error("TTLCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := ttlCache.Get(1); ok {
		t.Error("TTLCache.Get should not return ok for expired values")
	}

	ttlCache.Set(2, "2")
	ttlCache.Remove(2)
	_, ok := ttlCache.Get(2)
	if ok {
		t.Error("Get should return false for removed value")
	}

	ttlCache.StopReaping()

	ttlCache.Set(3, "3")

	time.Sleep(time.Second * 2)
	_, ok = ttlCache.Get(3)
	if ok {
		t.Error("Get should return false for expired value, even when reaping has ceased")
	}
}
