package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthHandler_GoogleCallback_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := testutil.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: cfg})

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &auth.GoogleUser{ID: "google-123", Email: "test@example.com", Name: "Test User", Picture: "http://pic.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	user, err := repo.FindUserByGoogleID(context.Background(), "google-123")
	require.NoError(t, err)
	assert.Equal(t, user.ID, sess.Values["user_id"])
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := testutil.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: cfg})

	mockProvider.On("GetAuthCodeURL", testifyMock.AnythingOfType("string"), "http", "example.com").Return("http://google.com/auth")

	err := h.GoogleLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "http://google.com/auth", rec.Header().Get("Location"))
}
