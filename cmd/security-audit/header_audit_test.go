package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHeaders(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		expectedPassed bool
	}{
		{
			name: "All Headers Present",
			headers: map[string]string{
				"Strict-Transport-Security": "max-age=63072000",
				"Content-Security-Policy":   "default-src 'self'",
				"X-Frame-Options":           "DENY",
			},
			expectedPassed: true,
		},
		{
			name: "Missing HSTS",
			headers: map[string]string{
				"Content-Security-Policy": "default-src 'self'",
				"X-Frame-Options":         "DENY",
			},
			expectedPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tt.headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			passed, skipped := checkHeaders(server.URL, server.Client())
			if passed != tt.expectedPassed {
				t.Errorf("expected passed=%v, got %v", tt.expectedPassed, passed)
			}
			if skipped {
				t.Errorf("expected skipped=false, got true")
			}
		})
	}
}

func TestCheckHeadersConnectionFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	passed, skipped := checkHeaders("https://invalid.localhost.nowhere:9999", server.Client())
	if passed != false {
		t.Errorf("Expected passed=false for connection failure")
	}
	if !skipped {
		t.Errorf("Expected skipped=true for connection failure")
	}
}
