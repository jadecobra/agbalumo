package cached

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// stubListingStore is a minimal stub satisfying domain.ListingStore for testing.
type stubListingStore struct {
	domain.ListingRepository
	getCountsFunc     func() (map[domain.Category]int, error)
	getLocationsFunc  func() ([]string, error)
	getCountsCalls    int
	getLocationsCalls int
}

func (s *stubListingStore) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	s.getCountsCalls++
	if s.getCountsFunc != nil {
		return s.getCountsFunc()
	}
	return map[domain.Category]int{domain.Business: 5, domain.Food: 3}, nil
}

func (s *stubListingStore) GetLocations(ctx context.Context) ([]string, error) {
	s.getLocationsCalls++
	if s.getLocationsFunc != nil {
		return s.getLocationsFunc()
	}
	return []string{"Lagos", "London"}, nil
}

func TestCachedGetCounts_CacheMiss(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	counts, err := cache.GetCounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getCountsCalls != 1 {
		t.Errorf("expected 1 call to underlying store, got %d", stub.getCountsCalls)
	}
	if counts[domain.Business] != 5 {
		t.Errorf("expected Business=5, got %d", counts[domain.Business])
	}
	if counts[domain.Food] != 3 {
		t.Errorf("expected Food=3, got %d", counts[domain.Food])
	}
}

func TestCachedGetCounts_CacheHit(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	// First call — cache miss
	_, _ = cache.GetCounts(context.Background())
	// Second call — cache hit
	counts, err := cache.GetCounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getCountsCalls != 1 {
		t.Errorf("expected 1 call to underlying store (cached), got %d", stub.getCountsCalls)
	}
	if counts[domain.Business] != 5 {
		t.Errorf("expected Business=5, got %d", counts[domain.Business])
	}
}

func TestCachedGetCounts_CacheExpired(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 1*time.Millisecond)

	// First call — cache miss
	_, _ = cache.GetCounts(context.Background())
	// Wait for TTL to expire
	time.Sleep(5 * time.Millisecond)
	// Second call — cache expired, should fetch again
	_, err := cache.GetCounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getCountsCalls != 2 {
		t.Errorf("expected 2 calls to underlying store (expired), got %d", stub.getCountsCalls)
	}
}

func TestCachedGetCounts_ErrorPassthrough(t *testing.T) {
	expectedErr := errors.New("db connection lost")
	stub := &stubListingStore{
		getCountsFunc: func() (map[domain.Category]int, error) {
			return nil, expectedErr
		},
	}
	cache := NewCachedListingStore(stub, 60*time.Second)

	_, err := cache.GetCounts(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestCachedGetCounts_ReturnsCopy(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	counts1, _ := cache.GetCounts(context.Background())
	// Mutate the returned map
	counts1[domain.Business] = 999

	// Get again — should still be the original value
	counts2, _ := cache.GetCounts(context.Background())
	if counts2[domain.Business] != 5 {
		t.Errorf("cache was mutated by caller: expected Business=5, got %d", counts2[domain.Business])
	}
}

func TestCachedGetLocations_CacheMiss(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	locs, err := cache.GetLocations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getLocationsCalls != 1 {
		t.Errorf("expected 1 call to underlying store, got %d", stub.getLocationsCalls)
	}
	if len(locs) != 2 || locs[0] != "Lagos" || locs[1] != "London" {
		t.Errorf("unexpected locations: %v", locs)
	}
}

func TestCachedGetLocations_CacheHit(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	// First call — cache miss
	_, _ = cache.GetLocations(context.Background())
	// Second call — cache hit
	locs, err := cache.GetLocations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getLocationsCalls != 1 {
		t.Errorf("expected 1 call to underlying store (cached), got %d", stub.getLocationsCalls)
	}
	if len(locs) != 2 || locs[0] != "Lagos" || locs[1] != "London" {
		t.Errorf("unexpected locations: %v", locs)
	}
}

func TestCachedGetLocations_CacheExpired(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 1*time.Millisecond)

	// First call — cache miss
	_, _ = cache.GetLocations(context.Background())
	// Wait for TTL to expire
	time.Sleep(5 * time.Millisecond)
	// Second call — cache expired, should fetch again
	_, err := cache.GetLocations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.getLocationsCalls != 2 {
		t.Errorf("expected 2 calls to underlying store (expired), got %d", stub.getLocationsCalls)
	}
}

func TestCachedGetLocations_ErrorPassthrough(t *testing.T) {
	expectedErr := errors.New("db connection lost")
	stub := &stubListingStore{
		getLocationsFunc: func() ([]string, error) {
			return nil, expectedErr
		},
	}
	cache := NewCachedListingStore(stub, 60*time.Second)

	_, err := cache.GetLocations(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestCachedGetLocations_ReturnsCopy(t *testing.T) {
	stub := &stubListingStore{}
	cache := NewCachedListingStore(stub, 60*time.Second)

	locs1, _ := cache.GetLocations(context.Background())
	// Mutate the returned slice
	locs1[0] = "MUTATED"

	// Get again — should still be the original value
	locs2, _ := cache.GetLocations(context.Background())
	if locs2[0] != "Lagos" {
		t.Errorf("cache was mutated by caller: expected Lagos, got %s", locs2[0])
	}
}
