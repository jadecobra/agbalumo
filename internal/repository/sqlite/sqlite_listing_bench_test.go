package sqlite_test

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"testing"
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
