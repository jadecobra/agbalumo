package sqlite_test

import (
	"context"
	"database/sql"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	_ "modernc.org/sqlite"
)

func TestReproCategoryRegression(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// 1. Seed core categories (simulates startup/seeder)
	coreCats := []domain.CategoryData{
		{ID: "Business", Name: "Business", IsSystem: true, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "Service", Name: "Service", IsSystem: true, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, c := range coreCats {
		if err := repo.UpsertCoreCategory(ctx, c); err != nil {
			t.Fatalf("Failed to seed core category %s: %v", c.ID, err)
		}
	}

	// 2. Verify they are active
	active, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(active) != 2 {
		t.Fatalf("Expected 2 active categories, got %d", len(active))
	}

	// 3. Add "church" as admin does
	name := "church"
	newCat := domain.CategoryData{
		ID:        "church", // strings.ToLower(...)
		Name:      name,
		Claimable: true,
		IsSystem:  false,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.SaveCategory(ctx, newCat)
	if err != nil {
		t.Fatalf("Failed to save church category: %v", err)
	}

	// 4. Check if core categories are still active
	activeAfter, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		t.Fatalf("GetCategories after failed: %v", err)
	}
	if len(activeAfter) != 3 {
		t.Errorf("Expected 3 active categories, got %d", len(activeAfter))
		for _, c := range activeAfter {
			t.Logf("  category: ID=%q Name=%q Active=%v", c.ID, c.Name, c.Active)
		}
	}
}

func TestCategoryCaseSensitivity(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// 1. Seed core category "Business"
	core := domain.CategoryData{ID: "Business", Name: "Business", IsSystem: true, Active: true}
	if err := repo.UpsertCoreCategory(ctx, core); err != nil {
		t.Fatalf("Failed to seed: %v", err)
	}

	// 2. Add custom category "business"
	custom := domain.CategoryData{ID: "business", Name: "business", IsSystem: false, Active: true}
	if err := repo.SaveCategory(ctx, custom); err != nil {
		t.Fatalf("Failed to save custom: %v", err)
	}

	// 3. Check if both exist
	all, _ := repo.GetCategories(ctx, domain.CategoryFilter{})
	if len(all) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(all))
		for _, c := range all {
			t.Logf("  ID=%q Name=%q", c.ID, c.Name)
		}
	}
}

func TestUpsertCoreCategory_ActiveOverwrite(t *testing.T) {
	repo, dbPath := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	// 1. Manually insert an INACTIVE core category
	// (Simulating an old DB state or manual deactivation)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open DB directly: %v", err)
	}
	defer func() { _ = db.Close() }()

	_, err = db.ExecContext(ctx,
		"INSERT INTO categories (id, name, active) VALUES (?, ?, ?)",
		"Business", "Business", 0)
	if err != nil {
		t.Fatalf("Failed to insert inactive category: %v", err)
	}

	// 2. Run UpsertCoreCategory which HAS active: true in the struct (from categories.json)
	core := domain.CategoryData{ID: "Business", Name: "Business", IsSystem: true, Active: true}
	if err := repo.UpsertCoreCategory(ctx, core); err != nil {
		t.Fatalf("UpsertCoreCategory failed: %v", err)
	}

	// 3. Check if it's now active
	got, _ := repo.GetCategory(ctx, "Business")
	if !got.Active {
		t.Errorf("BUG FOUND: UpsertCoreCategory DID NOT update active status to true!")
	}
}
