package seeder_test

import (
	"context"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/seeder"
)

func TestEnsureCategoriesSeeded_HappyPath(t *testing.T) {
	t.Parallel()
	repo, configPath, cleanup := setupSeeder(t)
	defer cleanup()
	ctx := context.Background()

	// Create a temp config file
	content := `[
		{"id": "business", "name": "Business", "claimable": true, "is_system": true, "active": true},
		{"id": "event", "name": "Event", "claimable": false, "is_system": true, "active": true, "requires_special_validation": true}
	]`
	if err := os.WriteFile( /*nolint:gosec*/ configPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	err := seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
	if err != nil {
		t.Fatalf("EnsureCategoriesSeeded failed: %v", err)
	}

	// Verify categories were seeded
	cats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(cats))
	}
}

func TestEnsureCategoriesSeeded_FileNotFound(t *testing.T) {
	t.Parallel()
	repo, _, cleanup := setupSeeder(t)
	defer cleanup()
	ctx := context.Background()

	// Should not error when file doesn't exist
	err := seeder.EnsureCategoriesSeeded(ctx, repo, "/nonexistent/path/categories.json")
	if err != nil {
		t.Errorf("Expected no error for missing file, got %v", err)
	}
}

func TestEnsureCategoriesSeeded_InvalidJSON(t *testing.T) {
	t.Parallel()
	repo, configPath, cleanup := setupSeeder(t)
	defer cleanup()
	ctx := context.Background()

	if err := os.WriteFile( /*nolint:gosec*/ configPath, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	err := seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestEnsureCategoriesSeeded_Idempotent(t *testing.T) {
	t.Parallel()
	repo, configPath, cleanup := setupSeeder(t)
	defer cleanup()
	ctx := context.Background()

	content := `[{"id": "food", "name": "Food", "claimable": false, "is_system": true, "active": true}]`
	if err := os.WriteFile( /*nolint:gosec*/ configPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Run twice
	_ = seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
	_ = seeder.EnsureCategoriesSeeded(ctx, repo, configPath)

	cats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("Expected 1 category after idempotent seed, got %d", len(cats))
	}
}
