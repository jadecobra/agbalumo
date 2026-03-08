package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStaticCacheHeaders(t *testing.T) {
	e := echo.New()
	mw := staticCacheHeaders()

	tests := []struct {
		name           string
		path           string
		wantCache      bool
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
