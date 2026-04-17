package auth

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides routing middleware for authentication
type AuthMiddleware struct {
	Repo domain.UserStore
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(repo domain.UserStore) *AuthMiddleware {
	return &AuthMiddleware{Repo: repo}
}

// OptionalAuth injects user into context if session exists
func (m *AuthMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := middleware.GetSession(c)
		if sess != nil {
			if userID, ok := sess.Values[domain.SessionKeyUserID].(string); ok {

				user, err := m.Repo.FindUserByID(c.Request().Context(), userID)
				if err == nil {
					c.Set(domain.CtxKeyUser, &user)

				}
			}
		}
		return next(c)
	}
}

// RequireAuth redirects to login if no active session
func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := middleware.GetSession(c)
		authSuccess := false
		if sess != nil {
			if _, ok := sess.Values["user_id"].(string); ok {
				authSuccess = true
			}
		}

		if !authSuccess {
			// Redirect to Google Login
			return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
		}
		return next(c)
	}
}
