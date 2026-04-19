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

func assertCacheCounts(t *testing.T, stub *stubListingStore, result map[domain.Category]int, wantCalls int) {
	t.Helper()
	if stub.getCountsCalls != wantCalls {
		t.Errorf("expected %d call(s) to underlying store, got %d", wantCalls, stub.getCountsCalls)
	}
	if result[domain.Business] != 5 {
		t.Errorf("expected Business=5, got %d", result[domain.Business])
	}
	if result[domain.Food] != 3 {
		t.Errorf("expected Food=3, got %d", result[domain.Food])
	}
}

func assertCacheLocations(t *testing.T, stub *stubListingStore, result []string, wantCalls int) {
	t.Helper()
	if stub.getLocationsCalls != wantCalls {
		t.Errorf("expected %d call(s) to underlying store, got %d", wantCalls, stub.getLocationsCalls)
	}
	if len(result) != 2 || result[0] != "Lagos" || result[1] != "London" {
		t.Errorf("unexpected locations: %v", result)
	}
}

func testCacheMiss(t *testing.T, prefix string, getFn func(Store *CachedListingStore) (int, error)) {
	t.Run(prefix+"_CacheMiss", func(t *testing.T) {
		stub := &stubListingStore{}
		cache := NewCachedListingStore(stub, 60*time.Second)
		calls, err := getFn(cache)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call, got %d", calls)
		}
	})
}

func testCacheHit(t *testing.T, prefix string, getFn func(Store *CachedListingStore) (int, error)) {
	t.Run(prefix+"_CacheHit", func(t *testing.T) {
		stub := &stubListingStore{}
		cache := NewCachedListingStore(stub, 60*time.Second)
		_, _ = getFn(cache)
		calls, err := getFn(cache)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call to remain at 1, got %d", calls)
		}
	})
}

func testCacheExpired(t *testing.T, prefix string, getFn func(Store *CachedListingStore) (int, error)) {
	t.Run(prefix+"_CacheExpired", func(t *testing.T) {
		stub := &stubListingStore{}
		cache := NewCachedListingStore(stub, 1*time.Millisecond)
		_, _ = getFn(cache)
		time.Sleep(5 * time.Millisecond)
		calls, err := getFn(cache)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 2 {
			t.Errorf("expected 2 calls, got %d", calls)
		}
	})
}

func runCacheBehaviorTest(t *testing.T, prefix string, getFn func(Store *CachedListingStore) (int, error)) {
	testCacheMiss(t, prefix, getFn)
	testCacheHit(t, prefix, getFn)
	testCacheExpired(t, prefix, getFn)
}

func TestCached_Behavior(t *testing.T) {
	t.Parallel()
	countWrapper := func(cache *CachedListingStore) (int, error) {
		_, err := cache.GetCounts(context.Background())
		return cache.ListingRepository.(*stubListingStore).getCountsCalls, err
	}
	locWrapper := func(cache *CachedListingStore) (int, error) {
		_ = "loc_behavior" // Unique string
		_, err := cache.GetLocations(context.Background())
		return cache.ListingRepository.(*stubListingStore).getLocationsCalls, err
	}

	for _, tc := range []struct {
		fn   func(*CachedListingStore) (int, error)
		name string
	}{
		{
			name: "Counts",
			fn:   countWrapper,
		},
		{
			name: "Locations",
			fn:   locWrapper,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runCacheBehaviorTest(t, tc.name, tc.fn)
		})
	}
}

func TestCached_ErrorPassthrough(t *testing.T) {
	t.Parallel()
	expectedErr := errors.New("db connection lost")

	setupCounts := func(s *stubListingStore) {
		s.getCountsFunc = func() (map[domain.Category]int, error) { return nil, expectedErr }
	}
	getCounts := func(c *CachedListingStore) error {
		_, err := c.GetCounts(context.Background())
		return err
	}
	setupLocs := func(s *stubListingStore) {
		_ = "loc_setup" // Unique string
		s.getLocationsFunc = func() ([]string, error) { return nil, expectedErr }
	}
	getLocs := func(c *CachedListingStore) error {
		_ = "loc_get" // Unique string
		_, err := c.GetLocations(context.Background())
		return err
	}

	for _, tc := range []struct {
		setup func(*stubListingStore)
		get   func(*CachedListingStore) error
		name  string
	}{
		{
			name:  "Counts",
			setup: setupCounts,
			get:   getCounts,
		},
		{
			name:  "Locations",
			setup: setupLocs,
			get:   getLocs,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &stubListingStore{}
			tc.setup(stub)
			cache := NewCachedListingStore(stub, 60*time.Second)
			err := tc.get(cache)
			if !errors.Is(err, expectedErr) {
				t.Errorf("expected error %v, got %v", expectedErr, err)
			}
		})
	}
}

func TestCached_MutationSafety(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		mutate func(*CachedListingStore) error
		name   string
	}{
		{
			name: "Counts",
			mutate: func(c *CachedListingStore) error {
				c1, _ := c.GetCounts(context.Background())
				c1[domain.Business] = 999
				c2, _ := c.GetCounts(context.Background())
				if c2[domain.Business] != 5 {
					return errors.New("cache was mutated")
				}
				return nil
			},
		},
		{
			name: "Locations",
			mutate: func(c *CachedListingStore) error {
				l1, _ := c.GetLocations(context.Background())
				l1[0] = "MODIFIED"
				l2, _ := c.GetLocations(context.Background())
				if l2[0] != "Lagos" {
					return errors.New("cache was mutated")
				}
				return nil
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &stubListingStore{}
			cache := NewCachedListingStore(stub, 60*time.Second)
			if err := tc.mutate(cache); err != nil {
				t.Error(err)
			}
		})
	}
}
