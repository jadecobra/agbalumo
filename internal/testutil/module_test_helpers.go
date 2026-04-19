package testutil

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
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

// SeedStandardData populates the environment's database with a representative set of listings.
func (e ModuleTestEnv) SeedStandardData(t *testing.T) {
	SeedStandardData(t, e.App.DB)
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

// SetupCSVUploadBody creates a multipart body containing a CSV file.
func SetupCSVUploadBody(t *testing.T, fieldName, fileName, content string) (*bytes.Buffer, string) {
	t.Helper()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, _ = part.Write([]byte(content))
	_ = writer.Close()
	return body, writer.FormDataContentType()
}

// AssertListingExists verifies that a listing with the given title exists in the database.
func AssertListingExists(t *testing.T, db domain.ListingRepository, title string) {
	t.Helper()
	listings, err := db.FindByTitle(context.Background(), title)
	if err != nil {
		t.Fatalf("failed to query database for title %s: %v", title, err)
	}
	if len(listings) == 0 {
		t.Errorf("expected listing with title %s to exist, but none were found", title)
	}
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
	return (len(body) > 0) && ((strings.Contains(body, `hx-swap-oob="true"`)) ||
		(strings.Contains(body, `hx-target`)) ||
		(strings.Contains(body, `hx-trigger`)))
}

// Re-implementing contains to avoid depending on it or using standard library
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
