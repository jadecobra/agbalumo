package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func okHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func TestSecureHeaders(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := middleware.SecureHeaders(okHandler)

	if err := h(c); err != nil {
		t.Fatalf("Middleware failed: %v", err)
	}

	headers := rec.Header()

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for k, v := range expectedHeaders {
		if got := headers.Get(k); got != v {
			t.Errorf("Expected header %s to be %q, got %q", k, v, got)
		}
	}

	if got := headers.Get("Content-Security-Policy"); got == "" {
		t.Error("Expected Content-Security-Policy header")
	} else {
		checkScriptSrc(t, got)
	}

	if perm := headers.Get("Permissions-Policy"); perm == "" {
		t.Error("Expected Permissions-Policy header")
	}
}

func checkScriptSrc(t *testing.T, csp string) {
	for _, part := range strings.Split(csp, ";") {
		trimmed := strings.TrimSpace(part)
		if strings.HasPrefix(trimmed, "script-src") && strings.Contains(trimmed, "'unsafe-inline'") {
			t.Error("script-src should NOT contain 'unsafe-inline'")
		}
	}
}

func TestCanonicalPath(t *testing.T) {
	t.Parallel()
	e := echo.New()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "valid path",
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid root path",
			path:       "/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "malformed path - missing leading slash",
			path:       "healthz",
			wantStatus: http.StatusNotImplemented,
		},
		{
			name:       "empty path",
			path:       "",
			wantStatus: http.StatusNotImplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "https://localhost:8443", nil)
			req.URL.Path = tt.path

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := middleware.CanonicalPath(okHandler)

			err := h(c)
			assertCanonicalPathResult(t, tt.path, tt.wantStatus, err, rec)
		})
	}
}

func assertCanonicalPathResult(t *testing.T, path string, wantStatus int, err error, rec *httptest.ResponseRecorder) {
	t.Helper()
	if wantStatus == http.StatusOK {
		if err != nil {
			t.Errorf("expected no error for path %q, got %v", path, err)
		}
		if rec.Code != wantStatus {
			t.Errorf("expected status %d for path %q, got %d", wantStatus, path, rec.Code)
		}
	} else {
		if err == nil {
			t.Errorf("expected error for path %q", path)
		} else {
			he, ok := err.(*echo.HTTPError)
			if !ok {
				t.Errorf("expected echo.HTTPError for path %q, got %T", path, err)
			} else if he.Code != wantStatus {
				t.Errorf("expected status %d for path %q, got %d", wantStatus, path, he.Code)
			}
		}
	}
}
