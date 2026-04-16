package sqlite_test

import (
	"context"
	"testing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/jadecobra/agbalumo/internal/seeder"
)

func BenchmarkSQLiteRepository_BulkInsertListings(b *testing.B) {
	repo, _ := testutil.SetupTestRepositoryUnique(b)
	ctx := context.Background()
	count := 10000
	listings := seeder.GenerateStressListings(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.BulkInsertListings(ctx, listings)
		if err != nil {
			b.Fatalf("BulkInsertListings failed: %v", err)
		}
	}
}
