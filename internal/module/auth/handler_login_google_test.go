package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleCallback_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true

	rec := performRegistration(t, app, map[string]string{
		"id":      "google-123",
		"email":   "test@example.com",
		"name":    "Test User",
		"picture": "http://pic.com",
	})

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	user, err := app.DB.FindUserByGoogleID(context.Background(), "google-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	// We can't easily check the session here because the helper doesn't expose the echo context
	// but we've verified the redirection and DB persistence.
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupTestContextWithSession(http.MethodGet, "/auth/google/login", nil)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true
	mockProvider := &testutil.MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	mockProvider.On("GetAuthCodeURL", testifyMock.AnythingOfType("string"), "http", "example.com").Return("http://google.com/auth")

	err := h.GoogleLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "http://google.com/auth", rec.Header().Get("Location"))
}

