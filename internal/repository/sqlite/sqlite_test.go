package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	dsn := dbPath + "?_time_format=sqlite"
	repo, err := sqlite.NewSQLiteRepository(dsn)
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
	defer func() { _ = db.Close() }()

	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	if repo == nil {
		t.Fatal("Expected repository, got nil")
	}

	// Verify we can use it
	ctx := context.Background()
	_, _, err = repo.FindAll(ctx, "All", "", "", "", false, 20, 0)
	if err == nil {
		t.Error("Expected error due to missing tables, got nil")
	}
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

func TestFindAll_Filtering(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed Data
	_ = repo.Save(ctx, domain.Listing{ID: "1", Title: "Jollof Rice", Type: "Business", Status: domain.ListingStatusApproved, IsActive: true, CreatedAt: time.Now()})
	_ = repo.Save(ctx, domain.Listing{ID: "2", Title: "Hair Braiding", Type: "Service", Status: domain.ListingStatusPending, IsActive: true, CreatedAt: time.Now()})
	_ = repo.Save(ctx, domain.Listing{ID: "3", Title: "Deleted Item", Type: "Product", Status: domain.ListingStatusRejected, IsActive: false, CreatedAt: time.Now()})

	// 1. Find All Active (Default for Public) - should only return Approved
	res, _, err := repo.FindAll(ctx, "", "", "", "", false, 20, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected 1 approved listing, got %d", len(res))
	}

	// 2. Filter by Type
	res, _, _ = repo.FindAll(ctx, "Business", "", "", "", false, 20, 0)
	if len(res) != 1 || res[0].Title != "Jollof Rice" {
		t.Errorf("Type filtering failed")
	}

	// 3. Search text (public) - should NOT find Pending "Braiding"
	res, _, _ = repo.FindAll(ctx, "", "Braiding", "", "", false, 20, 0)
	if len(res) != 0 {
		t.Errorf("Expected 0 results for public search of pending listing, got %d", len(res))
	}
	
	// 3b. Search text (admin) - should find Pending "Braiding"
	res, _, _ = repo.FindAll(ctx, "", "Braiding", "", "", true, 20, 0)
	if len(res) != 1 {
		t.Errorf("Expected 1 result for admin search of pending listing, got %d", len(res))
	}

	// 4. Include Inactive (Admin view)
	res, _, _ = repo.FindAll(ctx, "", "", "", "", true, 20, 0)
	if len(res) != 3 {
		t.Errorf("Expected 3 listings including inactive, got %d", len(res))
	}
}

func TestFindAll_Sorting_Featured(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed Data
	_ = repo.Save(ctx, domain.Listing{ID: "1", Title: "A-Normal", Status: domain.ListingStatusApproved, IsActive: true, Featured: false, CreatedAt: time.Now()})
	_ = repo.Save(ctx, domain.Listing{ID: "2", Title: "B-Featured", Status: domain.ListingStatusApproved, IsActive: true, Featured: true, CreatedAt: time.Now().Add(time.Second)})
	_ = repo.Save(ctx, domain.Listing{ID: "3", Title: "C-Normal", Status: domain.ListingStatusApproved, IsActive: true, Featured: false, CreatedAt: time.Now().Add(2 * time.Second)})

	// Test sort by featured DESC (Featured should be first)
	res, _, err := repo.FindAll(ctx, "", "", "featured", "DESC", true, 10, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(res) != 3 {
		t.Fatalf("Expected 3 listings, got %d", len(res))
	}

	if res[0].ID != "2" {
		t.Errorf("Expected first listing to be featured (ID 2), got ID %s", res[0].ID)
	}

	// Test sort by featured ASC (Featured should be last)
	resAsc, _, errAsc := repo.FindAll(ctx, "", "", "featured", "ASC", true, 10, 0)
	if errAsc != nil {
		t.Fatalf("FindAll failed: %v", errAsc)
	}
	
	if len(resAsc) != 3 {
		t.Fatalf("Expected 3 listings, got %d", len(resAsc))
	}

	if resAsc[2].ID != "2" {
		t.Errorf("Expected last listing to be featured (ID 2), got ID %s", resAsc[2].ID)
	}
}

func TestGetCounts(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{ID: "1", Type: "Business", IsActive: true})
	_ = repo.Save(ctx, domain.Listing{ID: "2", Type: "Business", IsActive: true})
	_ = repo.Save(ctx, domain.Listing{ID: "3", Type: "Service", IsActive: true})
	_ = repo.Save(ctx, domain.Listing{ID: "4", Type: "Business", IsActive: false}) // Inactive

	counts, err := repo.GetCounts(ctx)
	if err != nil {
		t.Fatalf("GetCounts failed: %v", err)
	}

	if counts["Business"] != 2 {
		t.Errorf("Expected 2 active Business, got %d", counts["Business"])
	}
	if counts["Service"] != 1 {
		t.Errorf("Expected 1 active Service, got %d", counts["Service"])
	}
}

