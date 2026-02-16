package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func TestSecureHeaders(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := middleware.SecureHeaders(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	if err := h(c); err != nil {
		t.Fatalf("Middleware failed: %v", err)
	}

	headers := rec.Header()
	
	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for k, v := range expectedHeaders {
		if got := headers.Get(k); got != v {
			t.Errorf("Expected header %s to be %q, got %q", k, v, got)
		}
	}

	if csp := headers.Get("Content-Security-Policy"); csp == "" {
		t.Error("Expected Content-Security-Policy header to be set")
	}

	if perm := headers.Get("Permissions-Policy"); perm == "" {
		t.Error("Expected Permissions-Policy header to be set")
	}
}
