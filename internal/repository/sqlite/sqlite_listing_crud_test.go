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

func TestFindEnrichmentTargets_FiltersAttemptedAt(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)
	eightDaysAgo := now.AddDate(0, 0, -8)

	l1 := domain.Listing{
		ID:         "unenriched-1",
		Title:      "Unenriched 1",
		WebsiteURL: "http://test.com",
		IsActive:   true,
	}
	l2 := domain.Listing{
		ID:                    "unenriched-2",
		Title:                 "Unenriched 2",
		WebsiteURL:            "http://test.com",
		EnrichmentAttemptedAt: &oneDayAgo,
		IsActive:              true,
	}
	l3 := domain.Listing{
		ID:                    "unenriched-3",
		Title:                 "Unenriched 3",
		WebsiteURL:            "http://test.com",
		EnrichmentAttemptedAt: &eightDaysAgo,
		IsActive:              true,
	}

	_ = repo.Save(ctx, l1)
	_ = repo.Save(ctx, l2)
	_ = repo.Save(ctx, l3)

	targets, err := repo.FindEnrichmentTargets(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to find enrichment targets: %v", err)
	}

	targetsMap := make(map[string]bool)
	for _, target := range targets {
		targetsMap[target.ID] = true
	}

	if !targetsMap["unenriched-1"] {
		t.Error("Expected to find unenriched-1 (attempted_at is NULL)")
	}
	if targetsMap["unenriched-2"] {
		t.Error("Did NOT expect to find unenriched-2 (attempted_at is 1 day ago)")
	}
	if !targetsMap["unenriched-3"] {
		t.Error("Expected to find unenriched-3 (attempted_at is 8 days ago)")
	}
}

func TestDeliveryPlatformsPersistence(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:                "dp-1",
		Title:             "Delivery Test",
		DeliveryPlatforms: `["UberEats", "DoorDash"]`,
		IsActive:          true,
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "dp-1")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
	if found.DeliveryPlatforms != l.DeliveryPlatforms {
		t.Errorf("Expected delivery platforms %q, got %q", l.DeliveryPlatforms, found.DeliveryPlatforms)
	}
}

