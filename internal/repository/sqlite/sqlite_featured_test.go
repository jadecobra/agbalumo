package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetFeaturedListings(t *testing.T) {
	repo, _ := newTestRepo(t)
	// defer cleanup is handled by t.TempDir() in newTestRepo automatically cleaning up files,
	// though DB connection closing isn't explicit in newTestRepo, it's file based.

	ctx := context.Background()

	// 1. Create mixed listings
	// - Active Business (Should appear)
	// - Inactive Business (Should NOT appear)
	// - Active Job (Should NOT appear)
	// - Active Product (Should appear)

	listings := []domain.Listing{
		{
			ID:          "feat-1",
			Title:       "Active Business",
			Type:        domain.Business,
			IsActive:    true,
			OwnerOrigin: "Nigeria",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-2",
			Title:       "Inactive Business",
			Type:        domain.Business,
			IsActive:    false,
			OwnerOrigin: "Ghana",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-3",
			Title:       "Active Job",
			Type:        domain.Job,
			IsActive:    true,
			OwnerOrigin: "Senegal",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-4",
			Title:       "Active Product",
			Type:        domain.Product,
			IsActive:    true,
			OwnerOrigin: "Togo",
			CreatedAt:   time.Now().Add(time.Hour), // Newer
		},
	}

	for _, l := range listings {
		err := repo.Save(ctx, l)
		assert.NoError(t, err)
	}

	// 2. Fetch Featured
	featured, err := repo.GetFeaturedListings(ctx)
	assert.NoError(t, err)

	// 3. Verify
	// Should have 2 items: feat-1 and feat-4
	assert.Len(t, featured, 2)

	// Feat-4 is newer, should be first
	assert.Equal(t, "feat-4", featured[0].ID)
	assert.Equal(t, "feat-1", featured[1].ID)
}
