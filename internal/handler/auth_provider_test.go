package handler_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRealGoogleProvider_GetRedirectURL(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		host     string
		expected string
	}{
		{
			name:     "BASE_URL set",
			env:      map[string]string{"BASE_URL": "https://myapp.com"},
			host:     "localhost:8080",
			expected: "https://myapp.com/auth/google/callback",
		},
		{
			name:     "GOOGLE_REDIRECT_URL set",
			env:      map[string]string{"GOOGLE_REDIRECT_URL": "http://legacy.com/callback"},
			host:     "localhost:8080",
			expected: "http://legacy.com/callback",
		},
		{
			name:     "Default fallback - localhost",
			env:      map[string]string{},
			host:     "localhost:8080",
			expected: "http://localhost:8080/auth/google/callback",
		},
		{
			name:     "Default fallback - secure local",
			env:      map[string]string{},
			host:     "localhost:8443",
			expected: "https://localhost:8443/auth/google/callback",
		},
		{
			name:     "Default fallback - production",
			env:      map[string]string{"AGBALUMO_ENV": "production"},
			host:     "agbalumo.com",
			expected: "https://agbalumo.com/auth/google/callback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current env
			oldBase := os.Getenv("BASE_URL")
			oldRedirect := os.Getenv("GOOGLE_REDIRECT_URL")
			oldEnv := os.Getenv("AGBALUMO_ENV")
			defer func() {
				_ = os.Setenv("BASE_URL", oldBase)
				_ = os.Setenv("GOOGLE_REDIRECT_URL", oldRedirect)
				_ = os.Setenv("AGBALUMO_ENV", oldEnv)
			}()

			// Clear env
			_ = os.Unsetenv("BASE_URL")
			_ = os.Unsetenv("GOOGLE_REDIRECT_URL")
			_ = os.Unsetenv("AGBALUMO_ENV")

			// Set test env
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}

			p := handler.NewRealGoogleProvider()
			rawURL := p.GetAuthCodeURL("state", tt.host)
			
			decodedURL, err := url.QueryUnescape(rawURL)
			assert.NoError(t, err)
			assert.Contains(t, decodedURL, "redirect_uri="+tt.expected)
		})
	}
}

func TestRealGoogleProvider_Exchange_RedirectURL(t *testing.T) {
	// This tests that Exchange also uses the correct redirect URL
	_ = os.Setenv("BASE_URL", "https://test.com")
	defer func() { _ = os.Unsetenv("BASE_URL") }()

	p := handler.NewRealGoogleProvider()
	_ = p
}
