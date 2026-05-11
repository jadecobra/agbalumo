package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func TestListingDeterminism(t *testing.T) {
	ctx := context.Background()
	repo := testutil.SetupTestRepository(t)
	t.Cleanup(func() {
		if err := repo.Close(); err != nil {
			t.Errorf("failed to close repo: %v", err)
		}
	})

	seedDeterminismData(t, repo, ctx)

	t.Run("FindAllByOwner determinism", func(t *testing.T) {
		assertOwnerDeterminism(t, repo, ctx)
	})

	t.Run("GetFeaturedListings determinism", func(t *testing.T) {
		assertFeaturedDeterminism(t, repo, ctx)
	})
}

func seedDeterminismData(t *testing.T, repo domain.ListingRepository, ctx context.Context) {
	t.Helper()
	now := time.Now()
	for i := 0; i < 10; i++ {
		l := domain.Listing{
			ID:          fmt.Sprintf("listing-%d", i),
			Title:       fmt.Sprintf("Listing %d", i),
			Type:        domain.Food,
			OwnerID:     "owner-1",
			OwnerOrigin: "test",
			CreatedAt:   now,
			Status:      domain.ListingStatusApproved,
			IsActive:    true,
		}
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("failed to save listing: %v", err)
		}
	}
}

func assertOwnerDeterminism(t *testing.T, repo domain.ListingReader, ctx context.Context) {
	t.Helper()
	results1, _, err := repo.FindAllByOwner(ctx, "owner-1", 10, 0)
	if err != nil {
		t.Fatalf("FindAllByOwner failed: %v", err)
	}

	results2, _, err := repo.FindAllByOwner(ctx, "owner-1", 10, 0)
	if err != nil {
		t.Fatalf("FindAllByOwner failed: %v", err)
	}

	if len(results1) == 0 {
		t.Fatalf("expected results, got 0")
	}

	for i := range results1 {
		if results1[i].ID != results2[i].ID {
			t.Errorf("non-deterministic order at index %d: %s vs %s", i, results1[i].ID, results2[i].ID)
		}
	}
}

func assertFeaturedDeterminism(t *testing.T, repo domain.ListingRepository, ctx context.Context) {
	t.Helper()
	for i := 0; i < 5; i++ {
		if err := repo.SetFeatured(ctx, fmt.Sprintf("listing-%d", i), true); err != nil {
			t.Fatalf("failed to set featured: %v", err)
		}
	}

	results1, err := repo.GetFeaturedListings(ctx, "Food", "")
	if err != nil {
		t.Fatalf("GetFeaturedListings failed: %v", err)
	}

	results2, err := repo.GetFeaturedListings(ctx, "Food", "")
	if err != nil {
		t.Fatalf("GetFeaturedListings failed: %v", err)
	}

	if len(results1) != 3 || len(results2) != 3 {
		t.Fatalf("expected 3 results, got %d and %d", len(results1), len(results2))
	}

	for i := range results1 {
		if results1[i].ID != results2[i].ID {
			t.Errorf("non-deterministic order at index %d: %s vs %s", i, results1[i].ID, results2[i].ID)
		}
	}
}
