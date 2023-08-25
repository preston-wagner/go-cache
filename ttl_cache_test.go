package cache

import (
	"testing"
	"time"
)

func TestTTLCache(t *testing.T) {
	// test purging based on TTL expiration
	ttlCache := NewTTLCache[int, string](time.Second, time.Second)

	ttlCache.Set(1, "1")
	if value, ok := ttlCache.Get(1); value != "1" || !ok {
		t.Error("TTLCache.Get failed")
	}

	time.Sleep(time.Second * 3)

	if _, ok := ttlCache.Get(1); ok {
		t.Error("TTLCache.Get should not return ok for expired values")
	}

	ttlCache.StopReaping()
}
