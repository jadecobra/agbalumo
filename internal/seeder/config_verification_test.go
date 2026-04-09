package seeder_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	_ "modernc.org/sqlite"
)

func TestEnsureCategoriesSeeded_Verification(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repo, configPath, cleanup := setupSeeder(t)
	t.Cleanup(cleanup)

	t.Run("MissingConfig", func(t *testing.T) {
		err := seeder.EnsureCategoriesSeeded(ctx, repo, "non-existent.json")
		if err != nil {
			t.Errorf("Expected no error when config is missing, got %v", err)
		}

		// ASSERTION: Verify DB is still empty
		dbCats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
		if err != nil {
			t.Fatalf("Failed to get categories: %v", err)
		}
		if len(dbCats) != 0 {
			t.Errorf("Expected 0 categories when config is missing, got %d", len(dbCats))
		}
	})

	verifyValidConfigRun(t, ctx, repo, configPath)
}

func verifyValidConfigRun(t *testing.T, ctx context.Context, repo *sqlite.SQLiteRepository, configPath string) {
	t.Run("ValidConfig", func(t *testing.T) {
		cats := []domain.CategoryData{
			{ID: "Ver-1", Name: "Verify One", IsSystem: true, Active: true},
			{ID: "Ver-2", Name: "Verify Two", IsSystem: true, Active: true},
		}
		data, _ := json.Marshal(cats)
		_ = os.WriteFile( /*nolint:gosec*/ configPath, data, 0600)

		err := seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
		if err != nil {
			t.Fatalf("Failed to seed categories: %v", err)
		}

		// Verify count
		dbCats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
		if err != nil {
			t.Fatalf("Failed to get categories: %v", err)
		}

		if len(dbCats) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(dbCats))
		}
	})
}
