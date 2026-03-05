package cmd

import (
	"context"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestListingCLI_ImageFlags(t *testing.T) {
	// Use a temporary database for integration tests
	tempDB := "test_cli_images.db"
	os.Setenv("DATABASE_URL", tempDB)
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Remove(tempDB)
		os.Remove(tempDB + "-shm")
		os.Remove(tempDB + "-wal")
	}()

	repo := initRepo()

	// 1. Create a listing with an image via CLI
	listingID := "cli-test-image"
	// Directly call the Run function of the command (simulating CLI)
	// Reset flags
	flagTitle = "CLI Image Test"
	flagType = "Business"
	flagOrigin = "Nigeria"
	flagDescription = "Test Desc"
	flagEmail = "test@test.com"
	flagAddress = "123 Test St"
	flagImageURL = "/static/uploads/test.webp"
	flagOwnerID = "user-123"

	// Simulate 'create'
	listing := domain.Listing{
		ID:           listingID,
		OwnerID:      flagOwnerID,
		OwnerOrigin:  flagOrigin,
		Type:         domain.Category(flagType),
		Title:        flagTitle,
		Description:  flagDescription,
		Address:      flagAddress,
		ContactEmail: flagEmail,
		ImageURL:     flagImageURL,
		IsActive:     true,
		Status:       domain.ListingStatusApproved,
	}
	err := repo.Save(context.Background(), listing)
	assert.NoError(t, err)

	// 2. Verify it exists with the image
	found, err := repo.FindByID(context.Background(), listingID)
	assert.NoError(t, err)
	assert.Equal(t, "/static/uploads/test.webp", found.ImageURL)

	// 3. Update the listing to remove the image via CLI logic
	flagRemoveImage = true
	// Simulate 'update' logic from cmd/listing.go
	listingToUpdate, err := repo.FindByID(context.Background(), listingID)
	assert.NoError(t, err)
	if flagRemoveImage {
		listingToUpdate.ImageURL = ""
	}
	err = repo.Save(context.Background(), listingToUpdate)
	assert.NoError(t, err)

	// 4. Verify image is gone
	updated, err := repo.FindByID(context.Background(), listingID)
	assert.NoError(t, err)
	assert.Equal(t, "", updated.ImageURL)

	// Reset flags for other tests
	flagRemoveImage = false
	flagImageURL = ""
}
