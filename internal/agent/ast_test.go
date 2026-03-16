package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractRoutes(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "server.go")

	code := `package cmd

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo) {
	// Root routes
	e.GET("/", func(c echo.Context) error { return nil })
	e.POST("/users", func(c echo.Context) error { return nil })
	
	// Item route with trailing slash normalization needed
	e.GET("/items/", func(c echo.Context) error { return nil })
	e.DELETE("/items/:id", func(c echo.Context) error { return nil })

	// Admin group
	adminGroup := e.Group("/admin")
	adminGroup.GET("", func(c echo.Context) error { return nil })
	adminGroup.GET("/dashboard", func(c echo.Context) error { return nil })

	// Nested group
	settingsGroup := adminGroup.Group("/settings")
	settingsGroup.POST("/update", func(c echo.Context) error { return nil })
}
`
	if err := os.WriteFile(sourceFile, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	routes, err := ExtractRoutes(sourceFile)
	if err != nil {
		t.Fatalf("ExtractRoutes failed: %v", err)
	}

	expectedRoutes := map[string]bool{
		"GET /":                 true,
		"POST /users":           true,
		"GET /items":            true,
		"DELETE /items/{id}":    true,
		"GET /admin":            true,
		"GET /admin/dashboard":  true,
		"POST /admin/settings/update": true,
	}

	if len(routes) != len(expectedRoutes) {
		t.Errorf("expected %d routes, got %d", len(expectedRoutes), len(routes))
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if !expectedRoutes[key] {
			t.Errorf("unexpected route found: %s", key)
		} else {
			delete(expectedRoutes, key)
		}
	}

	for k := range expectedRoutes {
		t.Errorf("missing expected route: %s", k)
	}
}
