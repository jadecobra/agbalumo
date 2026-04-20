package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func TestFindAll_CityFilter_Regression(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed data matching the issue: "Dallas" is in the address, but City field is empty.
	listings := []domain.Listing{
		{
			ID:          "dallas-1",
			Title:       "Lalibela Restaurant & Bar",
			Type:        domain.Food,
			City:        "Dallas",
			Address:     "9191 Forest Ln # 2, Dallas, TX 75243",
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "dallas-2",
			Title:       "African Food Depot",
			Type:        domain.Food,
			City:        "", // EMPTY CITY - This is the bug scenario
			Address:     "9751 Walnut St #112, Dallas, TX 75243",
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "dallas-3",
			Title:       "Southwest Farmers Market",
			Type:        domain.Food,
			City:        "", // EMPTY CITY
			Address:     "4460 W Walnut St, Garland, TX 75042", // This one is Garland, but close. Let's use one that says Dallas in address.
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "dallas-4",
			Title:       "Murphy's Mansion",
			Type:        domain.Food,
			City:        "", // EMPTY CITY
			Address:     "10051 Whitehurst Dr, Dallas, TX 75243",
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
	}

	for _, l := range listings {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to save listing %s: %v", l.ID, err)
		}
	}

	// EXECUTION: Filter by Dallas
	res, _, err := repo.FindAll(ctx, string(domain.Food), "", "Dallas", "", "", false, 20, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	// VERIFICATION: We expect 3 results if the fix for missing city data is implemented.
	// CURRENTLY (expected to fail): It will only return 1 (Lalibela).
	if len(res) < 3 {
		t.Errorf("Expected at least 3 Dallas listings, got %d. Only ID %s was found if 1.", len(res), func() string {
			if len(res) > 0 { return res[0].ID }
			return "none"
		}())
	}
}

func TestFindAll_SearchSuya_Regression(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed data with "Suya" in various fields
	listings := []domain.Listing{
		{
			ID:          "suya-1",
			Title:       "Brooklyn Suya",
			Description: "Best Suya in Brooklyn",
			Type:        domain.Food,
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "suya-2",
			Title:       "Mama J's",
			Description: "We serve authentic Nigerian Suya and Jollof.",
			Type:        domain.Food,
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
			CreatedAt:   time.Now(),
		},
	}

	for _, l := range listings {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to save listing %s: %v", l.ID, err)
		}
	}

	// EXECUTION: Search for "Suya"
	res, _, err := repo.FindAll(ctx, "", "Suya", "", "", "", false, 20, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	// VERIFICATION: Expected 2 results.
	// CURRENTLY (reported as zero results in production):
	if len(res) != 2 {
		t.Errorf("Expected 2 search results for 'Suya', got %d", len(res))
	}
}
