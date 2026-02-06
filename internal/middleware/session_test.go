package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func TestSessionMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))

	// Middleware under test
	mw := middleware.SessionMiddleware(store)

	// Next handler verifies session is present
	handler := func(c echo.Context) error {
		sess := middleware.GetSession(c)
		if sess == nil {
			return c.String(http.StatusInternalServerError, "Session is nil")
		}

		// Set a value
		sess.Values["foo"] = "bar"

		// Directly saving here since typical usage requires handlers to save if modifying
		// But middleware might not save automatically unless configured (the current impl doesn't auto-save on return)
		// We'll verify retrieval is enough.
		return c.String(http.StatusOK, "OK")
	}

	if err := mw(handler)(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}
}

func TestGetSession_Fallback(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// No middleware run, so session not in context
	sess := middleware.GetSession(c)
	if sess != nil {
		t.Error("Expected nil session when not in context")
	}
}
