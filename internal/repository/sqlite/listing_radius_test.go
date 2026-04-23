package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFindAll_RadiusSearch(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Dallas, TX coordinates (approx): 32.7767, -96.7970
	dallasLat := 32.7767
	dallasLng := -96.7970

	// 1. Dallas Listing (Inside radius)
	l1 := domain.Listing{
		ID:        "dallas-1",
		Title:     "Dallas Spot",
		Latitude:  dallasLat,
		Longitude: dallasLng,
		Status:    domain.ListingStatusApproved,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// 2. Fort Worth Listing (~32 miles away)
	// FTW coords: 32.7555, -97.3308
	l2 := domain.Listing{
		ID:        "fort-worth-1",
		Title:     "Fort Worth Spot",
		Latitude:  32.7555,
		Longitude: -97.3308,
		Status:    domain.ListingStatusApproved,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// 3. Austin Listing (~180 miles away)
	// Austin coords: 30.2672, -97.7431
	l3 := domain.Listing{
		ID:        "austin-1",
		Title:     "Austin Spot",
		Latitude:  30.2672,
		Longitude: -97.7431,
		Status:    domain.ListingStatusApproved,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	for _, l := range []domain.Listing{l1, l2, l3} {
		err := repo.Save(ctx, l)
		assert.NoError(t, err)
	}

	// Test 1: Search within 5 miles of Dallas
	// Should only find l1
	res, _, err := repo.FindAll(ctx, "", "", "", dallasLat, dallasLng, 5, "", "", false, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "dallas-1", res[0].ID)

	// Test 2: Search within 50 miles of Dallas
	// Should find l1 and l2
	res, _, err = repo.FindAll(ctx, "", "", "", dallasLat, dallasLng, 50, "", "", false, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	
	ids := map[string]bool{}
	for _, l := range res {
		ids[l.ID] = true
	}
	assert.True(t, ids["dallas-1"])
	assert.True(t, ids["fort-worth-1"])
	assert.False(t, ids["austin-1"])

	// Test 3: Search within 250 miles of Dallas
	// Should find all three
	res, _, err = repo.FindAll(ctx, "", "", "", dallasLat, dallasLng, 250, "", "", false, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
}
