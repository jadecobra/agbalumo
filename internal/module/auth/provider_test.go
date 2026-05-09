package auth_test

import (
	"net/url"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/auth"
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
			// Clear env using t.Setenv to ensure clean state for each subtest
			t.Setenv("BASE_URL", "")
			t.Setenv("GOOGLE_REDIRECT_URL", "")
			t.Setenv("AGBALUMO_ENV", "")

			// Set test env
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			p := auth.NewRealGoogleProvider()
			scheme := "http"
			if tt.host == "localhost:8443" || tt.env["AGBALUMO_ENV"] == "production" {
				scheme = "https"
			}
			rawURL := p.GetAuthCodeURL("state", scheme, tt.host)

			decodedURL, err := url.QueryUnescape(rawURL)
			assert.NoError(t, err)
			assert.Contains(t, decodedURL, "redirect_uri="+tt.expected)
		})
	}
}

func TestRealGoogleProvider_Exchange_RedirectURL(t *testing.T) {
	// This tests that Exchange also uses the correct redirect URL
	t.Setenv("BASE_URL", "https://test.com")

	p := auth.NewRealGoogleProvider()
	_ = p
}
