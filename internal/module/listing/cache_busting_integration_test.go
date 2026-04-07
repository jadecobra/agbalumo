package listing_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestListingHandler_HandleImageUpload_CacheBusting(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockImageService := &testutil.MockImageService{}
	mockGeocodingService := &testutil.MockGeocodingService{}
	app.ImageSvc = mockImageService
	app.GeocodingSvc = mockGeocodingService

	h := listmod.NewListingHandler(app)

	// Mock successful upload returning a clean URL
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).
		Return("/static/uploads/test.webp", nil)

	// Multipart form request with image
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Cache Busting Test")
	_ = writer.WriteField("type", "Business")
	_ = writer.WriteField("owner_origin", "Nigeria")
	_ = writer.WriteField("description", "This is a long enough description for validation purposes.")
	_ = writer.WriteField("contact_email", "test@test.com")
	_ = writer.WriteField("address", "123 Test St")
	_ = writer.WriteField("city", "Lagos")
	part, _ := writer.CreateFormFile("image", "test.jpg")
	_, _ = part.Write([]byte("fake image content"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "u1"})

	// Execute
	err := h.HandleCreate(c)
	assert.NoError(t, err)

	// Fetch from DB to check ImageURL
	all, _, _ := app.DB.FindAll(c.Request().Context(), "", "", "", "", false, 10, 0)
	assert.Equal(t, 1, len(all))
	assert.Contains(t, all[0].ImageURL, "/static/uploads/test.webp?t=")

	mockImageService.AssertExpectations(t)
}
