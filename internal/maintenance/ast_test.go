package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractRoutes(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "extract_routes_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a sample Go file with routes
	goCode := `
package test
import "github.com/labstack/echo/v4"
func Register(e *echo.Echo) {
	e.GET("/health", nil)
	v1 := e.Group("/api/v1")
	v1.POST("/login", nil)
	users := v1.Group("/users")
	users.GET("/:id", nil)
}
`
	err = os.WriteFile( /*nolint:gosec*/ filepath.Join(tmpDir, "routes.go"), []byte(goCode), 0600)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	routes, err := ExtractRoutes(tmpDir)
	if err != nil {
		t.Fatalf("ExtractRoutes failed: %v", err)
	}

	expected := []Route{
		{Method: "GET", Path: "/api/v1/users/:id"},
		{Method: "POST", Path: "/api/v1/login"},
		{Method: "GET", Path: "/health"},
	}

	verifyExtractedRoutes(t, expected, routes)
}

func verifyExtractedRoutes(t *testing.T, expected, routes []Route) {
	if len(routes) != len(expected) {
		t.Errorf("expected %d routes, got %d", len(expected), len(routes))
	}

	found := 0
	for _, exp := range expected {
		for _, got := range routes {
			if got.Method == exp.Method && got.Path == exp.Path {
				found++
				break
			}
		}
	}

	if found != len(expected) {
		t.Errorf("only found %d out of %d expected routes", found, len(expected))
		t.Logf("Found routes: %v", routes)
	}
}

func TestExtractRoutes_WithDomainConstants(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "extract_routes_domain_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	goCode := `
package test
import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)
func Register(e *echo.Echo) {
	admin := e.Group(domain.PathAdmin)
	admin.GET(domain.PathLogin, nil)
	admin.POST(domain.PathLogin, nil)
	
	e.GET(domain.PathListingID, nil)
	e.GET(domain.PathListingID+"/edit", nil)
	
	auth := e.Group("", nil)
	auth.DELETE(domain.PathListingID, nil)
	auth.POST(domain.PathListingID+"/claim", nil)
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "routes.go"), []byte(goCode), 0600)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	routes, err := ExtractRoutes(tmpDir)
	if err != nil {
		t.Fatalf("ExtractRoutes failed: %v", err)
	}

	expected := []Route{
		{Method: "GET", Path: "/admin/login"},
		{Method: "POST", Path: "/admin/login"},
		{Method: "GET", Path: "/listings/:id"},
		{Method: "GET", Path: "/listings/:id/edit"},
		{Method: "DELETE", Path: "/listings/:id"},
		{Method: "POST", Path: "/listings/:id/claim"},
	}

	verifyExtractedRoutes(t, expected, routes)
}
