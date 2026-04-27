package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestSaveAndFindByID(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	l := domain.Listing{
		ID:           "test-1",
		Title:        "Original Title",
		OwnerOrigin:  "Ghana",
		Type:         domain.Business,
		IsActive:     true,
		CreatedAt:    time.Now(),
		ContactEmail: "test@example.com",
	}

	ctx := context.Background()

	// 1. Save New
	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// 2. Find
	found, err := repo.FindByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.Title != "Original Title" {
		t.Errorf("Expected title 'Original Title', got '%s'", found.Title)
	}

	// 3. Update (Save existing)
	l.Title = "Updated Title"
	err = repo.Save(ctx, l)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	foundUpdated, err := repo.FindByID(ctx, "test-1")
	if err != nil {
		t.Fatalf("Failed to find updated: %v", err)
	}
	if foundUpdated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", foundUpdated.Title)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()
	_ = repo.Save(ctx, domain.Listing{ID: "del-me", Title: "Delete Me"})

	err := repo.Delete(ctx, "del-me")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, "del-me")
	if err == nil {
		t.Error("Expected error finding deleted listing, got nil")
	}

	// Delete non-existent
	err = repo.Delete(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error deleting non-existent listing, got nil")
	}
}

func TestHoursOfOperationPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:               "h-1",
		Title:            "Hours Test",
		HoursOfOperation: "Mon-Fri 9am-5pm",
		IsActive:         true,
	}

	_ = repo.Save(ctx, l)

	found, _ := repo.FindByID(ctx, "h-1")
	if found.HoursOfOperation != l.HoursOfOperation {
		t.Errorf("Expected hours %q, got %q", l.HoursOfOperation, found.HoursOfOperation)
	}
}

func TestCategoryErrors_Raw(t *testing.T) {
	t.Parallel()
	db, _ := sql.Open("sqlite", ":memory:")
	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	_ = db.Close()
	ctx := context.Background()

	_, err := repo.GetCategory(ctx, "any")
	if err == nil {
		t.Error("Expected error on closed DB")
	}
}

func TestEnrichmentAttemptedAtPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	now := time.Now().Truncate(time.Second) // SQLite might truncate sub-second precision
	l := domain.Listing{
		ID:                    "enrich-1",
		Title:                 "Enrich Test",
		EnrichmentAttemptedAt: &now,
		IsActive:              true,
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "enrich-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}

	if found.EnrichmentAttemptedAt == nil {
		t.Fatal("Expected EnrichmentAttemptedAt to be set, got nil")
	}

	if !found.EnrichmentAttemptedAt.Equal(now) {
		t.Errorf("Expected EnrichmentAttemptedAt %v, got %v", now, *found.EnrichmentAttemptedAt)
	}
}

