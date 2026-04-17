package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// performRegistration is a helper to simulate a successful Google OAuth callback and registration.
func performRegistration(t *testing.T, app *env.AppEnv, payload map[string]string) *httptest.ResponseRecorder {
	c, rec := testutil.SetupTestContextWithSession(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "random-state"})

	mockProvider := &testutil.MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{
		ID:      payload["id"],
		Email:   payload["email"],
		Name:    payload["name"],
		Picture: payload["picture"],
	}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	_ = h.GoogleCallback(c)

	return rec
}

// setupExistingUserAndRegister creates a user in the DB and then performs a simulated registration.
func setupExistingUserAndRegister(t *testing.T, app *env.AppEnv, user domain.User, payload map[string]string) *httptest.ResponseRecorder {
	_ = app.DB.SaveUser(context.Background(), user)
	return performRegistration(t, app, payload)
}
