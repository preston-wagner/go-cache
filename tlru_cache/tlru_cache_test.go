package tlru_cache_test

import (
	"testing"
	"time"

	"github.com/preston-wagner/go-cache/tlru_cache"
)

func TestTLRUCache(t *testing.T) {
	// test purging based on TTL expiration
	tlruCache := tlru_cache.NewTLRUCache[int, string](time.Second, time.Second, 0.9)

	tlruCache.Set(1, "1")
	if value, ok := tlruCache.Get(1); value != "1" || !ok {
		t.Error("tlruCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := tlruCache.Get(1); ok {
		t.Error("tlruCache.Get should not return ok for expired values")
	}

	tlruCache.StopReaping()

	// test purging based on memory usage
	tlruCache = tlru_cache.NewTLRUCache[int, string](time.Second*5, time.Second, 0)

	tlruCache.Set(1, "1")
	tlruCache.Set(2, "2")
	tlruCache.Set(3, "3")
	if value, ok := tlruCache.Get(1); value != "1" || !ok {
		t.Error("tlruCache.Get failed")
	}

	time.Sleep(time.Second * 2)

	if _, ok := tlruCache.Get(1); ok {
		t.Error("tlruCache.Get should not return ok for values purged by reaper")
	}
	// the 2nd value may or may not be reaped based on the timing of the test, but the first should always be reaped and the last should never be
	if _, ok := tlruCache.Get(3); !ok {
		t.Error("tlruCache.Get should return ok for values more recent than the average after reaping")
	}

	tlruCache.StopReaping()
}
