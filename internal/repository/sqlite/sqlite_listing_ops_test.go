package sqlite_test

import (
	"context"
	"database/sql"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
)

func TestGetCounts(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestExpireListings(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestGetPendingClaimRequests(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestGetClaimRequestByUserAndListing(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestGetMetrics(t *testing.T) {
	t.Parallel()
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestMigrationBackfillCity(t *testing.T) {
	t.Parallel()
	// 1. Create a raw DB and insert a row with NULL/empty city
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "backfill.db")
	db, err := sql.Open("sqlite", dbPath)
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
