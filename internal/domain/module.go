package domain

import "github.com/labstack/echo/v4"

// AuthMiddleware defines the required methods for route authentication middleware
type AuthMiddleware interface {
	OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc
	RequireAuth(next echo.HandlerFunc) echo.HandlerFunc
}

// Registrar defines the interface that all vertical slices (modules) must implement
// to register their HTTP routes with the main echo instance.
type Registrar interface {
	RegisterRoutes(e *echo.Echo, authMw AuthMiddleware)
}
