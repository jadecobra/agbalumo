package cached

import (
	"context"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// CachedListingStore wraps a ListingRepository and caches GetCounts results with a TTL.
type CachedListingStore struct {
	countsTime    time.Time
	locationsTime time.Time
	domain.ListingRepository
	counts    map[domain.Category]int
	locations []domain.Location
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewCachedListingStore creates a new cached wrapper around a ListingRepository.
func NewCachedListingStore(store domain.ListingRepository, ttl time.Duration) *CachedListingStore {
	return &CachedListingStore{
		ListingRepository: store,
		ttl:               ttl,
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
	counts, err := c.ListingRepository.GetCounts(ctx)
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

// GetLocations returns cached locations, refreshing from the underlying store
// if the cache is expired or empty.
func (c *CachedListingStore) GetLocations(ctx context.Context) ([]domain.Location, error) {
	c.mu.RLock()
	if c.locations != nil && time.Since(c.locationsTime) < c.ttl {
		// Cache hit — return a copy to prevent mutation
		result := make([]domain.Location, len(c.locations))
		copy(result, c.locations)
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	// Cache miss — fetch from underlying store
	locations, err := c.ListingRepository.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.mu.Lock()
	c.locations = locations
	c.locationsTime = time.Now()
	c.mu.Unlock()

	// Return a copy
	result := make([]domain.Location, len(locations))
	copy(result, locations)
	return result, nil
}
