package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const sessionContextKey = "session"

// SessionMiddleware returns a middleware that attaches a session to the context.
func SessionMiddleware(store sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session
			session, _ := store.Get(c.Request(), "auth_session")
			// Make session available in context
			c.Set(sessionContextKey, session)

			// Proceed
			err := next(c)

			// Save session before response (if not already committed)
			// Note: Gorilla sessions need Save called before writing response headers.
			// However, typical pattern is to call Save manually when modifying.
			// Or we can defer? But err handling...
			// For simple auth, we save in the Handler (Login/Logout).
			return err
		}
	}
}

// GetSession retrieves the session from the context.
func GetSession(c echo.Context) *sessions.Session {
	if sess, ok := c.Get(sessionContextKey).(*sessions.Session); ok {
		return sess
	}
	// Fallback if not in context (should catch in middleware but safe guard)
	return nil
}
