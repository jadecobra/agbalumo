package sqlite_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
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

func TestNewSQLiteRepositoryFromDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	// Do not close DB here if repo takes ownership, or close it after test?
	// The repo doesn't close DB on its own usually unless Close() is called?
	// Let's defer close.
	defer db.Close()

	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	if repo == nil {
		t.Fatal("Expected repository, got nil")
	}

	// Verify we can use it
	ctx := context.Background()
	// Should fail because no tables
	_, err = repo.FindAll(ctx, "All", "", false) // Fixed signature call
	if err == nil {
		t.Error("Expected error due to missing tables, got nil")
	}
}

func TestSaveAndFindByID(t *testing.T) {
	repo, _ := newTestRepo(t) // Modified call to ignore 2nd return value
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

func TestDelete(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// 1. Seed
	l := domain.Listing{
		ID:        "to-delete",
		Title:     "Delete Me",
		Type:      domain.Business,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	if err := repo.Save(ctx, l); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// 2. Delete
	if err := repo.Delete(ctx, "to-delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 3. Verify Deleted
	_, err := repo.FindByID(ctx, "to-delete")
	if err == nil {
		t.Error("Expected error finding deleted listing, got nil")
	}

	// 4. Delete Non-Existent
	if err := repo.Delete(ctx, "non-existent"); err == nil {
		t.Error("Expected error deleting non-existent listing, got nil")
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
	// Add Job Expiration Test
	expiredJob := domain.Listing{
		ID:           "job-exp",
		Type:         domain.Job,
		Title:        "Expired Job",
		JobStartDate: now.AddDate(0, 0, -91), // 91 days ago (limit is 90)
		IsActive:     true,
		CreatedAt:    now,
		OwnerOrigin:  "Ghana",
		ContactEmail: "test@test.com",
	}

	for _, l := range []domain.Listing{expiredRequest, activeRequest, expiredEvent, activeEvent, expiredJob} {
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
	// Request + Event + Job = 3
	if count != 3 {
		t.Errorf("Expected 3 listings to expire, got %d", count)
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

	l, _ = repo.FindByID(ctx, "job-exp")
	if l.IsActive {
		t.Error("Expired Job should be inactive")
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

func TestFindAllByOwner(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed
	user1Listings := []domain.Listing{
		{ID: "u1-1", Title: "U1 First", OwnerID: "user1", IsActive: true, CreatedAt: time.Now()},
		{ID: "u1-2", Title: "U1 Second", OwnerID: "user1", IsActive: false, CreatedAt: time.Now()}, // Should include inactive
	}
	user2Listing := domain.Listing{ID: "u2-1", Title: "U2 Only", OwnerID: "user2", IsActive: true, CreatedAt: time.Now()}

	for _, l := range append(user1Listings, user2Listing) {
		if err := repo.Save(ctx, l); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	// Test
	results, err := repo.FindAllByOwner(ctx, "user1")
	if err != nil {
		t.Fatalf("FindAllByOwner failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 listings for user1, got %d", len(results))
	}

	// Verify order or content if needed, but count is good start
	foundTitles := make(map[string]bool)
	for _, l := range results {
		foundTitles[l.Title] = true
	}
	if !foundTitles["U1 First"] || !foundTitles["U1 Second"] {
		t.Errorf("Expected listings not found. Got: %v", results)
	}

	// Test Empty
	resultsEmpty, err := repo.FindAllByOwner(ctx, "non-existent")
	if err != nil {
		t.Fatalf("FindAllByOwner empty failed: %v", err)
	}
	if len(resultsEmpty) != 0 {
		t.Errorf("Expected 0 results for non-existent user, got %d", len(resultsEmpty))
	}
}

func TestGetPendingListings(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed mixed status
	repo.Save(ctx, domain.Listing{ID: "1", Title: "Approved", Status: domain.ListingStatusApproved, CreatedAt: time.Now()})
	repo.Save(ctx, domain.Listing{ID: "2", Title: "Pending 1", Status: domain.ListingStatusPending, CreatedAt: time.Now()})
	repo.Save(ctx, domain.Listing{ID: "3", Title: "Rejected", Status: domain.ListingStatusRejected, CreatedAt: time.Now()})
	repo.Save(ctx, domain.Listing{ID: "4", Title: "Pending 2", Status: domain.ListingStatusPending, CreatedAt: time.Now().Add(time.Hour)})

	pending, err := repo.GetPendingListings(ctx)
	if err != nil {
		t.Fatalf("GetPendingListings failed: %v", err)
	}

	if len(pending) != 2 {
		t.Errorf("Expected 2 pending listings, got %d", len(pending))
	}

	// Verify order (Oldest first as per sqlite impl "ORDER BY created_at ASC")
	// Although the SQL says ASC, let's verify.
	// Pending 1 created now, Pending 2 created now+1h.
	// Pending 1 is older (smaller time), so it should be first.
	if pending[0].Title != "Pending 1" {
		t.Errorf("Expected first pending to be 'Pending 1', got '%s'", pending[0].Title)
	}
}

func TestGetUserCount(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Should be 0 initially
	c, err := repo.GetUserCount(ctx)
	if err != nil {
		t.Fatalf("Failed to get count: %v", err)
	}
	if c != 0 {
		t.Errorf("Expected 0 users, got %d", c)
	}

	// Add Users
	repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Email: "e1", CreatedAt: time.Now()})
	repo.SaveUser(ctx, domain.User{ID: "u2", GoogleID: "g2", Email: "e2", CreatedAt: time.Now()})

	c, err = repo.GetUserCount(ctx)
	if err != nil {
		t.Fatalf("Failed to get count after add: %v", err)
	}
	if c != 2 {
		t.Errorf("Expected 2 users, got %d", c)
	}
}
