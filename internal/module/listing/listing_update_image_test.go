package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleUpdate_ImageRemoval(t *testing.T) {

	t.Parallel()

	// Existing listing with an image
	existingListing := domain.Listing{
		ID:           "listing-123",
		OwnerID:      "user1",
		Title:        "Old Title",
		ImageURL:     "/static/uploads/listing-123.webp",
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		OwnerOrigin:  "Nigeria",
		ContactEmail: "owner@example.com",
		Type:         domain.Business,
		City:         "Lagos",
		Address:      "123 Street",
	}
	mockImageService := &testutil.MockImageService{}
	// Expect DeleteImage to be called
	mockImageService.On("DeleteImage", testifyMock.Anything, existingListing.ImageURL).Return(nil)
	// UploadImage might be called with nil if no image is uploaded
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return("", nil).Maybe()

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.ImageSvc = mockImageService
	_ = app.DB.Save(context.Background(), existingListing)

	h := listmod.NewListingHandler(app)

	// Body with remove_image=true and required fields
	body := "title=New+Title&remove_image=true&owner_origin=Nigeria&description=Cool&contact_email=test@test.com&city=Lagos&address=123+Street&type=Business"
	c, rec := setupTestContext(http.MethodPut, "/listings/listing-123", strings.NewReader(body))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("listing-123")
	c.Set("User", domain.User{ID: "user1", Email: "owner@example.com"})

	// Execute
	err := h.HandleUpdate(c)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify DB state
	updated, _ := app.DB.FindByID(context.Background(), "listing-123")
	assert.Empty(t, updated.ImageURL)
	assert.Equal(t, "New Title", updated.Title)

	mockImageService.AssertExpectations(t)
}
