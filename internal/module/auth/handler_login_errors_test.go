package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleCallback_Errors(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	t.Run("StateMismatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=wrong", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "valid"})

		repo := testutil.SetupTestRepository(t)
		h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: config.LoadConfig()})
		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ExchangeError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=s&code=c", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "s"})

		repo := testutil.SetupTestRepository(t)
		mockProvider := &MockGoogleProvider{}
		h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})
		mockProvider.On("Exchange", testifyMock.Anything, "c", "http", "example.com").Return(nil, assert.AnError)

		err := h.GoogleCallback(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
