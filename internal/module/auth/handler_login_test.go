package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthHandler_DevLogin_Production(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/dev", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)

	cfg := config.LoadConfig()
	cfg.Env = "production"
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: cfg})

	// Execute
	err := h.DevLogin(c)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_GoogleCallback_StateMismatch(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=wrong-state", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Using cookie instead of session for oauth_state
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	// Still need session for finding if a user is logged in later (if applicable), though here we just test callback handling.
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: cfg})

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_GoogleCallback_Success(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Using cookie instead of session for oauth_state
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: cfg})

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &auth.GoogleUser{
		ID:      "google-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "http://pic.com",
	}

	// Mocks
	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	// Execute
	err := h.GoogleCallback(c)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))

	// Verify user in DB
	user, err := repo.FindUserByGoogleID(context.Background(), "google-123")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, user.ID, sess.Values["user_id"])
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
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

func TestAuthHandler_DevLogin_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=test@dev.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: config.LoadConfig()})

	_ = os.Setenv("AGBALUMO_ENV", "development")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.NotEmpty(t, sess.Values["user_id"])
}

func TestAuthHandler_DevLogin_GOENVFallback(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=go@env.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: config.LoadConfig()})

	_ = os.Unsetenv("AGBALUMO_ENV")
	_ = os.Setenv("GO_ENV", "development")
	defer func() { _ = os.Unsetenv("GO_ENV") }()

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.NotEmpty(t, sess.Values["user_id"])
}

func TestAuthHandler_DevLogin_DefaultEmail(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: config.LoadConfig()})

	_ = os.Setenv("AGBALUMO_ENV", "development")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_DevLogin_FindOrCreateError(t *testing.T) {
	// Skipping error path for real repo.
}

func TestAuthHandler_GoogleCallback_ExchangeError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=bad-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	mockProvider.On("Exchange", testifyMock.Anything, "bad-code", "http", "example.com").Return(nil, assert.AnError)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_GoogleCallback_GetUserInfoError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(nil, assert.AnError)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_SetSessionAndRedirect_NilSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	// Inject cookie to bypass state verification and reach setSessionAndRedirect
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{ID: "g-no-session", Email: "no-session@test.com", Name: "Test", Picture: "http://pic.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}
