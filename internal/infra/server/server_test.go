package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStaticCacheHeaders(t *testing.T) {
	t.Parallel()
	e := echo.New()
	mw := StaticCacheHeaders()

	tests := []struct {
		name      string
		path      string
		wantCache bool
	}{
		{"CSS file", "/static/css/style.css", true},
		{"JS file", "/static/js/app.js", true},
		{"PNG file", "/static/img/logo.png", true},
		{"JPG file", "/static/img/photo.jpg", true},
		{"JPEG file", "/static/img/photo.jpeg", true},
		{"SVG file", "/static/img/icon.svg", true},
		{"WOFF2 file", "/static/fonts/font.woff2", true},
		{"WOFF file", "/static/fonts/font.woff", true},
		{"HTML file", "/home", false},
		{"No extension", "/static/data", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := mw(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			err := handler(c)
			assert.NoError(t, err)

			cacheControl := rec.Header().Get("Cache-Control")
			if tt.wantCache {
				assert.Equal(t, "public, max-age=31536000, immutable", cacheControl)
			} else {
				assert.Empty(t, cacheControl)
			}
		})
	}
}

func TestPublicRoutes_HeadRequest(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	setupMiddleware(e, app.Cfg)
	setupRoutes(e, app)

	// Seed a test listing
	listingID := "5e195713-5a55-4f0b-b145-9741d3583660"
	testutil.SaveTestListing(t, app.DB, listingID, "Mama Put Dallas")

	// 1. HEAD /listings/{id} should return 200 OK with empty body
	req := httptest.NewRequest(http.MethodHead, "/listings/"+listingID, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.String())

	// 2. HEAD / (home) should return 200 OK with empty body
	reqHome := httptest.NewRequest(http.MethodHead, "/", nil)
	recHome := httptest.NewRecorder()
	e.ServeHTTP(recHome, reqHome)

	assert.Equal(t, http.StatusOK, recHome.Code)
	assert.Empty(t, recHome.Body.String())

	// 3. HEAD to admin route should redirect to login (302/307)
	adminReq := httptest.NewRequest(http.MethodHead, "/admin", nil)
	adminRec := httptest.NewRecorder()
	e.ServeHTTP(adminRec, adminReq)

	assert.True(t, adminRec.Code == http.StatusFound || adminRec.Code == http.StatusTemporaryRedirect, "expected redirect status, got: %d", adminRec.Code)
	assert.Contains(t, adminRec.Header().Get("Location"), "/auth/google/login")
}

func TestHomepage_OpenGraphSharingMeta_WithUTM(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewRealTemplate(t)}
	setupMiddleware(e, app.Cfg)
	setupRoutes(e, app)

	// Homepage with ?action=post and UTM parameters
	req := httptest.NewRequest(http.MethodGet, "/?action=post&utm_source=facebook&utm_medium=social&utm_campaign=quality_index", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()

	assert.Contains(t, body, `<meta property="og:image" content="https://agbalumo.com/og.png" />`)
	assert.Contains(t, body, `<meta property="og:image:width" content="1200" />`)
	assert.Contains(t, body, `<meta property="og:image:height" content="630" />`)
	assert.Contains(t, body, `<meta property="og:image:alt" content="agbalumo, find African food in DFW in under 60 seconds" />`)
	assert.Contains(t, body, `<meta name="twitter:image" content="https://agbalumo.com/og.png" />`)

	// Also verify /og.png serves the brand card via server routing
	ogReq := httptest.NewRequest(http.MethodGet, "/og.png", nil)
	ogRec := httptest.NewRecorder()
	e.ServeHTTP(ogRec, ogReq)

	assert.Equal(t, http.StatusOK, ogRec.Code)
	assert.Equal(t, "image/png", ogRec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=86400", ogRec.Header().Get("Cache-Control"))
	assert.True(t, len(ogRec.Body.Bytes()) >= 8)

	// HEAD /og.png
	ogHeadReq := httptest.NewRequest(http.MethodHead, "/og.png", nil)
	ogHeadRec := httptest.NewRecorder()
	e.ServeHTTP(ogHeadRec, ogHeadReq)

	assert.Equal(t, http.StatusOK, ogHeadRec.Code)
	assert.Equal(t, "image/png", ogHeadRec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=86400", ogHeadRec.Header().Get("Cache-Control"))
	assert.Empty(t, ogHeadRec.Body.String())
}
