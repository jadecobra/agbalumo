package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleUpdate_ImageRemoval(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// Existing listing with an image
	existingListing := domain.Listing{
		ID:       "listing-123",
		OwnerID:  "user1",
		Title:    "Old Title",
		ImageURL: "/static/uploads/listing-123.webp",
	}

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "listing-123").Return(existingListing, nil)
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()

	// Expect Save with ImageURL = ""
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.ID == "listing-123" && l.ImageURL == ""
	})).Return(nil)

	mockImageService := &mock.MockImageService{}
	// Expect DeleteImage to be called
	mockImageService.On("DeleteImage", testifyMock.Anything, existingListing.ImageURL).Return(nil)
	// No upload expected
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return("", nil).Maybe()

	h := handler.NewListingHandler(mockRepo, mockImageService)

	// Body with remove_image=true and required fields
	body := "title=New+Title&remove_image=true&owner_origin=Nigeria&description=Cool&contact_email=test@test.com&address=123+Street&type=Business"
	req := httptest.NewRequest(http.MethodPut, "/listings/listing-123", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("listing-123")
	c.Set("User", domain.User{ID: "user1", Email: "owner@example.com"})

	// Execute
	err := h.HandleUpdate(c)
	if err != nil {
		t.Fatalf("HandleUpdate failed: %v", err)
	}

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	mockRepo.AssertExpectations(t)
	mockImageService.AssertExpectations(t)
}
