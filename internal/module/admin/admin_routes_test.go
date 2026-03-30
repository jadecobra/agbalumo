package admin_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/labstack/echo/v4"
)

// mockAuthMiddleware implements domain.AuthMiddleware for testing.
type mockAuthMiddleware struct{}

func (m *mockAuthMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (m *mockAuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func TestAdminHandler_RegisterRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	handler := &admin.AdminHandler{}
	authMw := &mockAuthMiddleware{}

	// Execute
	handler.RegisterRoutes(e, authMw)

	// Verify
	routes := e.Routes()
	expectedRoutes := map[string]string{
		"/admin/login":                   http.MethodGet,
		"/admin/login*":                  http.MethodPost, // Note: the POST route is on the subgroup without trailing slash, Echo might represent it differently or as /admin/login
		"/admin":                         http.MethodGet,
		"/admin/users":                   http.MethodGet,
		"/admin/listings":                http.MethodGet,
		"/admin/claims/:id/approve":      http.MethodPost,
		"/admin/claims/:id/reject":       http.MethodPost,
		"/admin/listings/bulk":           http.MethodPost,
		"/admin/listings/:id/row":        http.MethodGet,
		"/admin/listings/delete-confirm": http.MethodGet,
		"/admin/listings/delete":         http.MethodPost,
		"/admin/listings/:id/featured":   http.MethodPost,
		"/admin/upload":                  http.MethodPost,
		"/admin/listings/export":         http.MethodGet,
		"/admin/categories":              http.MethodPost,
	}

	// Build a map of registered routes for easy lookup
	registered := make(map[string]map[string]bool)
	for _, r := range routes {
		if registered[r.Path] == nil {
			registered[r.Path] = make(map[string]bool)
		}
		registered[r.Path][r.Method] = true
	}

	for path, method := range expectedRoutes {
		// Echo might register paths with or without trailing slash depending on group definition.
		// For the "" POST route on the group "/admin/login", it's usually registered as "/admin/login".
		// We'll normalize the lookup.
		checkPath := path
		if path == "/admin/login*" {
			checkPath = "/admin/login"
		}

		methods, exists := registered[checkPath]
		if !exists {
			t.Errorf("Expected route %s not registered", checkPath)
			continue
		}
		if !methods[method] {
			t.Errorf("Expected route %s to handle method %s, but it doesn't", checkPath, method)
		}
	}
}
