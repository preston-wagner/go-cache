package cache

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	// test purging based on cache expiration
	lruCache := NewLRUCache[int, string](time.Second, 1)

	lruCache.Set(1, "1")
	if value, ok := lruCache.Get(1); value != "1" || !ok {
		t.Error("lruCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := lruCache.Get(1); ok {
		t.Error("lruCache.Get should not return ok for expired values")
	}

	lruCache.StopReaping()

	// test purging based on memory usage
	lruCache = NewLRUCache[int, string](time.Second*5, 0)

	lruCache.Set(1, "1")
	lruCache.Set(2, "2")
	lruCache.Set(3, "3")
	if value, ok := lruCache.Get(1); value != "1" || !ok {
		t.Error("lruCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := lruCache.Get(1); ok {
		t.Error("lruCache.Get should not return ok for values purged by reaper")
	}
	// the 2nd value may or may not be reaped based on the timing of the test, but the first should always be reaped and the last should never be
	if _, ok := lruCache.Get(3); !ok {
		t.Error("lruCache.Get should return ok for values more recent than the average after reaping")
	}

	lruCache.StopReaping()
}
