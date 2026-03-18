package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractRoutes(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile1 := filepath.Join(tempDir, "server.go")
	sourceFile2 := filepath.Join(tempDir, "handler.go")

	code1 := `package cmd
import "github.com/labstack/echo/v4"
func setupRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error { return nil })
	e.POST("/users", func(c echo.Context) error { return nil })
	
	e.GET("/items/", func(c echo.Context) error { return nil })
	e.DELETE("/items/:id", func(c echo.Context) error { return nil })

	adminGroup := e.Group("/admin")
	adminGroup.GET("", func(c echo.Context) error { return nil })
	adminGroup.GET("/dashboard", func(c echo.Context) error { return nil })

	settingsGroup := adminGroup.Group("/settings")
	settingsGroup.POST("/update", func(c echo.Context) error { return nil })
}
`
	code2 := `package handler
import "github.com/labstack/echo/v4"
func RegisterUserRoutes(g *echo.Group) {
	g.GET("/profile", func(c echo.Context) error { return nil })
	g.PUT("/profile/update", func(c echo.Context) error { return nil })
}
`
	if err := os.WriteFile(sourceFile1, []byte(code1), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(sourceFile2, []byte(code2), 0644); err != nil {
		t.Fatalf("failed to write test file 2: %v", err)
	}

	routes, err := ExtractRoutes(tempDir)
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
		"GET /profile":          true,
		"PUT /profile/update":   true,
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
