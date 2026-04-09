package sqlite_test

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetFeaturedListings(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

	featured, err := repo.GetFeaturedListings(ctx, "")
	assert.NoError(t, err)

	assert.Len(t, featured, 2)

	assert.Equal(t, "feat-4", featured[0].ID)
	assert.Equal(t, "feat-1", featured[1].ID)
}

func TestGetFeaturedListings_LimitThree(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
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

	featured, err := repo.GetFeaturedListings(ctx, "")
	assert.NoError(t, err)

	assert.Len(t, featured, 3, "GetFeaturedListings should return exactly 3 items limit")
}

func TestGetFeaturedListings_CategoryFilter(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	listings := []domain.Listing{
		{
			ID: "biz-1", Title: "Biz 1", Type: domain.Business, IsActive: true, Featured: true,
		},
		{
			ID: "biz-2", Title: "Biz 2", Type: domain.Business, IsActive: true, Featured: true,
		},
		{
			ID: "event-1", Title: "Event 1", Type: domain.Event, IsActive: true, Featured: true,
		},
		{
			ID: "event-2", Title: "Event 2", Type: domain.Event, IsActive: true, Featured: true,
		},
		{
			ID: "job-orig", Title: "Job 1", Type: domain.Job, IsActive: true, Featured: true,
		},
	}

	for _, l := range listings {
		err := repo.Save(ctx, l)
		assert.NoError(t, err)
	}

	// Filter by Business
	bizFeatured, err := repo.GetFeaturedListings(ctx, string(domain.Business))
	assert.NoError(t, err)
	assert.Len(t, bizFeatured, 2)
	assert.Equal(t, domain.Business, bizFeatured[0].Type)
	assert.Equal(t, domain.Business, bizFeatured[1].Type)

	// Filter by Event
	eventFeatured, err := repo.GetFeaturedListings(ctx, string(domain.Event))
	assert.NoError(t, err)
	assert.Len(t, eventFeatured, 2)
	assert.Equal(t, domain.Event, eventFeatured[0].Type)
	assert.Equal(t, domain.Event, eventFeatured[1].Type)

	// Filter by Job
	jobFeatured, err := repo.GetFeaturedListings(ctx, string(domain.Job))
	assert.NoError(t, err)
	assert.Len(t, jobFeatured, 1)
	assert.Equal(t, domain.Job, jobFeatured[0].Type)

	// No filter, should return 3 items (mixed limit)
	allFeatured, err := repo.GetFeaturedListings(ctx, "")
	assert.NoError(t, err)
	assert.Len(t, allFeatured, 3)
}

func TestSetFeatured(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

	featured, err := repo.GetFeaturedListings(ctx, "")
	assert.NoError(t, err)
	assert.Len(t, featured, 1)
	assert.Equal(t, "feat-test", featured[0].ID)

	err = repo.SetFeatured(ctx, "feat-test", false)
	assert.NoError(t, err)

	featured, err = repo.GetFeaturedListings(ctx, "")
	assert.NoError(t, err)
	assert.Len(t, featured, 0)
}
