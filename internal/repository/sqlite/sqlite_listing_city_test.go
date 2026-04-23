package sqlite_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFindAll_CityFiltering(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed Data
	saveTestListing(t, ctx, repo, domain.Listing{ID: "1", Title: "Food Houston", Type: domain.Food, City: "Houston", IsActive: true, Status: domain.ListingStatusApproved})
	saveTestListing(t, ctx, repo, domain.Listing{ID: "2", Title: "Service Houston", Type: domain.Service, City: "Houston", IsActive: true, Status: domain.ListingStatusApproved})
	saveTestListing(t, ctx, repo, domain.Listing{ID: "3", Title: "Food Dallas", Type: domain.Food, City: "Dallas", IsActive: true, Status: domain.ListingStatusApproved})

	// 1. Filter by City only
	res, _, err := repo.FindAll(ctx, "", "", "Houston", 0.0, 0.0, 0.0, "", "", false, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	// 2. Filter by City + Category
	res, _, err = repo.FindAll(ctx, string(domain.Food), "", "Houston", 0.0, 0.0, 0.0, "", "", false, 10, 0)
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "1", res[0].ID)
	}

	// 3. Filter by City + Search
	res, _, err = repo.FindAll(ctx, "", "Service", "Houston", 0.0, 0.0, 0.0, "", "", false, 10, 0)
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "2", res[0].ID)
	}
}

func TestGetFeaturedListings_CityAware(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// Seed Data
	saveTestListing(t, ctx, repo, domain.Listing{ID: "f1", Title: "Feat Food Houston", Type: domain.Food, City: "Houston", Featured: true, IsActive: true, Status: domain.ListingStatusApproved})
	saveTestListing(t, ctx, repo, domain.Listing{ID: "f2", Title: "Feat Food Dallas", Type: domain.Food, City: "Dallas", Featured: true, IsActive: true, Status: domain.ListingStatusApproved})
	saveTestListing(t, ctx, repo, domain.Listing{ID: "n1", Title: "Normal Houston", Type: domain.Food, City: "Houston", Featured: false, IsActive: true, Status: domain.ListingStatusApproved})

	// 1. Featured Food in Houston
	res, err := repo.GetFeaturedListings(ctx, string(domain.Food), "Houston")
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "f1", res[0].ID)
	}

	// 2. All Featured in Houston
	res, err = repo.GetFeaturedListings(ctx, "", "Houston")
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "f1", res[0].ID)
	}

	// 3. Global Featured
	res, err = repo.GetFeaturedListings(ctx, "", "")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
