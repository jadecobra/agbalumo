package handler_test

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleCreate_WithImage(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Image Listing")
	_ = writer.WriteField("type", "Business")
	_ = writer.WriteField("owner_origin", "Ghana")
	_ = writer.WriteField("description", "Desc")
	_ = writer.WriteField("contact_email", "img@example.com")
	_ = writer.WriteField("address", "123 Image St")

	part, _ := writer.CreateFormFile("image", "test.png")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	_ = png.Encode(part, img)
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.ImageURL != ""
	})).Return(nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	c.Set("User", domain.User{ID: "u1"})

	if err := h.HandleCreate(c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleCreate_InvalidDates(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid Deadline",
			body:           "title=T&type=Request&deadline_date=invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			h := handler.NewListingHandler(nil, nil, "")
			c.Set("User", domain.User{ID: "u1"})
			_ = h.HandleCreate(c)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleCreate_ImageUploadError(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	_, _ = part.Write([]byte("fake"))
	_ = writer.Close()

	c, rec := setupTestContext(http.MethodPost, "/listings", body)
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()

	mockImageService := &MockImageService{}
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return("", errors.New("upload fail"))

	h := handler.NewListingHandler(mockRepo, mockImageService, "")
	c.Set("User", domain.User{ID: "u1"})

	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandleProfile_NoUser(t *testing.T) {
	c, rec := setupTestContext(http.MethodGet, "/profile", nil)
	h := handler.NewListingHandler(nil, nil, "")
	_ = h.HandleProfile(c)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}
