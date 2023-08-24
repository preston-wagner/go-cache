package cache

import (
	"sync"
	"time"
)

type ttlCacheValue[VALUE_TYPE any] struct {
	createdAt time.Time
	value     VALUE_TYPE
}

// a cache that stores values for a fixed amount of time
type TTLCache[KEY_TYPE comparable, VALUE_TYPE any] struct {
	cache         map[KEY_TYPE]ttlCacheValue[VALUE_TYPE]
	ttl           time.Duration
	reapFrequency time.Duration
	lock          *sync.RWMutex
	done          chan bool
}

func NewTTLCache[KEY_TYPE comparable, VALUE_TYPE any](ttl, reapFrequency time.Duration) *TTLCache[KEY_TYPE, VALUE_TYPE] {
	cache := &TTLCache[KEY_TYPE, VALUE_TYPE]{
		cache:         map[KEY_TYPE]ttlCacheValue[VALUE_TYPE]{},
		ttl:           ttl,
		reapFrequency: reapFrequency,
		lock:          &sync.RWMutex{},
		done:          make(chan bool),
	}

	cache.StartReaping()

	return cache
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) Set(key KEY_TYPE, value VALUE_TYPE) {
	ttlCache.lock.Lock()
	defer ttlCache.lock.Unlock()
	ttlCache.cache[key] = ttlCacheValue[VALUE_TYPE]{
		value:     value,
		createdAt: time.Now(),
	}
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) Get(key KEY_TYPE) (VALUE_TYPE, bool) {
	ttlCache.lock.RLock()
	defer ttlCache.lock.RUnlock()
	element, ok := ttlCache.cache[key]
	return element.value, ok
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) Remove(key KEY_TYPE) {
	ttlCache.lock.Lock()
	defer ttlCache.lock.Unlock()
	delete(ttlCache.cache, key)
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) StartReaping() {
	go func() {
		ticker := time.NewTicker(ttlCache.reapFrequency)
		select {
		case <-ticker.C:
			ttlCache.Reap()
		case <-ttlCache.done:
			ticker.Stop()
			return
		}
	}()
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) StopReaping() {
	go func() {
		ttlCache.done <- true
	}()
}

func (ttlCache *TTLCache[KEY_TYPE, VALUE_TYPE]) Reap() {
	ttlCache.lock.Lock()
	defer ttlCache.lock.Unlock()
	expired := time.Now().Add(-ttlCache.ttl)
	for key, element := range ttlCache.cache {
		if expired.After(element.createdAt) {
			// Note: deleting from a map doesn't shrink the underlying array, so while memory can be re-used it is not freed
			delete(ttlCache.cache, key)
		}
	}
}
