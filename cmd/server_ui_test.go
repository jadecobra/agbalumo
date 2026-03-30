package cmd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/cmd"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	_ = os.Setenv("AGBALUMO_ENV", "development")
	// Keep ENV=test for test compatibility but set high rate limits to avoid 429 in tests
	_ = os.Setenv("RATE_LIMIT_RATE", "10000")
	_ = os.Setenv("RATE_LIMIT_BURST", "20000")
	_ = os.Setenv("DATABASE_URL", "file:test_ui.db?mode=memory&cache=shared")
	_ = os.Setenv("GOOGLE_CLIENT_ID", "dummy_client_id")
	_ = os.Setenv("GOOGLE_CLIENT_SECRET", "dummy_client_secret")
	// SetupServer handles seeding as long as ENV != "production"
	var err error

	// We need to change to the project root directory so template paths work
	_ = os.Chdir("..")

	e, err = cmd.SetupServer()
	if err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

func getSessionCookie(rec *httptest.ResponseRecorder) string {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "auth_session" {
			return cookie.String()
		}
	}
	return ""
}

func TestMobileFilterBottomSheet(t *testing.T) {
	// RED TEST: This test currently fails because the panel is "floating" (inset-x-4)
	// rather than being a bottom sheet (bottom-0 left-0 right-0).
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()

	// Expected Bottom Sheet behavior for mobile
	assert.Contains(t, body, "fixed md:absolute")
	assert.Contains(t, body, "bottom-0 md:bottom-auto")
	assert.Contains(t, body, "left-0 right-0")
	assert.NotContains(t, body, "inset-x-4")
	assert.Contains(t, body, "style=\"z-index: 120;\"")     // Higher z-index as inline style
	assert.Contains(t, body, "style=\"max-height: 55vh;\"") // Reduced height as inline style
	assert.Contains(t, body, "rounded-t-3xl")
	assert.Contains(t, body, "bg-earth-dark/20 rounded-full")
}
