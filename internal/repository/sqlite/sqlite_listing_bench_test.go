package sqlite_test

import (
	"testing"
)

func BenchmarkSQLiteRepository_FindAll(b *testing.B) {
	repo, ctx, cleanup := setupBenchmarkDB(b, 100)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.FindAll(ctx, "", "", "", "", true, 50, 0)
	}
}

func BenchmarkSQLiteRepository_FindByTitle(b *testing.B) {
	repo, ctx, cleanup := setupBenchmarkDB(b, 100)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByTitle(ctx, "Listing 50")
	}
}
