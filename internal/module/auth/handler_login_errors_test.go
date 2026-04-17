package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleCallback_Errors(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	t.Run("StateMismatch", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=wrong", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "valid"})


		app, cleanup := testutil.SetupTestAppEnv(t)
		defer cleanup()
		h := auth.NewAuthHandler(app)
		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ExchangeError", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=s&code=c", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "s"})


		app, cleanup := testutil.SetupTestAppEnv(t)
		defer cleanup()
		mockProvider := &MockGoogleProvider{}
		h := auth.NewAuthHandler(app)
		h.GoogleProvider = mockProvider
		mockProvider.On("Exchange", testifyMock.Anything, "c", "http", "example.com").Return(nil, assert.AnError)

		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
