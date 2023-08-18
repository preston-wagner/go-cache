package cache

import (
	"sync"
	"time"

	"github.com/nuvi/unicycle"
)

type tlruCacheValue[VALUE_TYPE any] struct {
	value        VALUE_TYPE
	lastAccessed time.Time
	addedAt      time.Time
}

// a Time-aware Least Recently Used cache: when the memory threshold specified is exceeded, the cache will start deleting the least-recently accessed values first, as well as any expired values
type TLRUCache[KEY_TYPE comparable, VALUE_TYPE any] struct {
	cache       map[KEY_TYPE]tlruCacheValue[VALUE_TYPE]
	ttl         time.Duration
	memoryLimit float64 // a value between 0 and 1, representing the threshold at which the cache should start purging old values (0.9 would mean "start deleting when 90% of available ram has been used")
	lock        *sync.RWMutex
	done        chan bool
}

func NewTLRUCache[KEY_TYPE comparable, VALUE_TYPE any](ttl time.Duration, memoryLimit float64) *TLRUCache[KEY_TYPE, VALUE_TYPE] {
	tlruCache := &TLRUCache[KEY_TYPE, VALUE_TYPE]{
		cache:       make(map[KEY_TYPE]tlruCacheValue[VALUE_TYPE]),
		ttl:         ttl,
		memoryLimit: memoryLimit,
		lock:        &sync.RWMutex{},
		done:        make(chan bool),
	}
	tlruCache.StartReaping()
	return tlruCache
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) Set(key KEY_TYPE, value VALUE_TYPE) {
	tlruCache.lock.Lock()
	defer tlruCache.lock.Unlock()
	tlruCache.cache[key] = tlruCacheValue[VALUE_TYPE]{
		value:        value,
		lastAccessed: time.Now(),
		addedAt:      time.Now(),
	}
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) Get(key KEY_TYPE) (VALUE_TYPE, bool) {
	tlruCache.lock.Lock()
	defer tlruCache.lock.Unlock()
	value, ok := tlruCache.cache[key]
	if ok {
		if value.addedAt.Add(tlruCache.ttl).Before(time.Now()) {
			delete(tlruCache.cache, key)
			return value.value, false
		} else {
			value.addedAt = time.Now()
			tlruCache.cache[key] = value
			return value.value, ok
		}
	} else {
		return unicycle.ZeroValue[VALUE_TYPE](), ok
	}
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) Remove(key KEY_TYPE) {
	tlruCache.lock.Lock()
	defer tlruCache.lock.Unlock()
	delete(tlruCache.cache, key)
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) Len() int {
	tlruCache.lock.RLock()
	defer tlruCache.lock.RUnlock()
	return len(tlruCache.cache)
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) StartReaping() {
	go func() {
		ticker := time.NewTicker(time.Second)
		select {
		case <-tlruCache.done:
			return
		case <-ticker.C:
			tlruCache.Reap()
		}
	}()
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) StopReaping() {
	go func() {
		tlruCache.done <- true
	}()
}

func (tlruCache *TLRUCache[KEY_TYPE, VALUE_TYPE]) Reap() {
	if tlruCache.Len() > 0 {
		tlruCache.lock.Lock()
		defer tlruCache.lock.Unlock()

		now := time.Now()
		isNotExpired := func(key KEY_TYPE, value tlruCacheValue[VALUE_TYPE]) bool {
			return now.Before(value.addedAt.Add(tlruCache.ttl))
		}

		filter := isNotExpired

		if MemUsage() >= tlruCache.memoryLimit {
			var total, count int64
			for _, value := range tlruCache.cache {
				total += value.lastAccessed.UnixMicro()
				count++
			}
			avg := time.UnixMicro(total / count)

			filter = func(key KEY_TYPE, value tlruCacheValue[VALUE_TYPE]) bool {
				return isNotExpired(key, value) && avg.Before(value.lastAccessed)
			}
		}

		tlruCache.cache = unicycle.FilterMap(tlruCache.cache, filter)
	}
}
