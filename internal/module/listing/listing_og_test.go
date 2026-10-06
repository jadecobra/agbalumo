package listing_test

import (
	"bytes"
	"image"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleOGImage_KnownListing(t *testing.T) {
	t.Parallel()

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	testutil.SaveTestListing(t, env.App.DB, "aria-suya-og", "Aria Suya Kitchen", func(l *domain.Listing) {
		l.City = "Arlington"
		l.Type = domain.Food
		l.Rating = 4.8
		l.ReviewCount = 120
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/listings/aria-suya-og/og.png", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("aria-suya-og")

	h := listing.NewListingHandler(env.App)
	err := h.HandleOGImage(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=86400", rec.Header().Get("Cache-Control"))

	// Check PNG magic bytes
	body := rec.Body.Bytes()
	require.True(t, len(body) >= 8)
	expectedMagic := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	assert.Equal(t, expectedMagic, body[:8])

	// Check dimensions
	cfg, format, err := image.DecodeConfig(bytes.NewReader(body))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 1200, cfg.Width)
	assert.Equal(t, 630, cfg.Height)
}

func TestHandleOGImage_UnknownListing(t *testing.T) {
	t.Parallel()

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/listings/unknown-id/og.png", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("unknown-id")

	h := listing.NewListingHandler(env.App)
	err := h.HandleOGImage(c)
	require.NoError(t, err)

	// Missing listing still returns 200 with fallback card (never 404 or 5xx for crawlers)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=86400", rec.Header().Get("Cache-Control"))

	body := rec.Body.Bytes()
	require.True(t, len(body) >= 8)
	expectedMagic := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	assert.Equal(t, expectedMagic, body[:8])

	cfg, format, err := image.DecodeConfig(bytes.NewReader(body))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 1200, cfg.Width)
	assert.Equal(t, 630, cfg.Height)
}
