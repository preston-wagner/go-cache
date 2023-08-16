package cache

import (
	"testing"
	"time"
)

func TestTLRUCache(t *testing.T) {
	// test purging based on cache expiration
	TLRUCache := NewTLRUCache[int, string](time.Second, 1)

	TLRUCache.Set(1, "1")
	if value, ok := TLRUCache.Get(1); value != "1" || !ok {
		t.Error("TLRUCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := TLRUCache.Get(1); ok {
		t.Error("TLRUCache.Get should not return ok for expired values")
	}

	TLRUCache.StopReaping()

	// test purging based on memory usage
	TLRUCache = NewTLRUCache[int, string](time.Second*5, 0)

	TLRUCache.Set(1, "1")
	TLRUCache.Set(2, "2")
	TLRUCache.Set(3, "3")
	if value, ok := TLRUCache.Get(1); value != "1" || !ok {
		t.Error("TLRUCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := TLRUCache.Get(1); ok {
		t.Error("TLRUCache.Get should not return ok for values purged by reaper")
	}
	// the 2nd value may or may not be reaped based on the timing of the test, but the first should always be reaped and the last should never be
	if _, ok := TLRUCache.Get(3); !ok {
		t.Error("TLRUCache.Get should return ok for values more recent than the average after reaping")
	}

	TLRUCache.StopReaping()
}
