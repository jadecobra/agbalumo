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
	_ = os.Setenv("BASE_URL", "http://192.168.1.5:8080")
	defer func() { _ = os.Unsetenv("BASE_URL") }()
	_ = os.Unsetenv("GOOGLE_REDIRECT_URL")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("http", "localhost:8080")

	assert.Equal(t, "http://192.168.1.5:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_GoogleRedirectURL(t *testing.T) {
	_ = os.Unsetenv("BASE_URL")
	_ = os.Setenv("GOOGLE_REDIRECT_URL", "https://custom.example.com/callback")
	defer func() { _ = os.Unsetenv("GOOGLE_REDIRECT_URL") }()

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "localhost:8080")

	assert.Equal(t, "https://custom.example.com/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTPS(t *testing.T) {
	_ = os.Unsetenv("BASE_URL")
	_ = os.Unsetenv("GOOGLE_REDIRECT_URL")
	_ = os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "localhost:8443")

	assert.Equal(t, "https://localhost:8443/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTP(t *testing.T) {
	_ = os.Unsetenv("BASE_URL")
	_ = os.Unsetenv("GOOGLE_REDIRECT_URL")
	_ = os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("http", "localhost:8080")

	assert.Equal(t, "http://localhost:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_Production(t *testing.T) {
	_ = os.Unsetenv("BASE_URL")
	_ = os.Unsetenv("GOOGLE_REDIRECT_URL")
	_ = os.Setenv("AGBALUMO_ENV", "production")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("https", "agbalumo.fly.dev")

	assert.Equal(t, "https://agbalumo.fly.dev/auth/google/callback", got)
}

// --- RealGoogleProvider.GetAuthCodeURL Test ---

func TestRealGoogleProvider_GetAuthCodeURL(t *testing.T) {
	_ = os.Unsetenv("BASE_URL")
	_ = os.Unsetenv("GOOGLE_REDIRECT_URL")
	_ = os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	url := p.GetAuthCodeURL("test-state", "http", "localhost:8080")

	assert.NotEmpty(t, url)
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "redirect_uri=")
}

// --- RealGoogleProvider.GetUserInfo Tests ---

func TestRealGoogleProvider_GetUserInfo_Success(t *testing.T) {
	expected := GoogleUser{
		ID:      "123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "http://pic.com/avatar.jpg",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), "access_token=test-token")
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
