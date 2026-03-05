package handler_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestListingHandler_HandleImageUpload_CacheBusting(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	mockRepo := &mock.MockListingRepository{}
	mockImageService := &mock.MockImageService{}

	h := handler.NewListingHandler(mockRepo, mockImageService)

	// Mock successful upload returning a clean URL
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).
		Return("/static/uploads/test.webp", nil)

	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()

	var savedListing domain.Listing
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		savedListing = l
		return true
	})).Return(nil)

	// Multipart form request with image
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Cache Busting Test")
	writer.WriteField("type", "Business")
	writer.WriteField("owner_origin", "Nigeria")
	writer.WriteField("description", "Desc")
	writer.WriteField("contact_email", "test@test.com")
	writer.WriteField("address", "123 Test St")
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "u1"})

	// Execute
	err := h.HandleCreate(c)
	assert.NoError(t, err)

	// Assert cache-busting parameter exists
	assert.Contains(t, savedListing.ImageURL, "/static/uploads/test.webp?t=")

	mockRepo.AssertExpectations(t)
	mockImageService.AssertExpectations(t)
}
