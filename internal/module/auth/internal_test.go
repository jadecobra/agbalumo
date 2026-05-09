package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// --- RealGoogleProvider.getRedirectURL Tests ---

func TestRealGoogleProvider_getRedirectURL_BaseURL(t *testing.T) {
	t.Setenv("BASE_URL", "http://192.168.1.5:8080")

	// Ensure clean env for other vars
	for _, key := range []string{"GOOGLE_REDIRECT_URL"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("http", "localhost:8080")

	assert.Equal(t, "http://192.168.1.5:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_GoogleRedirectURL(t *testing.T) {
	t.Setenv("GOOGLE_REDIRECT_URL", "https://custom.example.com/callback")

	// Ensure clean env for other vars
	for _, key := range []string{"BASE_URL"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "localhost:8080")

	assert.Equal(t, "https://custom.example.com/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTPS(t *testing.T) {
	for _, key := range []string{"BASE_URL", "GOOGLE_REDIRECT_URL", "AGBALUMO_ENV"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "localhost:8443")

	assert.Equal(t, "https://localhost:8443/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTP(t *testing.T) {
	for _, key := range []string{"BASE_URL", "GOOGLE_REDIRECT_URL", "AGBALUMO_ENV"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("http", "localhost:8080")

	assert.Equal(t, "http://localhost:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_Production(t *testing.T) {
	t.Setenv("AGBALUMO_ENV", "production")

	// Ensure clean env for other vars
	for _, key := range []string{"BASE_URL", "GOOGLE_REDIRECT_URL"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "agbalumo.fly.dev")

	assert.Equal(t, "https://agbalumo.fly.dev/auth/google/callback", got)
}

// --- RealGoogleProvider.GetAuthCodeURL Test ---

func TestRealGoogleProvider_GetAuthCodeURL(t *testing.T) {
	for _, key := range []string{"BASE_URL", "GOOGLE_REDIRECT_URL", "AGBALUMO_ENV"} {
		old, existed := os.LookupEnv(key)
		_ = os.Unsetenv(key)
		t.Cleanup(func() {
			if existed {
				_ = os.Setenv(key, old)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}

	p := NewRealGoogleProvider()
	url := p.GetAuthCodeURL("test-state", "http", "localhost:8080")

	assert.NotEmpty(t, url)
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "redirect_uri=")
}

// --- RealGoogleProvider.GetUserInfo Tests ---

func TestRealGoogleProvider_GetUserInfo_Success(t *testing.T) {
	t.Parallel()
	expected := GoogleUser{
		ID:      "123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "http://pic.com/avatar.jpg",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	p := NewRealGoogleProvider()
	p.UserInfoURL = server.URL

	token := &oauth2.Token{AccessToken: "test-token"}
	user, err := p.GetUserInfo(t.Context(), token)

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Picture, user.Picture)
}

func TestRealGoogleProvider_GetUserInfo_BadStatus(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "server error")
	}))
	defer server.Close()

	p := NewRealGoogleProvider()
	p.UserInfoURL = server.URL
	token := &oauth2.Token{AccessToken: "test-token"}
	user, err := p.GetUserInfo(t.Context(), token)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to fetch user info")
}

func TestRealGoogleProvider_GetUserInfo_BadJSON(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "not valid json{{{")
	}))
	defer server.Close()

	p := NewRealGoogleProvider()
	p.UserInfoURL = server.URL
	token := &oauth2.Token{AccessToken: "test-token"}
	user, err := p.GetUserInfo(t.Context(), token)

	assert.Error(t, err)
	assert.Nil(t, user)
}
