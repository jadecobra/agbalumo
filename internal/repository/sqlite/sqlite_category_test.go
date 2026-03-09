package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func TestSaveCategory(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	cat := domain.CategoryData{
		ID:        "test-cat",
		Name:      "Test Category",
		Claimable: true,
		IsSystem:  false,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.SaveCategory(ctx, cat); err != nil {
		t.Fatalf("SaveCategory failed: %v", err)
	}

	// Verify it was saved
	got, err := repo.GetCategory(ctx, "test-cat")
	if err != nil {
		t.Fatalf("GetCategory failed: %v", err)
	}
	if got.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", got.Name)
	}
	if !got.Claimable {
		t.Error("Expected claimable to be true")
	}
	if !got.Active {
		t.Error("Expected active to be true")
	}
}

func TestSaveCategory_Upsert(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	cat := domain.CategoryData{
		ID:        "upsert-cat",
		Name:      "Original",
		Claimable: false,
		IsSystem:  false,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.SaveCategory(ctx, cat); err != nil {
		t.Fatalf("SaveCategory failed: %v", err)
	}

	// Update via upsert
	cat.Name = "Updated"
	cat.Claimable = true
	cat.UpdatedAt = time.Now()
	if err := repo.SaveCategory(ctx, cat); err != nil {
		t.Fatalf("SaveCategory update failed: %v", err)
	}

	got, err := repo.GetCategory(ctx, "upsert-cat")
	if err != nil {
		t.Fatalf("GetCategory failed: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", got.Name)
	}
	if !got.Claimable {
		t.Error("Expected claimable to be true after update")
	}
}

func TestGetCategories(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed categories
	cats := []domain.CategoryData{
		{ID: "bus", Name: "Business", Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "svc", Name: "Service", Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "old", Name: "Old Category", Active: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, c := range cats {
		if err := repo.SaveCategory(ctx, c); err != nil {
			t.Fatalf("Failed to seed category %s: %v", c.Name, err)
		}
	}

	// Get all (no filter)
	all, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err != nil {
		t.Fatalf("GetCategories all failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(all))
	}

	// Get active only
	active, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		t.Fatalf("GetCategories active failed: %v", err)
	}
	if len(active) != 2 {
		t.Errorf("Expected 2 active categories, got %d", len(active))
	}

	// Verify ordering (ASC by name)
	if len(active) >= 2 && active[0].Name > active[1].Name {
		t.Errorf("Expected categories ordered by name ASC, got %s before %s", active[0].Name, active[1].Name)
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetCategory(ctx, "non-existent")
	if err != domain.ErrCategoryNotFound {
		t.Errorf("Expected ErrCategoryNotFound, got %v", err)
	}
}

func TestUpsertCoreCategory(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	core := domain.CategoryData{
		ID:                        "job",
		Name:                      "Job",
		Claimable:                 false,
		IsSystem:                  true,
		Active:                    true,
		RequiresSpecialValidation: true,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	if err := repo.UpsertCoreCategory(ctx, core); err != nil {
		t.Fatalf("UpsertCoreCategory failed: %v", err)
	}

	got, err := repo.GetCategory(ctx, "job")
	if err != nil {
		t.Fatalf("GetCategory failed: %v", err)
	}
	if got.Name != "Job" {
		t.Errorf("Expected name 'Job', got '%s'", got.Name)
	}
	if !got.IsSystem {
		t.Error("Expected is_system to be true")
	}
	if !got.RequiresSpecialValidation {
		t.Error("Expected requires_special_validation to be true")
	}

	// Upsert again with updated name - should now also update active if we want to enforce it
	core.Name = "Job Posting"
	core.Active = false // admin disabled it (simulated)
	err = repo.UpsertCoreCategory(ctx, core)
	if err != nil {
		t.Fatalf("UpsertCoreCategory update failed: %v", err)
	}

	got, err = repo.GetCategory(ctx, "job")
	if err != nil {
		t.Fatalf("GetCategory after upsert failed: %v", err)
	}
	if got.Name != "Job Posting" {
		t.Errorf("Expected updated name 'Job Posting', got '%s'", got.Name)
	}
	// With the fix, active SHOULD now be false because UpsertCoreCategory overrides it
	if got.Active {
		t.Error("UpsertCoreCategory should override 'active' flag with the provided value")
	}
}

func TestCategoryErrors(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	_ = db.Close()
	ctx := context.Background()

	checkError := func(name string, err error) {
		if err == nil {
			t.Errorf("%s: Expected error on closed DB, got nil", name)
		}
	}

	checkError("SaveCategory", repo.SaveCategory(ctx, domain.CategoryData{ID: "x"}))
	_, err = repo.GetCategories(ctx, domain.CategoryFilter{})
	checkError("GetCategories", err)
	_, err = repo.GetCategory(ctx, "x")
	checkError("GetCategory", err)
	checkError("UpsertCoreCategory", repo.UpsertCoreCategory(ctx, domain.CategoryData{ID: "x"}))
	_, err = repo.GetLocations(ctx)
	checkError("GetLocations", err)
}

func TestGetLocations(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed listings with cities
	_ = repo.SaveCategory(ctx, domain.CategoryData{ID: "bus", Name: "Business", Active: true})
	listings := []domain.Listing{
		{ID: "1", Title: "Place A", City: "Houston", Type: "Business", IsActive: true},
		{ID: "2", Title: "Place B", City: "Dallas", Type: "Business", IsActive: true},
		{ID: "3", Title: "Place C", City: "Houston", Type: "Business", IsActive: true},
		{ID: "4", Title: "Place D", City: "", Type: "Business", IsActive: true},        // no city
		{ID: "5", Title: "Place E", City: "Austin", Type: "Business", IsActive: false}, // inactive
	}
	for _, l := range listings {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to seed listing: %v", err)
		}
	}

	locations, err := repo.GetLocations(ctx)
	if err != nil {
		t.Fatalf("GetLocations failed: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 unique active locations (Houston, Dallas), got %d: %v", len(locations), locations)
	}
}

func TestGetCategories_EmptyDB(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	cats, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(cats) != 0 {
		t.Errorf("Expected 0 categories on empty DB, got %d", len(cats))
	}
}

// TestSaveCategory_PreservesExistingCategories verifies that adding a new category
// does NOT overwrite pre-existing categories. This is the regression test for the
// bug where HandleAddCategory used an empty ID, causing ON CONFLICT(id) overwrites.
func TestSaveCategory_PreservesExistingCategories(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed the core categories (simulates what EnsureCategoriesSeeded does)
	coreCategories := []domain.CategoryData{
		{ID: "Business", Name: "Business", IsSystem: true, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "Service", Name: "Service", IsSystem: true, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "Food", Name: "Food", IsSystem: true, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, c := range coreCategories {
		if err := repo.SaveCategory(ctx, c); err != nil {
			t.Fatalf("Failed to seed core category %s: %v", c.ID, err)
		}
	}

	// Now add a new admin-created category with a proper unique ID
	newCat := domain.CategoryData{
		ID:        "music",
		Name:      "Music",
		IsSystem:  false,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.SaveCategory(ctx, newCat); err != nil {
		t.Fatalf("Failed to save new category: %v", err)
	}

	// Verify ALL 4 categories exist
	all, err := repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	if len(all) != 4 {
		t.Errorf("Expected 4 categories after adding one new, got %d", len(all))
		for _, c := range all {
			t.Logf("  category: ID=%q Name=%q", c.ID, c.Name)
		}
	}

	// Verify each core category still exists
	for _, expected := range coreCategories {
		got, err := repo.GetCategory(ctx, expected.ID)
		if err != nil {
			t.Errorf("Core category %q missing after adding new category: %v", expected.ID, err)
			continue
		}
		if got.Name != expected.Name {
			t.Errorf("Core category %q name changed from %q to %q", expected.ID, expected.Name, got.Name)
		}
	}
}
