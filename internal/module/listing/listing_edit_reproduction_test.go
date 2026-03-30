package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListingHandler_HandleUpdate_Reproduction(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	ctx := context.Background()

	// 1. Create a listing to edit
	listing := domain.Listing{
		ID:          "test-id-123",
		OwnerID:     "user-1",
		Title:       "Original Title",
		Description: "Original Description",
		Type:        domain.Business,
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	}
	err := repo.Save(ctx, listing)
	assert.NoError(t, err)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:  repo,
		CategoryStore: repo,
		ListingSvc:    listingSvc,
		ImageService:  nil,
		GeocodingSvc:  &MockGeocodingService{},
		Config:        &config.Config{},
	})

	// 2. Prepare update data
	updatedTitle := "Updated Title"

	// 3. Simulate PUT request with multipart/form-data (as HTMX does when enctype is set)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("title", updatedTitle)
	_ = w.WriteField("type", "Business")
	_ = w.WriteField("owner_origin", "Nigeria")
	_ = w.WriteField("description", "Updated Description")
	_ = w.WriteField("city", "Lagos")
	_ = w.WriteField("address", "123 Street")
	_ = w.WriteField("contact_email", "test@example.com")
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPut, "/listings/test-id-123", &b)
	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-id-123")
	c.Set("User", domain.User{ID: "user-1"})

	// 4. Call handler
	// Simulate the redundant bind in HandleUpdate to see if it breaks anything
	err = h.HandleUpdate(c)
	assert.NoError(t, err)

	// 5. Verify response and DB
	if rec.Code != http.StatusOK {
		t.Logf("Response Body: %s", rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)

	updatedListing, err := repo.FindByID(ctx, "test-id-123")
	assert.NoError(t, err)
	assert.Equal(t, updatedTitle, updatedListing.Title, "Title should be updated in database")
}

func TestListingHandler_HandleUpdate_AdminSource(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	ctx := context.Background()

	// 1. Create a listing to edit
	listing := domain.Listing{
		ID:          "test-id-admin",
		OwnerID:     "user-1",
		Title:       "Original Title",
		Description: "Original Description",
		Type:        domain.Business,
		IsActive:    true,
		Status:      domain.ListingStatusApproved,
	}
	err := repo.Save(ctx, listing)
	assert.NoError(t, err)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:  repo,
		CategoryStore: repo,
		ListingSvc:    listingSvc,
		ImageService:  nil,
		GeocodingSvc:  &MockGeocodingService{},
		Config:        &config.Config{},
	})

	// 2. Prepare update data
	updatedTitle := "Updated Title Admin"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("title", updatedTitle)
	_ = w.WriteField("type", "Business")
	_ = w.WriteField("owner_origin", "Nigeria")
	_ = w.WriteField("description", "Updated Description")
	_ = w.WriteField("city", "Lagos")
	_ = w.WriteField("address", "123 Street")
	_ = w.WriteField("contact_email", "test@example.com")
	_ = w.Close()

	// 3. Request with source=admin
	req := httptest.NewRequest(http.MethodPut, "/listings/test-id-admin?source=admin", &b)
	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-id-admin")
	c.Set("User", domain.User{ID: "admin-1", Role: domain.UserRoleAdmin})

	// 4. Call handler
	err = h.HandleUpdate(c)
	assert.NoError(t, err)

	// 5. Verify response and DB
	if rec.Code != http.StatusOK {
		t.Logf("Response Body: %s", rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)

	// Ensure the response contains the HX-Trigger header and no content
	assert.Equal(t, "listing-updated-test-id-admin", rec.Header().Get("HX-Trigger"), "Response should trigger listing-updated event")
	assert.Empty(t, rec.Body.String(), "Response should be empty for admin source update")

	updatedListing, err := repo.FindByID(ctx, "test-id-admin")
	assert.NoError(t, err)
	assert.Equal(t, updatedTitle, updatedListing.Title, "Title should be updated in database")
}
