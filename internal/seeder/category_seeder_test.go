package seeder_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	_ "github.com/mattn/go-sqlite3"
)

func newTestRepoForSeeder(t *testing.T) *sqlite.SQLiteRepository {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	dsn := dbPath + "?_time_format=sqlite"
	repo, err := sqlite.NewSQLiteRepository(dsn)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	return repo
}

func TestEnsureCategoriesSeeded_HappyPath(t *testing.T) {
	repo := newTestRepoForSeeder(t)
	ctx := context.Background()

	// Create a temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "categories.json")
	content := `[
		{"id": "business", "name": "Business", "claimable": true, "is_system": true, "active": true},
		{"id": "event", "name": "Event", "claimable": false, "is_system": true, "active": true, "requires_special_validation": true}
	]`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
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
	repo := newTestRepoForSeeder(t)
	ctx := context.Background()

	// Should not error when file doesn't exist
	err := seeder.EnsureCategoriesSeeded(ctx, repo, "/nonexistent/path/categories.json")
	if err != nil {
		t.Errorf("Expected no error for missing file, got %v", err)
	}
}

func TestEnsureCategoriesSeeded_InvalidJSON(t *testing.T) {
	repo := newTestRepoForSeeder(t)
	ctx := context.Background()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(configPath, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	err := seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestEnsureCategoriesSeeded_Idempotent(t *testing.T) {
	repo := newTestRepoForSeeder(t)
	ctx := context.Background()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "categories.json")
	content := `[{"id": "food", "name": "Food", "claimable": false, "is_system": true, "active": true}]`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Run twice
	seeder.EnsureCategoriesSeeded(ctx, repo, configPath)
	seeder.EnsureCategoriesSeeded(ctx, repo, configPath)

	cats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("Expected 1 category after idempotent seed, got %d", len(cats))
	}
}
