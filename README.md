# go-cache
A library of generic, unit-tested caching algorithms.

## TTLCache
Possibly the simplest, most common form of caching, `TTLCache` stores a mapping of keys (of any comparable type) to values (of any type) with an expiration time.
A background goroutine is responsible for removing expired entries.
Note: The underlying map is never re-allocated, so while expired entries are deleted and the memory can be re-used, it is never freed for use by anything other than the cache itself.
