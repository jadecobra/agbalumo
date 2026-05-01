package sqlite_test

import (
	"context"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

const (
	testUserID    = "user-1"
	testListingID = "listing-1"
)

func TestSaveListing_SaveAndRetrieve(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	// Setup: Save user and listing
	testutil.SaveTestUser(t, repo, testUserID, "user1@example.com", domain.UserRoleUser)
	testutil.SaveTestListing(t, repo, testListingID, "Test Listing")

	// Action: Save listing
	err := repo.SaveListing(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("Failed to save listing: %v", err)
	}

	// Verify: Get saved listings
	saved, err := repo.GetSavedListings(ctx, testUserID)
	if err != nil {
		t.Fatalf("Failed to get saved listings: %v", err)
	}

	if len(saved) != 1 {
		t.Errorf("Expected 1 saved listing, got %d", len(saved))
	} else if saved[0].ListingID != testListingID {
		t.Errorf("Expected listing ID %s, got %s", testListingID, saved[0].ListingID)
	}
}

func TestSaveListing_Unsave(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	// Setup: Save user, listing, and save the listing
	testutil.SaveTestUser(t, repo, testUserID, "user1@example.com", domain.UserRoleUser)
	testutil.SaveTestListing(t, repo, testListingID, "Test Listing")
	_ = repo.SaveListing(ctx, testUserID, testListingID)

	// Action: Unsave listing
	err := repo.UnsaveListing(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("Failed to unsave listing: %v", err)
	}

	// Verify: Check if saved
	saved, err := repo.IsListingSaved(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("Failed to check if listing is saved: %v", err)
	}
	if saved {
		t.Error("Expected listing to be unsaved, but it's still saved")
	}
}

func TestSaveListing_DuplicateSave(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	// Setup: Save user and listing
	testutil.SaveTestUser(t, repo, testUserID, "user1@example.com", domain.UserRoleUser)
	testutil.SaveTestListing(t, repo, testListingID, "Test Listing")

	// Action: Save listing twice
	err := repo.SaveListing(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("First save failed: %v", err)
	}
	err = repo.SaveListing(ctx, testUserID, testListingID)
	if err != nil {
		t.Errorf("Second save (duplicate) failed: %v. Expected no error (INSERT OR IGNORE)", err)
	}

	// Verify: Count should still be 1
	saved, err := repo.GetSavedListingIDs(ctx, testUserID)
	if err != nil {
		t.Fatalf("Failed to get saved IDs: %v", err)
	}
	if len(saved) != 1 {
		t.Errorf("Expected 1 saved ID after duplicate save, got %d", len(saved))
	}
}

func TestSaveListing_IsListingSaved(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	testutil.SaveTestUser(t, repo, testUserID, "user1@example.com", domain.UserRoleUser)
	testutil.SaveTestListing(t, repo, testListingID, "Test Listing")

	// Before saving
	saved, err := repo.IsListingSaved(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("IsListingSaved failed: %v", err)
	}
	if saved {
		t.Error("Expected False, got True")
	}

	// After saving
	_ = repo.SaveListing(ctx, testUserID, testListingID)
	saved, err = repo.IsListingSaved(ctx, testUserID, testListingID)
	if err != nil {
		t.Fatalf("IsListingSaved failed: %v", err)
	}
	if !saved {
		t.Error("Expected True, got False")
	}
}

func TestSaveListing_GetSavedListingIDs(t *testing.T) {
	repo, _ := testutil.SetupTestRepositoryUnique(t)
	defer func() { _ = repo.Close() }()

	ctx := context.Background()

	testutil.SaveTestUser(t, repo, testUserID, "user1@example.com", domain.UserRoleUser)
	listings := []string{"l1", "l2", "l3"}
	for _, id := range listings {
		testutil.SaveTestListing(t, repo, id, "Title "+id)
		_ = repo.SaveListing(ctx, testUserID, id)
	}

	// Verify: Get IDs
	ids, err := repo.GetSavedListingIDs(ctx, testUserID)
	if err != nil {
		t.Fatalf("Failed to get saved IDs: %v", err)
	}

	if len(ids) != 3 {
		t.Errorf("Expected 3 IDs, got %d", len(ids))
	}

	// Should be ordered by created_at DESC
	found := make(map[string]bool)
	for _, id := range ids {
		found[id] = true
	}
	for _, id := range listings {
		if !found[id] {
			t.Errorf("Expected ID %s not found in result", id)
		}
	}
}
