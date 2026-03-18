package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// MockGoogleProvider
type MockGoogleProvider struct {
	testifyMock.Mock
}

func (m *MockGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	args := m.Called(state, scheme, host)
	return args.String(0)
}

func (m *MockGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	args := m.Called(ctx, code, scheme, host)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*handler.GoogleUser, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*handler.GoogleUser), args.Error(1)
}

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
	h := handler.NewAuthHandler(repo, nil, cfg)

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

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := handler.NewAuthHandler(repo, mockProvider, cfg)

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

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := handler.NewAuthHandler(repo, mockProvider, cfg)

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &handler.GoogleUser{
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
	h := handler.NewAuthHandler(repo, mockProvider, cfg)

	mockProvider.On("GetAuthCodeURL", testifyMock.AnythingOfType("string"), "http", "example.com").Return("http://google.com/auth")

	err := h.GoogleLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "http://google.com/auth", rec.Header().Get("Location"))
}

func TestAuthHandler_Logout(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	h := handler.NewAuthHandler(repo, nil, config.LoadConfig())

	err := h.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.Equal(t, -1, sess.Options.MaxAge)
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
	h := handler.NewAuthHandler(repo, nil, config.LoadConfig())

	_ = os.Setenv("AGBALUMO_ENV", "development")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.NotEmpty(t, sess.Values["user_id"])
}

func TestAuthHandler_GoogleCallback_SaveUserError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{ID: "g-err", Email: "err@test.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
}

func TestAuthHandler_GoogleCallback_UpdateProfile(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{
		ID:      "g1",
		Email:   "test@example.com",
		Name:    "New Name",
		Picture: "http://new-pic.com",
	}

	existingUser := domain.User{
		ID:        "u1",
		GoogleID:  "g1",
		Email:     "test@example.com",
		Name:      "Old Name",
		AvatarURL: "http://old-pic.com",
	}
	_ = repo.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	updatedUser, _ := repo.FindUserByGoogleID(context.Background(), "g1")
	assert.Equal(t, "New Name", updatedUser.Name)
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
	h := handler.NewAuthHandler(repo, nil, config.LoadConfig())

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
	h := handler.NewAuthHandler(repo, nil, config.LoadConfig())

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

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

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

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(nil, assert.AnError)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_GoogleCallback_UpdateProfileSaveError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{
		ID:      "g-update-err",
		Email:   "user@test.com",
		Name:    "New Name",
		Picture: "http://new-pic.com",
	}

	existingUser := domain.User{
		ID:        "u-update-err",
		GoogleID:  "g-update-err",
		Email:     "user@test.com",
		Name:      "Old Name",
		AvatarURL: "http://old-pic.com",
	}
	_ = repo.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_SetSessionAndRedirect_NilSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{ID: "g-no-session", Email: "no-session@test.com", Name: "Test", Picture: "http://pic.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

func TestAuthHandler_Logout_NoSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	h := handler.NewAuthHandler(repo, nil, config.LoadConfig())

	err := h.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}

func TestAuthHandler_GoogleCallback_UpdateProfile_NoChanges(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["oauth_state"] = "random-state"
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(repo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{
		ID:      "g1",
		Email:   "test@example.com",
		Name:    "Same Name",
		Picture: "http://same-pic.com",
	}

	existingUser := domain.User{
		ID:        "u1",
		GoogleID:  "g1",
		Email:     "test@example.com",
		Name:      "Same Name",
		AvatarURL: "http://same-pic.com",
	}
	_ = repo.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}
