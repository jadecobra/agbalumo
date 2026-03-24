package cmd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestListingCommands_RunCoverage(t *testing.T) {
	tempDB := "test_cli_commands.db"
	_ = os.Setenv("DATABASE_URL", tempDB)
	defer func() {
		_ = os.Unsetenv("DATABASE_URL")
		_ = os.Remove(tempDB)
		_ = os.Remove(tempDB + "-shm")
		_ = os.Remove(tempDB + "-wal")
	}()

	repo := initRepo()

	// Ensure the repo is clean
	listings, _, _ := repo.FindAll(context.Background(), "", "", "", "", false, 100, 0)
	for _, l := range listings {
		_ = repo.Delete(context.Background(), l.ID)
	}

	t.Run("Create", func(t *testing.T) {
		flagTitle = "CLI Test Create"
		flagType = "Service"
		flagOrigin = "Ghana"
		flagDescription = "Create description"
		flagCity = "Accra"
		flagAddress = "123 Create St"
		flagEmail = "create@test.com"
		flagPhone = "123"
		flagWhatsApp = "456"
		flagWebsite = "url"
		flagImageURL = "img"
		flagOwnerID = "owner-1"
		flagDeadline = "2026-12-31"
		flagEventStart = "2026-01-01T10:00"
		flagEventEnd = "2026-01-01T12:00"
		flagJobStart = "2026-02-01T09:00"
		flagSkills = "Go"
		flagApplyURL = "job.com"
		flagCompany = "TestCo"
		flagPayRange = "100k"
		flagText = false // output json

		listingCreateCmd.Run(listingCreateCmd, nil)

		time.Sleep(100 * time.Millisecond) // Give SQLite a moment if needed, though Save is sync
		
		allListings, _, err := repo.FindAll(context.Background(), "", "", "", "", false, 10, 0)
		assert.NoError(t, err)

		var found domain.Listing
		for _, l := range allListings {
			if l.Title == "CLI Test Create" {
				found = l
			}
		}
		assert.NotEmpty(t, found.ID)
		assert.Equal(t, domain.Category("Service"), found.Type)
		assert.Equal(t, "Accra", found.City)
		
		// Run with flagText = true to hit the format print path
		flagTitle = "CLI Test Create Text"
		flagText = true
		listingCreateCmd.Run(listingCreateCmd, nil)
	})

	t.Run("List", func(t *testing.T) {
		// Output text
		flagText = true
		listingListCmd.Run(listingListCmd, nil)

		// Output JSON
		flagText = false
		listingListCmd.Run(listingListCmd, nil)
	})

	t.Run("Get", func(t *testing.T) {
		allListings, _, _ := repo.FindAll(context.Background(), "", "", "", "", false, 10, 0)
		if len(allListings) > 0 {
			targetID := allListings[0].ID
			
			// JSON format
			flagText = false
			listingGetCmd.Run(listingGetCmd, []string{targetID})
			
			// Text format
			flagText = true
			listingGetCmd.Run(listingGetCmd, []string{targetID})
		}
	})

	t.Run("Update", func(t *testing.T) {
		allListings, _, _ := repo.FindAll(context.Background(), "", "", "", "", false, 10, 0)
		if len(allListings) > 0 {
			targetID := allListings[0].ID
			
			// Set flags
			flagTitle = "Updated Title"
			flagCity = "Updated City"
			flagDescription = "Updated description"
			flagAddress = "Updated address"
			flagEmail = "updated@test.com"
			flagPhone = "updated phone"
			flagWhatsApp = "updated wa"
			flagWebsite = "updated site"
			flagImageURL = "updated img"
			flagRemoveImage = true
			flagDeadline = "2027-12-31"
			flagEventStart = "2027-01-01T10:00"
			flagEventEnd = "2027-01-01T12:00"
			flagSkills = "Updated Skills"
			flagJobStart = "2027-02-01T09:00"
			flagApplyURL = "updated.com"
			flagCompany = "UpdatedCo"
			flagPayRange = "200k"
			
			// JSON output
			flagText = false
			listingUpdateCmd.Run(listingUpdateCmd, []string{targetID})
			
			updated, err := repo.FindByID(context.Background(), targetID)
			assert.NoError(t, err)
			assert.Equal(t, "Updated Title", updated.Title)
			assert.Equal(t, "Updated City", updated.City)
			assert.Empty(t, updated.ImageURL)
			
			// Hit the text output path
			flagTitle = "Updated Title 2"
			flagRemoveImage = false
			flagText = true
			listingUpdateCmd.Run(listingUpdateCmd, []string{targetID})
		}
	})

	t.Run("Backfill", func(t *testing.T) {
		// Ensure we don't crash the test runner with os.Exit(1)
		_ = os.Setenv("GOOGLE_MAPS_API_KEY", "dummy_key_for_test")
		defer func() { _ = os.Unsetenv("GOOGLE_MAPS_API_KEY") }()

		// Create a listing missing City but has address
		l := domain.Listing{
			ID: "backfill-test-1",
			Title: "Backfill Test",
			Address: "123 Test Ave, Fake City",
			// City intentionally blank
		}
		_ = repo.Save(context.Background(), l)

		// Run backfill
		listingBackfillCitiesCmd.Run(listingBackfillCitiesCmd, nil)
		
		// The dummy key will fail geocoding, which is fine, we just want coverage of the code paths
		// checking if errorCount increments properly and branch coverage.
	})

	t.Run("Delete", func(t *testing.T) {
		allListings, _, _ := repo.FindAll(context.Background(), "", "", "", "", false, 10, 0)
		for _, l := range allListings {
			listingDeleteCmd.Run(listingDeleteCmd, []string{l.ID})
			
			_, err := repo.FindByID(context.Background(), l.ID)
			assert.Error(t, err)
		}
	})
}
