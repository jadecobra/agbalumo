package auth_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleCallback_Errors(t *testing.T) {
	t.Parallel()

	t.Run("StateMismatch", func(t *testing.T) {
		t.Parallel()
		c, rec := testutil.SetupTestContextWithSession(http.MethodGet, "/auth/google/callback?state=wrong", nil)
		c.Request().AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "valid"})

		app, cleanup := testutil.SetupTestAppEnv(t)
		defer cleanup()
		h := auth.NewAuthHandler(app)
		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ExchangeError", func(t *testing.T) {
		t.Parallel()
		c, rec := testutil.SetupTestContextWithSession(http.MethodGet, "/auth/google/callback?state=s&code=c", nil)
		c.Request().AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "s"})

		app, cleanup := testutil.SetupTestAppEnv(t)
		defer cleanup()
		mockProvider := &testutil.MockGoogleProvider{}
		h := auth.NewAuthHandler(app)
		h.GoogleProvider = mockProvider
		mockProvider.On("Exchange", testifyMock.Anything, "c", "http", "example.com").Return(nil, assert.AnError)

		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

