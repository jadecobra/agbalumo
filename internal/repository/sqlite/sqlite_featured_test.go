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
	ctx := context.Background()

	listings := []domain.Listing{
		{
			ID:          "feat-1",
			Title:       "Featured Business",
			Type:        domain.Business,
			IsActive:    true,
			Featured:    true,
			OwnerOrigin: "Nigeria",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-2",
			Title:       "Non-Featured Business",
			Type:        domain.Business,
			IsActive:    true,
			Featured:    false,
			OwnerOrigin: "Ghana",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-3",
			Title:       "Featured but Inactive",
			Type:        domain.Job,
			IsActive:    false,
			Featured:    true,
			OwnerOrigin: "Senegal",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "feat-4",
			Title:       "Another Featured",
			Type:        domain.Product,
			IsActive:    true,
			Featured:    true,
			OwnerOrigin: "Togo",
			CreatedAt:   time.Now().Add(time.Hour),
		},
	}

	for _, l := range listings {
		err := repo.Save(ctx, l)
		assert.NoError(t, err)
	}

	featured, err := repo.GetFeaturedListings(ctx)
	assert.NoError(t, err)

	assert.Len(t, featured, 2)

	assert.Equal(t, "feat-4", featured[0].ID)
	assert.Equal(t, "feat-1", featured[1].ID)
}

func TestGetFeaturedListings_LimitFive(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		l := domain.Listing{
			ID:          string(rune('a' + i)),
			Title:       "Featured",
			Type:        domain.Business,
			IsActive:    true,
			Featured:    true,
			OwnerOrigin: "Nigeria",
			CreatedAt:   time.Now().Add(time.Duration(i) * time.Hour),
		}
		err := repo.Save(ctx, l)
		assert.NoError(t, err)
	}

	featured, err := repo.GetFeaturedListings(ctx)
	assert.NoError(t, err)

	assert.LessOrEqual(t, len(featured), 5, "GetFeaturedListings should return at most 5 items")
}

func TestSetFeatured(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:          "feat-test",
		Title:       "Test Listing",
		Type:        domain.Business,
		IsActive:    true,
		Featured:    false,
		OwnerOrigin: "Nigeria",
		CreatedAt:   time.Now(),
	}
	err := repo.Save(ctx, l)
	assert.NoError(t, err)

	err = repo.SetFeatured(ctx, "feat-test", true)
	assert.NoError(t, err)

	featured, err := repo.GetFeaturedListings(ctx)
	assert.NoError(t, err)
	assert.Len(t, featured, 1)
	assert.Equal(t, "feat-test", featured[0].ID)

	err = repo.SetFeatured(ctx, "feat-test", false)
	assert.NoError(t, err)

	featured, err = repo.GetFeaturedListings(ctx)
	assert.NoError(t, err)
	assert.Len(t, featured, 0)
}
