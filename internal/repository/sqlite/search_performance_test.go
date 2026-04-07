package sqlite_test

import (
	"context"
	"fmt"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func BenchmarkSearchPerformance(b *testing.B) {
	// Setup a unique in-memory DB for the benchmark using standardized helper
	repo, _ := testutil.SetupTestRepositoryUnique(b)
	ctx := context.Background()

	// Seed listings (default 10000, can be overridden for smoke tests)
	numListings := 10000
	if val := os.Getenv("BENCH_LISTINGS"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			numListings = n
		}
	}

	b.Logf("Seeding %d listings...", numListings)
	if err := seedBenchmarkData(ctx, repo, numListings); err != nil {
		b.Fatalf("Failed to seed: %v", err)
	}
	b.Log("Seeding complete.")

	b.Run("FindAll_Default_Page1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = repo.FindAll(ctx, "", "", "", "", false, 30, 0)
		}
	})

	b.Run("FindAll_Search_Page1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = repo.FindAll(ctx, "", "ghana", "", "", false, 30, 0)
		}
	})

	b.Run("FindAll_Filter_Page1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = repo.FindAll(ctx, string(domain.Business), "", "", "", false, 30, 0)
		}
	})

	b.Run("FindAll_Search_Filter_Page1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = repo.FindAll(ctx, string(domain.Business), "ghana", "", "", false, 30, 0)
		}
	})

	b.Run("FindAll_Deep_Pagination", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = repo.FindAll(ctx, "", "", "", "", false, 30, 5000)
		}
	})
}

func seedBenchmarkData(ctx context.Context, repo *sqlite.SQLiteRepository, numListings int) error {
	for i := 0; i < numListings; i++ {
		l := domain.Listing{
			ID:          fmt.Sprintf("listing-%d", i),
			Title:       fmt.Sprintf("Business Listing %d", i),
			Description: fmt.Sprintf("Description for business %d with some common keywords like food service ghana nigeria.", i),
			Type:        domain.Business,
			OwnerOrigin: "Nigeria",
			City:        "Houston",
			Address:     fmt.Sprintf("%d Main St", i),
			IsActive:    true,
			Status:      domain.ListingStatusApproved,
			CreatedAt:   time.Now().Add(time.Duration(-i) * time.Hour),
		}
		if i%5 == 0 {
			l.Type = domain.Service
		}
		if i%10 == 0 {
			l.Status = domain.ListingStatusPending
		}
		if i%20 == 0 {
			l.IsActive = false
		}

		if err := repo.Save(ctx, l); err != nil {
			return err
		}
	}
	return nil
}
