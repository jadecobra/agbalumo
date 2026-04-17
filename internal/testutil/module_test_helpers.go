package testutil

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/labstack/echo/v4"
)

// ModuleTestEnv carries both the app environment and the cleanup function for convenience.
type ModuleTestEnv struct {
	App     *env.AppEnv
	Cleanup func()
}

// SetupTestModuleEnv provides a unified way to initialize the application environment for module tests.
func SetupTestModuleEnv(t *testing.T) ModuleTestEnv {
	t.Helper()
	app, cleanup := SetupTestAppEnv(t)
	return ModuleTestEnv{
		App:     app,
		Cleanup: cleanup,
	}
}

func SetupModuleContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	return SetupTestContext(method, target, body)
}
 
// SetupAdminContext prepares an Echo context with an admin user and session.
func SetupAdminContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	c, rec := SetupTestContextWithSession(method, target, body)
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	return c, rec
}

// AssertHTMXResponse checks for common HTMX response headers or attributes.
func AssertHTMXResponse(t testing.TB, body string) {
	t.Helper()
	// Common check for OOB swaps or HTMX-specific markers
	if !containsHTMXMarkers(body) {
		t.Error("response does not appear to be a valid HTMX fragment swap")
	}
}

func containsHTMXMarkers(body string) bool {
	return (len(body) > 0) && (
		(contains(body, `hx-swap-oob="true"`)) || 
		(contains(body, `hx-target`)) || 
		(contains(body, `hx-trigger`)))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (contains(s[1:], substr) || s[:len(substr)] == substr)))
}
