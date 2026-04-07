package sqlite_test

import (
	"github.com/jadecobra/agbalumo/internal/testutil"
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestFindAll_Filtering(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestListingRepository_FTS(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestFindByTitle_CaseInsensitive(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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

func TestFindAllByOwner(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
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
