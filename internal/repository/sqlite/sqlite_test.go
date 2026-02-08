package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func newTestRepo(t *testing.T) (*sqlite.SQLiteRepository, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	return repo, dbPath
}

func TestSaveAndFindByID(t *testing.T) {
	repo, _ := newTestRepo(t)
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
	if err := repo.Save(ctx, l); err != nil {
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

func TestFindAll_Filtering(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed Data
	repo.Save(ctx, domain.Listing{ID: "1", Title: "Jollof Rice", Type: "Business", IsActive: true, CreatedAt: time.Now()})
	repo.Save(ctx, domain.Listing{ID: "2", Title: "Hair Braiding", Type: "Service", IsActive: true, CreatedAt: time.Now()})
	repo.Save(ctx, domain.Listing{ID: "3", Title: "Deleted Item", Type: "Product", IsActive: false, CreatedAt: time.Now()})

	// 1. Find All Active (Default for Public)
	// Query: empty, Type: empty, IncludeInactive: false
	allActive, err := repo.FindAll(ctx, "", "", false)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(allActive) != 2 {
		t.Errorf("Expected 2 active listings, got %d", len(allActive))
	}

	// 2. Find With Inactive (Admin View)
	allAdmin, err := repo.FindAll(ctx, "", "", true)
	if err != nil {
		t.Fatalf("FindAll Admin failed: %v", err)
	}
	if len(allAdmin) != 3 {
		t.Errorf("Expected 3 listings (incl inactive), got %d", len(allAdmin))
	}

	// 3. Filter by Type
	services, err := repo.FindAll(ctx, "Service", "", false)
	if err != nil {
		t.Fatalf("FindAll Type failed: %v", err)
	}
	if len(services) != 1 || services[0].Title != "Hair Braiding" {
		t.Errorf("Expected 1 service 'Hair Braiding', got %v", services)
	}

	// 4. Search Query (LIKE)
	searchRes, err := repo.FindAll(ctx, "", "Rice", false)
	if err != nil {
		t.Fatalf("FindAll Search failed: %v", err)
	}
	if len(searchRes) != 1 || searchRes[0].Title != "Jollof Rice" {
		t.Errorf("Expected 1 result 'Jollof Rice', got %v", searchRes)
	}
}

func TestGetCounts(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed Data
	seedListings := []domain.Listing{
		{ID: "1", Title: "Food 1", Type: domain.Food, IsActive: true, CreatedAt: time.Now(), OwnerOrigin: "Nigeria", ContactEmail: "a@b.com"},
		{ID: "2", Title: "Food 2", Type: domain.Food, IsActive: true, CreatedAt: time.Now(), OwnerOrigin: "Nigeria", ContactEmail: "a@b.com"},
		{ID: "3", Title: "Business 1", Type: domain.Business, IsActive: true, CreatedAt: time.Now(), OwnerOrigin: "Nigeria", ContactEmail: "a@b.com"},
		{ID: "4", Title: "Inactive Service", Type: domain.Service, IsActive: false, CreatedAt: time.Now(), OwnerOrigin: "Nigeria", ContactEmail: "a@b.com"},
	}

	for _, l := range seedListings {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to seed listing: %v", err)
		}
	}

	counts, err := repo.GetCounts(ctx)
	if err != nil {
		t.Fatalf("GetCounts failed: %v", err)
	}

	if counts[domain.Food] != 2 {
		t.Errorf("Expected 2 Food, got %d", counts[domain.Food])
	}
	if counts[domain.Business] != 1 {
		t.Errorf("Expected 1 Business, got %d", counts[domain.Business])
	}
	if counts[domain.Service] != 0 {
		t.Errorf("Expected 0 Service (inactive), got %d", counts[domain.Service])
	}
}

func TestExpireListings(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// 1. Seed Data
	now := time.Now().UTC()
	expiredRequest := domain.Listing{
		ID:           "req-exp",
		Type:         domain.Request,
		Title:        "Expired Request",
		Deadline:     now.Add(-1 * time.Hour),
		IsActive:     true,
		CreatedAt:    now,
		OwnerOrigin:  "Ghana",
		ContactEmail: "test@test.com",
	}
	activeRequest := domain.Listing{
		ID:           "req-act",
		Type:         domain.Request,
		Title:        "Active Request",
		Deadline:     now.Add(1 * time.Hour),
		IsActive:     true,
		CreatedAt:    now,
		OwnerOrigin:  "Ghana",
		ContactEmail: "test@test.com",
	}
	expiredEvent := domain.Listing{
		ID:           "evt-exp",
		Type:         domain.Event,
		Title:        "Expired Event",
		EventEnd:     now.Add(-1 * time.Hour),
		IsActive:     true,
		CreatedAt:    now,
		OwnerOrigin:  "Ghana",
		ContactEmail: "test@test.com",
	}
	activeEvent := domain.Listing{
		ID:           "evt-act",
		Type:         domain.Event,
		Title:        "Active Event",
		EventEnd:     now.Add(1 * time.Hour),
		IsActive:     true,
		CreatedAt:    now,
		OwnerOrigin:  "Ghana",
		ContactEmail: "test@test.com",
	}

	for _, l := range []domain.Listing{expiredRequest, activeRequest, expiredEvent, activeEvent} {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to save %s: %v", l.Title, err)
		}
	}

	// 2. Run Expiration
	count, err := repo.ExpireListings(ctx)
	if err != nil {
		t.Fatalf("ExpireListings failed: %v", err)
	}

	// 3. Assertions
	if count != 2 {
		t.Errorf("Expected 2 listing to expire, got %d", count)
	}

	// Verify specific items
	l, _ := repo.FindByID(ctx, "req-exp")
	if l.IsActive {
		t.Error("Expired Request should be inactive")
	}

	l, _ = repo.FindByID(ctx, "req-act")
	if !l.IsActive {
		t.Error("Active Request should be active")
	}

	l, _ = repo.FindByID(ctx, "evt-exp")
	if l.IsActive {
		t.Error("Expired Event should be inactive")
	}
}

func TestHoursOfOperationPersistence(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	l := domain.Listing{
		ID:               "hours-persist",
		Type:             domain.Business,
		Title:            "Open Late",
		HoursOfOperation: "Mon-Sat 10AM-10PM",
		IsActive:         true,
		CreatedAt:        time.Now(),
		OwnerOrigin:      "Nigeria",
		ContactEmail:     "late@example.com",
		Address:          "123 Late St",
	}

	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	found, err := repo.FindByID(ctx, "hours-persist")
	if err != nil {
		t.Fatalf("Failed to find: %v", err)
	}

	if found.HoursOfOperation != "Mon-Sat 10AM-10PM" {
		t.Errorf("Expected HoursOfOperation 'Mon-Sat 10AM-10PM', got '%s'", found.HoursOfOperation)
	}
}


