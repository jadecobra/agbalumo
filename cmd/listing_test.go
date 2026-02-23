package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

var testDeadline = time.Now().Add(24 * time.Hour)

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("generateID() should return non-empty string")
	}

	if id1 == id2 {
		t.Error("generateID() should return unique IDs")
	}
}

func TestGetDatabaseURL(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		url := getDatabaseURL()
		if url != "agbalumo.db" {
			t.Errorf("getDatabaseURL() = %v, want agbalumo.db", url)
		}
	})

	t.Run("from env", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "custom.db")
		defer os.Unsetenv("DATABASE_URL")
		url := getDatabaseURL()
		if url != "custom.db" {
			t.Errorf("getDatabaseURL() = %v, want custom.db", url)
		}
	})
}

func TestInitRepo(t *testing.T) {
	repo := initRepo()
	if repo == nil {
		t.Error("initRepo() should return non-nil repo")
	}
}

func TestPrintListing(t *testing.T) {
	listing := domain.Listing{
		ID:              "test-id-123",
		OwnerID:         "owner-456",
		OwnerOrigin:     "Nigeria",
		Type:            domain.Business,
		Title:           "Test Business",
		Description:     "A test business description",
		City:            "Lagos",
		Address:         "123 Test Street",
		ContactEmail:    "test@example.com",
		ContactPhone:    "+2341234567890",
		ContactWhatsApp: "+2340987654321",
		WebsiteURL:      "https://test.com",
		Status:          domain.ListingStatusApproved,
		Featured:        true,
		IsActive:        true,
	}

	printListing(listing)
}

func TestPrintListingWithDeadline(t *testing.T) {
	listing := domain.Listing{
		ID:       "test-id-456",
		Title:    "Test Event",
		Type:     domain.Event,
		City:     "Abuja",
		Status:   domain.ListingStatusPending,
		Deadline: testDeadline,
	}

	printListing(listing)
}

func TestPrintListingSummary(t *testing.T) {
	listing := domain.Listing{
		ID:     "summary-test-id-789",
		Title:  "Summary Test Business",
		Type:   domain.Business,
		City:   "Ibadan",
		Status: domain.ListingStatusApproved,
	}

	printListingSummary(listing)
}
