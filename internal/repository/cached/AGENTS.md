# Cached Repository: Agent Guidance

This package provides a thread-safe caching wrapper for the `domain.ListingRepository` using an in-memory TTL strategy.

## Implementation Details

- **Concurrency**: Uses `sync.RWMutex` to allow multiple concurrent readers but exclusive writer access during cache misses/refreshes.
- **Copy-on-Return**: To prevent mutations of the cached state by callers, always return a **deep copy** of map and slice results.
- **TTL Strategy**: Tracks `countsTime` and `locationsTime` separately for granular invalidation.

## Adding New Cached Methods

1.  Follow the pattern of `GetCounts`: 
    - Attempt a `RLock` read first.
    - If miss, release and fetch from the underlying store.
    - Use a full `Lock` for the update.
    - Return a copy.

## Common Pitfalls

- **State Leakage**: Returning the raw slice/map directly from the struct allows external packages to modify the cache content. **NEVER** return pointers to internal slices/maps.
- **Lock Contention**: Keep the write-locked section as small as possible. Fetch the data from the repository *before* acquiring the write lock.
