package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestAuthHandler_GoogleCallback_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &auth.GoogleUser{ID: "google-123", Email: "test@example.com", Name: "Test User", Picture: "http://pic.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	user, err := app.DB.FindUserByGoogleID(context.Background(), "google-123")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, sess.Values["user_id"])
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	mockProvider.On("GetAuthCodeURL", testifyMock.AnythingOfType("string"), "http", "example.com").Return("http://google.com/auth")

	err := h.GoogleLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "http://google.com/auth", rec.Header().Get("Location"))
}
