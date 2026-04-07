package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

// saveTestListing is a helper that wraps repo.Save and calls t.Fatal on error.
func saveTestListing(t *testing.T, ctx context.Context, repo domain.ListingWriter, listing domain.Listing) {
	t.Helper()
	if listing.ID == "" {
		listing.ID = "test-id"
	}
	if listing.Title == "" {
		listing.Title = "Test Listing"
	}
	if listing.CreatedAt.IsZero() {
		listing.CreatedAt = time.Now()
	}
	if err := repo.Save(ctx, listing); err != nil {
		t.Fatalf("Failed to save test listing %q: %v", listing.ID, err)
	}
}

// saveTestCategory is a helper that wraps repo.SaveCategory and calls t.Fatal on error.
func saveTestCategory(t *testing.T, ctx context.Context, repo domain.CategoryStore, cat domain.CategoryData) {
	t.Helper()
	if cat.ID == "" {
		cat.ID = "test-cat"
	}
	if cat.Name == "" {
		cat.Name = "Test Category"
	}
	if cat.CreatedAt.IsZero() {
		cat.CreatedAt = time.Now()
	}
	if cat.UpdatedAt.IsZero() {
		cat.UpdatedAt = time.Now()
	}
	if err := repo.SaveCategory(ctx, cat); err != nil {
		t.Fatalf("Failed to save test category %q: %v", cat.ID, err)
	}
}

// setupBenchmarkDB seeds the DB with n listings for benchmarks.
func setupBenchmarkDB(b *testing.B, n int) (domain.ListingRepository, context.Context, func()) {
	b.Helper()
	repo, _ := testutil.SetupTestRepositoryUnique(b)
	ctx := context.Background()

	for i := 0; i < n; i++ {
		_ = repo.Save(ctx, domain.Listing{
			ID:          fmt.Sprintf("l%d", i),
			Title:       fmt.Sprintf("Listing %d", i),
			Type:        domain.Business,
			OwnerOrigin: "Nigeria",
			Address:     "123 St",
			Description: "Desc",
			IsActive:    true,
		})
	}
	return repo, ctx, func() { _ = repo.Close() }
}
