package handler

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
	os.Setenv("BASE_URL", "http://192.168.1.5:8080")
	defer os.Unsetenv("BASE_URL")
	os.Unsetenv("GOOGLE_REDIRECT_URL")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("localhost:8080")

	assert.Equal(t, "http://192.168.1.5:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_GoogleRedirectURL(t *testing.T) {
	os.Unsetenv("BASE_URL")
	os.Setenv("GOOGLE_REDIRECT_URL", "https://custom.example.com/callback")
	defer os.Unsetenv("GOOGLE_REDIRECT_URL")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("localhost:8080")

	assert.Equal(t, "https://custom.example.com/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTPS(t *testing.T) {
	os.Unsetenv("BASE_URL")
	os.Unsetenv("GOOGLE_REDIRECT_URL")
	os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("localhost:8443")

	assert.Equal(t, "https://localhost:8443/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_DynamicHTTP(t *testing.T) {
	os.Unsetenv("BASE_URL")
	os.Unsetenv("GOOGLE_REDIRECT_URL")
	os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("localhost:8080")

	assert.Equal(t, "http://localhost:8080/auth/google/callback", got)
}

func TestRealGoogleProvider_getRedirectURL_Production(t *testing.T) {
	os.Unsetenv("BASE_URL")
	os.Unsetenv("GOOGLE_REDIRECT_URL")
	os.Setenv("AGBALUMO_ENV", "production")
	defer os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	got := p.getRedirectURL("agbalumo.fly.dev")

	assert.Equal(t, "https://agbalumo.fly.dev/auth/google/callback", got)
}

// --- RealGoogleProvider.GetAuthCodeURL Test ---

func TestRealGoogleProvider_GetAuthCodeURL(t *testing.T) {
	os.Unsetenv("BASE_URL")
	os.Unsetenv("GOOGLE_REDIRECT_URL")
	os.Unsetenv("AGBALUMO_ENV")

	p := NewRealGoogleProvider()
	url := p.GetAuthCodeURL("test-state", "localhost:8080")

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
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	// Override package-level var
	original := googleUserInfoURL
	googleUserInfoURL = server.URL
	defer func() { googleUserInfoURL = original }()

	p := NewRealGoogleProvider()
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
		fmt.Fprint(w, "server error")
	}))
	defer server.Close()

	original := googleUserInfoURL
	googleUserInfoURL = server.URL
	defer func() { googleUserInfoURL = original }()

	p := NewRealGoogleProvider()
	token := &oauth2.Token{AccessToken: "test-token"}
	user, err := p.GetUserInfo(t.Context(), token)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to fetch user info")
}

func TestRealGoogleProvider_GetUserInfo_BadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "not valid json{{{")
	}))
	defer server.Close()

	original := googleUserInfoURL
	googleUserInfoURL = server.URL
	defer func() { googleUserInfoURL = original }()

	p := NewRealGoogleProvider()
	token := &oauth2.Token{AccessToken: "test-token"}
	user, err := p.GetUserInfo(t.Context(), token)

	assert.Error(t, err)
	assert.Nil(t, user)
}
