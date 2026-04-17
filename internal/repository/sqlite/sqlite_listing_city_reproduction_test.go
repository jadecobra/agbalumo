package sqlite_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func TestFindAll_CityFiltering_Reproduction(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed Data: One food listing in Houston, one service in Dallas
	saveTestListing(t, ctx, repo, domain.Listing{
		ID:       "1",
		Title:    "Ada's Jollof",
		Type:     domain.Food,
		City:     "Houston",
		Status:   domain.ListingStatusApproved,
		IsActive: true,
	})
	saveTestListing(t, ctx, repo, domain.Listing{
		ID:       "2",
		Title:    "Houston Hair Braiding",
		Type:     domain.Service,
		City:     "Dallas", // Mentions Houston in Title but is in Dallas
		Status:   domain.ListingStatusApproved,
		IsActive: true,
	})

	// Current problem: We can only filter by city using the search text (q)
	// which is inaccurate as it picks up titles/descriptions.
	// Also we can't combine it easily with Type if we want a structured filter.
	
	// We want to filter by City="Houston" (exact)
	// For now, let's see what happens when we search for "Houston" using the current interface
	res, _, err := repo.FindAll(ctx, "", "Houston", "", "", false, 20, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	// It will return BOTH because ID 2 has "Houston" in the title.
	// This proves that we need a dedicated city filter.
	if len(res) > 1 {
		t.Errorf("Expected only 1 listing in Houston, but got %d (FTS returned too much noise)", len(res))
	}
}
