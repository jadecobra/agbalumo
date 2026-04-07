package user

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetUser retrieves the authenticated user from the context.
// Returns (user, true) if user exists and is valid, (nil, false) otherwise.
// Handles both domain.User value and *domain.User pointer.
func GetUser(c echo.Context) (*domain.User, bool) {
	user := c.Get("User")
	if user == nil {
		return nil, false
	}
	// Handle both domain.User value and *domain.User pointer
	switch u := user.(type) {
	case *domain.User:
		return u, true
	case domain.User:
		return &u, true
	default:
		return nil, false
	}
}

// MustUser retrieves the authenticated user from the context.
// Panics if user is not present - use only when auth is guaranteed by middleware.
func MustUser(c echo.Context) *domain.User {
	user, ok := GetUser(c)
	if !ok {
		panic("user not authenticated")
	}
	return user
}

// RequireUser retrieves the authenticated user from the context.
// Redirects to login and returns the redirect error if user is missing.
func RequireUser(c echo.Context) (*domain.User, error) {
	u, ok := GetUser(c)
	if !ok || u == nil {
		return nil, c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}
	return u, nil
}
