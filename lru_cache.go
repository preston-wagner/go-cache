package cache

import (
	"sync"
	"time"

	"github.com/preston-wagner/unicycle"
)

type cacheValue[VALUE_TYPE any] struct {
	value        VALUE_TYPE
	lastAccessed time.Time
	addedAt      time.Time
}

type LRUCache[KEY_TYPE comparable, VALUE_TYPE any] struct {
	cache       map[KEY_TYPE]cacheValue[VALUE_TYPE]
	ttl         time.Duration
	memoryLimit float64 // a value between 0 and 1, representing the threshold at which the cache should start purging old values (0.9 would mean "start deleting when 90% of available ram has been used")
	lock        *sync.RWMutex
	done        chan bool
}

func NewLRUCache[KEY_TYPE comparable, VALUE_TYPE any](ttl time.Duration, memoryLimit float64) *LRUCache[KEY_TYPE, VALUE_TYPE] {
	lruCache := &LRUCache[KEY_TYPE, VALUE_TYPE]{
		cache:       make(map[KEY_TYPE]cacheValue[VALUE_TYPE]),
		ttl:         ttl,
		memoryLimit: memoryLimit,
		lock:        &sync.RWMutex{},
		done:        make(chan bool),
	}
	lruCache.StartReaping()
	return lruCache
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) Set(key KEY_TYPE, value VALUE_TYPE) {
	lruCache.lock.Lock()
	defer lruCache.lock.Unlock()
	lruCache.cache[key] = cacheValue[VALUE_TYPE]{
		value:        value,
		lastAccessed: time.Now(),
		addedAt:      time.Now(),
	}
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) Get(key KEY_TYPE) (VALUE_TYPE, bool) {
	lruCache.lock.Lock()
	defer lruCache.lock.Unlock()
	value, ok := lruCache.cache[key]
	if ok {
		if value.addedAt.Add(lruCache.ttl).Before(time.Now()) {
			delete(lruCache.cache, key)
			return value.value, false
		} else {
			value.addedAt = time.Now()
			lruCache.cache[key] = value
			return value.value, ok
		}
	} else {
		return unicycle.ZeroValue[VALUE_TYPE](), ok
	}
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) Remove(key KEY_TYPE) {
	lruCache.lock.Lock()
	defer lruCache.lock.Unlock()
	delete(lruCache.cache, key)
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) Len() int {
	lruCache.lock.RLock()
	defer lruCache.lock.RUnlock()
	return len(lruCache.cache)
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) StartReaping() {
	go func() {
		ticker := time.NewTicker(time.Second)
		select {
		case <-lruCache.done:
			return
		case <-ticker.C:
			lruCache.Reap()
		}
	}()
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) StopReaping() {
	go func() {
		lruCache.done <- true
	}()
}

func (lruCache *LRUCache[KEY_TYPE, VALUE_TYPE]) Reap() {
	if lruCache.Len() > 0 {
		lruCache.lock.Lock()
		defer lruCache.lock.Unlock()

		now := time.Now()
		isNotExpired := func(key KEY_TYPE, value cacheValue[VALUE_TYPE]) bool {
			return now.Before(value.addedAt.Add(lruCache.ttl))
		}

		filter := isNotExpired

		if MemUsage() >= lruCache.memoryLimit {
			var total, count int64
			for _, value := range lruCache.cache {
				total += value.lastAccessed.UnixMicro()
				count++
			}
			avg := time.UnixMicro(total / count)

			filter = func(key KEY_TYPE, value cacheValue[VALUE_TYPE]) bool {
				return isNotExpired(key, value) && avg.Before(value.lastAccessed)
			}
		}

		lruCache.cache = unicycle.FilterMap(lruCache.cache, filter)
	}
}
