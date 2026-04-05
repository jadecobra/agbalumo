package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleCreate_WithImage(t *testing.T) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Image Listing")
	_ = writer.WriteField("type", "Business")
	_ = writer.WriteField("owner_origin", "Ghana")
	_ = writer.WriteField("description", "Desc")
	_ = writer.WriteField("contact_email", "img@example.com")
	_ = writer.WriteField("address", "123 Image St")
	_ = writer.WriteField("city", "Accra")

	part, _ := writer.CreateFormFile("image", "test.png")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	_ = png.Encode(part, img)
	_ = writer.Close()

	c, rec := setupTestContext(http.MethodPost, "/listings", body)
	c.Request().Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := listmod.NewListingHandler(app)
	c.Set("User", domain.User{ID: "u1"})

	if err := h.HandleCreate(c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify DB state
	listings, _ := app.DB.FindByTitle(context.Background(), "Image Listing")
	assert.Len(t, listings, 1)
	assert.NotEmpty(t, listings[0].ImageURL)
}

func TestHandleCreate_InvalidDates(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid Deadline",
			body:           "title=T&type=Request&city=Lagos&deadline_date=invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, cleanup := testutil.SetupTestAppEnv(t)
			defer cleanup()
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			h := listmod.NewListingHandler(app)
			c.Set("User", domain.User{ID: "u1"})
			_ = h.HandleCreate(c)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleCreate_ImageUploadError(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockImageService := &MockImageService{}
	app.ImageSvc = mockImageService
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return("", errors.New("upload fail"))

	c, rec := setupTestContext(http.MethodPost, "/listings", nil)

	h := listmod.NewListingHandler(app)
	c.Set("User", domain.User{ID: "u1"})

	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandleProfile_NoUser(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	c, rec := setupTestContext(http.MethodGet, "/profile", nil)
	h := listmod.NewListingHandler(app)
	_ = h.HandleProfile(c)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}