func TestDelete(t *testing.T) {
	repo, _ := newTestRepo(t)
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

func TestExpireListings(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	now := time.Now().UTC()

	// 1. Request past deadline
	_ = repo.Save(ctx, domain.Listing{ID: "exp-1", Type: domain.Request, Deadline: now.Add(-1 * time.Hour), IsActive: true})
	// 2. Event past end time
	_ = repo.Save(ctx, domain.Listing{ID: "exp-2", Type: domain.Event, EventEnd: now.Add(-1 * time.Hour), IsActive: true})
	// 3. Active listing (future deadline)
	_ = repo.Save(ctx, domain.Listing{ID: "act-1", Type: domain.Request, Deadline: now.Add(1 * time.Hour), IsActive: true})

	expired, err := repo.ExpireListings(ctx)
	if err != nil {
		t.Fatalf("ExpireListings failed: %v", err)
	}

	if expired != 2 {
		t.Errorf("Expected 2 expired listings, got %d", expired)
	}

	l, _ := repo.FindByID(ctx, "act-1")
	if !l.IsActive {
		t.Error("Expected future listing to remain active")
	}
}

func TestHoursOfOperationPersistence(t *testing.T) {
	repo, _ := newTestRepo(t)
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

func TestFindAllByOwner(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	_ = repo.Save(ctx, domain.Listing{ID: "1", OwnerID: "user-1", Title: "L1"})
	_ = repo.Save(ctx, domain.Listing{ID: "2", OwnerID: "user-1", Title: "L2"})
	_ = repo.Save(ctx, domain.Listing{ID: "3", OwnerID: "user-2", Title: "L3"})

	res, _, err := repo.FindAllByOwner(ctx, "user-1", 10, 0)
	if err != nil {
		t.Fatalf("FindAllByOwner failed: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("Expected 2 listings for user-1, got %d", len(res))
	}
}

func TestGetPendingClaimRequests(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Save a listing so we can reference it
	_ = repo.Save(ctx, domain.Listing{ID: "l1", Title: "Unclaimed Business", Type: domain.Business, IsActive: true, CreatedAt: time.Now()})

	// Save a pending claim request
	_ = repo.SaveClaimRequest(ctx, domain.ClaimRequest{
		ID: "cr1", ListingID: "l1", UserID: "u1",
		Status: domain.ClaimStatusPending, CreatedAt: time.Now(),
	})

	// Also save an approved one (should not appear)
	_ = repo.SaveClaimRequest(ctx, domain.ClaimRequest{
		ID: "cr2", ListingID: "l1", UserID: "u2",
		Status: domain.ClaimStatusApproved, CreatedAt: time.Now(),
	})

	claims, err := repo.GetPendingClaimRequests(ctx)
	if err != nil {
		t.Fatalf("GetPendingClaimRequests failed: %v", err)
	}

	if len(claims) != 1 {
		t.Errorf("Expected 1 pending claim, got %d", len(claims))
	}
	if claims[0].ID != "cr1" {
		t.Errorf("Expected claim cr1, got %s", claims[0].ID)
	}
}

func TestUpdateClaimRequestStatus_Approve(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Create listing and a user
	_ = repo.Save(ctx, domain.Listing{ID: "l1", Title: "Biz", Type: domain.Business, IsActive: true, CreatedAt: time.Now()})
	_ = repo.SaveUser(ctx, domain.User{ID: "u1", GoogleID: "g1", Email: "u1@x.com", CreatedAt: time.Now()})

	// Create pending claim
	_ = repo.SaveClaimRequest(ctx, domain.ClaimRequest{
		ID: "cr1", ListingID: "l1", UserID: "u1",
		Status: domain.ClaimStatusPending, CreatedAt: time.Now(),
	})

	// Approve it
	err := repo.UpdateClaimRequestStatus(ctx, "cr1", domain.ClaimStatusApproved)
	if err != nil {
		t.Fatalf("UpdateClaimRequestStatus failed: %v", err)
	}

	// Verify status updated
	claims, _ := repo.GetPendingClaimRequests(ctx)
	if len(claims) != 0 {
		t.Error("Expected 0 pending claims after approval")
	}

	// Verify owner transferred
	l, _ := repo.FindByID(ctx, "l1")
	if l.OwnerID != "u1" {
		t.Errorf("Expected owner to be u1, got %s", l.OwnerID)
	}
}


func TestListingRepository_FTS(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Use trigram search (needs at least 3 chars usually)
	_ = repo.Save(ctx, domain.Listing{ID: "1", Title: "Jollof Rice", City: "Houston", IsActive: true})
	_ = repo.Save(ctx, domain.Listing{ID: "2", Title: "Egusi Soup", City: "Dallas", IsActive: true})

	tests := []struct {
		query string
		want  string
	}{
		{"Jollof", "1"},
		{"Soup", "2"},
		{"Houston", "1"},
		{"Egusi", "2"},
	}

	for _, tt := range tests {
		res, _, _ := repo.FindAll(ctx, "", tt.query, "", "", false, 10, 0)
		if len(res) != 1 || res[0].ID != tt.want {
			t.Errorf("FTS query %q failed: got %d results, want ID %s", tt.query, len(res), tt.want)
		}
	}
}

func TestSQLiteRepository_InitializationErrors(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file")
	_ = os.WriteFile(filePath, []byte("content"), 0644)

	// Try to use that file as a directory for the DB
	_, err := sqlite.NewSQLiteRepository(filePath)
	if err == nil {
		t.Error("Expected error initializing repo on a plain file without write access as DB, got nil")
	}
}

func TestGetMetrics(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	now := time.Now().UTC()

	// 1. Listings Growth
	// Yesterday
	_ = repo.Save(ctx, domain.Listing{ID: "g1", Title: "L1", Type: domain.Business, IsActive: true, CreatedAt: now.Add(-24 * time.Hour)})
	// Today
	_ = repo.Save(ctx, domain.Listing{ID: "g2", Title: "L2", Type: domain.Business, IsActive: true, CreatedAt: now})
	_ = repo.Save(ctx, domain.Listing{ID: "g3", Title: "L3", Type: domain.Business, IsActive: true, CreatedAt: now})
	// Old (outside 30 days)
	_ = repo.Save(ctx, domain.Listing{ID: "g4", Title: "Old", Type: domain.Business, IsActive: true, CreatedAt: now.Add(-31 * 24 * time.Hour)})

	metrics, err := repo.GetListingGrowth(ctx)
	if err != nil {
		t.Fatalf("GetListingGrowth failed: %v", err)
	}
	// Should have 2 entries (yesterday and today)
	if len(metrics) != 2 {
		t.Errorf("Expected 2 days of metrics, got %d", len(metrics))
	}
}


func TestCategoryErrors_Raw(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := sqlite.NewSQLiteRepositoryFromDB(db)
	_ = db.Close()
	ctx := context.Background()

	_, err := repo.GetCategory(ctx, "any")
	if err == nil {
		t.Error("Expected error on closed DB")
	}
}

func TestFindByTitle_CaseInsensitive(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	// Seed data
	_ = repo.Save(ctx, domain.Listing{ID: "1", Title: "Unique Title", Type: domain.Business, IsActive: true, CreatedAt: time.Now()})
	_ = repo.Save(ctx, domain.Listing{ID: "2", Title: "Duplicate Title", Type: domain.Service, IsActive: true, CreatedAt: time.Now()})
	_ = repo.Save(ctx, domain.Listing{ID: "3", Title: "Duplicate Title", Type: domain.Product, IsActive: true, CreatedAt: time.Now()})

	// Test unique title
	res, _ := repo.FindByTitle(ctx, "Unique Title")
	if len(res) != 1 {
		t.Errorf("Expected 1 result for Unique Title, got %d", len(res))
	}

	// Test duplicates
	res, _ = repo.FindByTitle(ctx, "Duplicate Title")
	if len(res) != 2 {
		t.Errorf("Expected 2 results for Duplicate Title, got %d", len(res))
	}
}

func TestTitleExists(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	title := "Unique Listing"
	_ = repo.Save(ctx, domain.Listing{ID: "1", Title: title, IsActive: true})

	exists, err := repo.TitleExists(ctx, title)
	if err != nil {
		t.Fatalf("TitleExists failed: %v", err)
	}
	if !exists {
		t.Error("Expected title to exist")
	}

	exists, err = repo.TitleExists(ctx, "Non-existent Title")
	if err != nil {
		t.Fatalf("TitleExists failed: %v", err)
	}
	if exists {
		t.Error("Expected title to not exist")
	}
}

func TestGetClaimRequestByUserAndListing(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	userID := "u1"
	listingID := "l1"

	// 1. Not found
	_, err := repo.GetClaimRequestByUserAndListing(ctx, userID, listingID)
	if err == nil {
		t.Error("Expected error when claim request not found, got nil")
	}

	// 2. Found
	req := domain.ClaimRequest{
		ID:        "cr1",
		UserID:    userID,
		ListingID: listingID,
		Status:    domain.ClaimStatusPending,
		CreatedAt: time.Now(),
	}
	_ = repo.SaveClaimRequest(ctx, req)

	found, err := repo.GetClaimRequestByUserAndListing(ctx, userID, listingID)
	if err != nil {
		t.Fatalf("GetClaimRequestByUserAndListing failed: %v", err)
	}
	if found.ID != req.ID {
		t.Errorf("Expected ID %s, got %s", req.ID, found.ID)
	}
}

func BenchmarkSQLiteRepository_FindAll(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")
	repo, _ := sqlite.NewSQLiteRepository(dbPath + "?_time_format=sqlite")
	ctx := context.Background()

	// Seed 100 listings
	for i := 0; i < 100; i++ {
		_ = repo.Save(ctx, domain.Listing{
			ID:          fmt.Sprintf("l%d", i),
			Title:       fmt.Sprintf("Listing %d", i),
			Type:        domain.Business,
			OwnerOrigin: "Nigeria",
			Address:     "123 St",
			Description: "Desc",
			IsActive:    true,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.FindAll(ctx, "", "", "", "", true, 50, 0)
	}
}

func BenchmarkSQLiteRepository_FindByTitle(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench_title.db")
	repo, _ := sqlite.NewSQLiteRepository(dbPath + "?_time_format=sqlite")
	ctx := context.Background()

	// Seed 100 listings
	for i := 0; i < 100; i++ {
		_ = repo.Save(ctx, domain.Listing{
			ID:          fmt.Sprintf("l%d", i),
			Title:       fmt.Sprintf("Listing %d", i),
			Type:        domain.Business,
			OwnerOrigin: "Nigeria",
			Address:     "123 St",
			Description: "Desc",
			IsActive:    true,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByTitle(ctx, "Listing 50")
	}
}

func TestMigrationBackfillCity(t *testing.T) {
	// 1. Create a raw DB and insert a row with NULL/empty city
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "backfill.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open raw db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE listings (
		id TEXT PRIMARY KEY,
		owner_id TEXT,
		title TEXT,
		description TEXT,
		type TEXT,
		owner_origin TEXT,
		city TEXT,
		address TEXT,
		hours_of_operation TEXT DEFAULT '',
		is_active BOOLEAN,
		created_at DATETIME,
		image_url TEXT,
		contact_email TEXT,
		contact_phone TEXT,
		contact_whatsapp TEXT,
		website_url TEXT,
		deadline DATETIME,
		skills TEXT,
		job_start_date DATETIME,
		job_apply_url TEXT,
		company TEXT,
		pay_range TEXT,
		status TEXT DEFAULT 'Approved',
		featured BOOLEAN DEFAULT 0
	);`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO listings (id, owner_id, owner_origin, type, title, description, city, address, is_active, status, contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at) 
		VALUES ('1', 'owner1', 'Nigeria', 'Business', 'Empty City', 'Desc', '', '123 St', 1, 'Approved', '', '', '', '', '', CURRENT_TIMESTAMP);`)
	if err != nil {
		t.Fatalf("Failed to insert empty city: %v", err)
	}
	_, err = db.Exec(`INSERT INTO listings (id, owner_id, owner_origin, type, title, description, city, address, is_active, status, contact_email, contact_phone, contact_whatsapp, website_url, image_url, created_at) 
		VALUES ('2', 'owner1', 'Nigeria', 'Business', 'NULL City', 'Desc', NULL, '123 St', 1, 'Approved', '', '', '', '', '', CURRENT_TIMESTAMP);`)
	if err != nil {
		t.Fatalf("Failed to insert NULL city: %v", err)
	}
	_ = db.Close()

	// 2. Initialize repository (triggers migration)
	repo, err := sqlite.NewSQLiteRepository(dbPath + "?_time_format=sqlite")
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	// 3. Verify backfill
	ctx := context.Background()
	l1, err1 := repo.FindByID(ctx, "1")
	if err1 != nil {
		t.Fatalf("L1: FindByID failed: %v", err1)
	}
	if l1.City != "" {
		t.Errorf("Expected city '' for ID 1, got %q", l1.City)
	}

	l2, err2 := repo.FindByID(ctx, "2")
	if err2 != nil {
		t.Fatalf("L2: FindByID failed: %v", err2)
	}
	if l2.City != "" {
		t.Errorf("Expected city '' for ID 2, got %q", l2.City)
	}
}
