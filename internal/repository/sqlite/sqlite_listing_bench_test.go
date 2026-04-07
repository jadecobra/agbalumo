package sqlite_test

import (
	"github.com/jadecobra/agbalumo/internal/testutil"
	"context"
	"fmt"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func BenchmarkSQLiteRepository_FindAll(b *testing.B) {
	repo, _ := testutil.SetupTestRepositoryUnique(b)
	ctx := context.Background()

	// Seed 100 listings
	for i := 0; i < 100; i++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.FindAll(ctx, "", "", "", "", true, 50, 0)
	}
}

func BenchmarkSQLiteRepository_FindByTitle(b *testing.B) {
	repo, _ := testutil.SetupTestRepositoryUnique(b)
	ctx := context.Background()

	// Seed 100 listings
	for i := 0; i < 100; i++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByTitle(ctx, "Listing 50")
	}
}
