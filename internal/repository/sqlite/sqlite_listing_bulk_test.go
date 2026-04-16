package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/jadecobra/agbalumo/internal/seeder"
)

func TestSQLiteRepository_BulkInsertListings(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// 1. Generate 100 smoke-test listings (reduced from 10,000 to keep CI fast)
	count := 100
	listings := seeder.GenerateStressListings(count)

	// 2. Insert using BulkInsertListings
	start := time.Now()
	err := repo.BulkInsertListings(ctx, listings)
	if err != nil {
		t.Fatalf("BulkInsertListings failed: %v", err)
	}
	duration := time.Since(start)
	t.Logf("Bulk inserted %d listings in %v", count, duration)

	// 3. Verify total count increased appropriately
	_, totalCount, err := repo.FindAll(ctx, "", "", "", "", true, 1, 0)
	if err != nil {
		t.Fatalf("Failed to count listings: %v", err)
	}

	if totalCount != count {
		t.Errorf("Expected %d listings to be inserted, but got %d", count, totalCount)
	}
}
