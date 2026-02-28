package cached

import (
	"context"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// CachedListingStore wraps a ListingStore and caches GetCounts results with a TTL.
type CachedListingStore struct {
	domain.ListingStore
	mu         sync.RWMutex
	counts     map[domain.Category]int
	countsTime time.Time
	ttl        time.Duration
}

// NewCachedListingStore creates a new cached wrapper around a ListingStore.
func NewCachedListingStore(store domain.ListingStore, ttl time.Duration) *CachedListingStore {
	return &CachedListingStore{
		ListingStore: store,
		ttl:          ttl,
	}
}

// GetCounts returns cached category counts, refreshing from the underlying store
// if the cache is expired or empty.
func (c *CachedListingStore) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	c.mu.RLock()
	if c.counts != nil && time.Since(c.countsTime) < c.ttl {
		// Cache hit — return a copy to prevent mutation
		result := make(map[domain.Category]int, len(c.counts))
		for k, v := range c.counts {
			result[k] = v
		}
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	// Cache miss — fetch from underlying store
	counts, err := c.ListingStore.GetCounts(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.mu.Lock()
	c.counts = counts
	c.countsTime = time.Now()
	c.mu.Unlock()

	// Return a copy
	result := make(map[domain.Category]int, len(counts))
	for k, v := range counts {
		result[k] = v
	}
	return result, nil
}
