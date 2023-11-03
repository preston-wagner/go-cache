# Time-aware, Least-Recently Used (TLRU) Cache
https://en.wikipedia.org/wiki/Cache_replacement_policies#Time-aware,_least-recently_used
Values are removed from the cache when:
- the value has been kept for longer than the time-to-live allows
- the cache detects that the provided memory usage threshold has been exceeded, and the value is older than the average of all current values
