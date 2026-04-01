package cmd

import (
	"os"
	"path/filepath"
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
		_ = os.Unsetenv("DATABASE_URL")
		url := getDatabaseURL()
		if url != ".tester/data/agbalumo.db" {
			t.Errorf("getDatabaseURL() = %v, want .tester/data/agbalumo.db", url)
		}
	})

	t.Run("from env", func(t *testing.T) {
		_ = os.Setenv("DATABASE_URL", "custom.db")
		defer func() { _ = os.Unsetenv("DATABASE_URL") }()
		url := getDatabaseURL()
		if url != "custom.db" {
			t.Errorf("getDatabaseURL() = %v, want custom.db", url)
		}
	})
}

func TestInitRepo(t *testing.T) {
	// Isolate database for CI parity
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_init.db")
	_ = os.Setenv("DATABASE_URL", dbPath)
	defer func() { _ = os.Unsetenv("DATABASE_URL") }()

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

	printListing(rootCmd, listing)
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

	printListing(rootCmd, listing)
}

func TestPrintListingSummary(t *testing.T) {
	listing := domain.Listing{
		ID:     "summary-test-id-789",
		Title:  "Summary Test Business",
		Type:   domain.Business,
		City:   "Ibadan",
		Status: domain.ListingStatusApproved,
	}

	printListingSummary(rootCmd, listing)
}

func TestPrintListingAllFields(t *testing.T) {
	listing := domain.Listing{
		ID:               "full-test-id",
		Title:            "Full Test Job",
		Type:             domain.Job,
		OwnerOrigin:      "Ghana",
		Status:           domain.ListingStatusApproved,
		Featured:         true,
		Description:      "A very detailed job description",
		City:             "Accra",
		Address:          "456 Job Lane",
		HoursOfOperation: "Mon-Fri 8am-5pm",
		ContactEmail:     "job@test.com",
		ContactPhone:     "+123456789",
		ContactWhatsApp:  "+987654321",
		WebsiteURL:       "https://job.test.com",
		ImageURL:         "https://job.test.com/logo.png",
		CreatedAt:        time.Now(),
		Deadline:         time.Now().Add(30 * 24 * time.Hour),
		EventStart:       time.Now().Add(24 * time.Hour),
		EventEnd:         time.Now().Add(48 * time.Hour),
		Skills:           "Go, CLI, Testing",
		JobStartDate:     time.Now().Add(60 * 24 * time.Hour),
		JobApplyURL:      "https://job.test.com/apply",
		Company:          "CLI Corp",
		PayRange:         "$120k-$150k",
	}

	printListing(rootCmd, listing)
}
