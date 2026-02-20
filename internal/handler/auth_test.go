package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// MockGoogleProvider
type MockGoogleProvider struct {
	testifyMock.Mock
}

func (m *MockGoogleProvider) GetAuthCodeURL(state string, host string) string {
	args := m.Called(state, host)
	return args.String(0)
}

func (m *MockGoogleProvider) Exchange(ctx context.Context, code string, host string) (*oauth2.Token, error) {
	args := m.Called(ctx, code, host)
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

	mockRepo := &mock.MockListingRepository{}
	// Mock findOrCreateUser if needed, but DevLogin only calls it if env != production.
	// Since we set env=production, it should return 403 before calling repo.

	cfg := config.LoadConfig()
	cfg.Env = "production"
	h := handler.NewAuthHandler(mockRepo, nil, cfg)

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

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

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

	// Setup Session Store in Context
	store := sessions.NewCookieStore([]byte("secret"))
	// We need to inject the session into the request context as middleware would
	// But handler calls customMiddleware.GetSession(c).
	// Let's rely on injecting the store into the context if possible, or just mock the session logic?
	// The handler uses: session, _ := store.Get(c.Request(), "session-name") via echo.Context?
	// No, it uses customMiddleware.GetSession(c).
	// We must register the middleware on the echo instance or manually set the session value.

	// Better: Register the session middleware on the Echo instance and use e.NewContext
	// But e.NewContext does not run middleware.
	// We have to invoke the middleware chain.

	// Alternative: Just unit test the handler logic?
	// The handler calls `h.setSessionAndRedirect`.
	// If session is missing, it returns 500.

	// To make GetSession work, we need "session" in c.Get("session").
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &handler.GoogleUser{
		ID:      "google-123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "http://pic.com",
	}
	user := domain.User{
		ID:        "user-123",
		GoogleID:  "google-123",
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "http://pic.com",
		CreatedAt: time.Now(),
	}

	// Mocks
	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "google-123").Return(user, nil)

	// Execute
	err := h.GoogleCallback(c)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.Equal(t, "user-123", sess.Values["user_id"])
}

func TestAuthHandler_RequireAuth_Redirect(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Session (Empty)
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	handlerFunc := h.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/google/login", rec.Header().Get("Location"))
}

func TestAuthHandler_RequireAuth_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Session (With UserID)
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["user_id"] = "user-123"
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	handlerFunc := h.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	mockProvider.On("GetAuthCodeURL", "random-state", testifyMock.Anything).Return("http://google.com/auth")

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
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	err := h.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
	assert.Equal(t, -1, sess.Options.MaxAge)
}

func TestAuthHandler_OptionalAuth(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["user_id"] = "user-123"
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	user := domain.User{ID: "user-123"}
	mockRepo.On("FindUserByID", testifyMock.Anything, "user-123").Return(user, nil)

	handlerFunc := h.OptionalAuth(func(c echo.Context) error {
		u := c.Get("User")
		if u == nil {
			return c.String(http.StatusOK, "No User")
		}
		return c.String(http.StatusOK, "Has User")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, "Has User", rec.Body.String())
}

func TestAuthHandler_DevLogin_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// Session Middleware
	store := sessions.NewCookieStore([]byte("secret"))
	// Mock retrieval via context injection
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=test@dev.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	// Environment Development
	os.Setenv("AGBALUMO_ENV", "development")
	defer os.Unsetenv("AGBALUMO_ENV")

	// Mock findOrCreateUser logic by mocking FindUserByGoogleID
	// It will try to find "dev-test@dev.com"
	// If not found, it calls SaveUser.

	// Scenario: New User
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "dev-test@dev.com").Return(domain.User{}, assert.AnError) // Return error to trigger creation
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.AnythingOfType("domain.User")).Return(nil)

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
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{ID: "g-err", Email: "err@test.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	// Mock DB Error
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g-err").Return(domain.User{}, assert.AnError) // Not found
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)             // Save fails

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAuthHandler_GoogleCallback_UpdateProfile(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{
		ID:      "g1",
		Email:   "test@example.com",
		Name:    "New Name",
		Picture: "http://new-pic.com",
	}

	// Existing user with OLD details
	existingUser := domain.User{
		ID:        "u1",
		GoogleID:  "g1",
		Email:     "test@example.com",
		Name:      "Old Name",
		AvatarURL: "http://old-pic.com",
	}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g1").Return(existingUser, nil)
	// Expect SaveUser with NEW details
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.MatchedBy(func(u domain.User) bool {
		return u.Name == "New Name" && u.AvatarURL == "http://new-pic.com"
	})).Return(nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

// --- DevLogin Edge Cases ---

func TestAuthHandler_DevLogin_GOENVFallback(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=go@env.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	// Only set GO_ENV, NOT AGBALUMO_ENV
	os.Unsetenv("AGBALUMO_ENV")
	os.Setenv("GO_ENV", "development")
	defer os.Unsetenv("GO_ENV")

	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "dev-go@env.com").Return(domain.User{}, assert.AnError)
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.AnythingOfType("domain.User")).Return(nil)

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.NotEmpty(t, sess.Values["user_id"])
}

func TestAuthHandler_DevLogin_DefaultEmail(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	store := sessions.NewCookieStore([]byte("secret"))
	// No email param
	req := httptest.NewRequest(http.MethodGet, "/auth/dev", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	os.Setenv("AGBALUMO_ENV", "development")
	defer os.Unsetenv("AGBALUMO_ENV")

	// Default email is "dev@agbalumo.com", so googleID is "dev-dev@agbalumo.com"
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "dev-dev@agbalumo.com").Return(domain.User{}, assert.AnError)
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.AnythingOfType("domain.User")).Return(nil)

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_DevLogin_FindOrCreateError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=err@test.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	os.Setenv("AGBALUMO_ENV", "development")
	defer os.Unsetenv("AGBALUMO_ENV")

	// findOrCreateUser: user not found, then save fails
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "dev-err@test.com").Return(domain.User{}, assert.AnError)
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)

	err := h.DevLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

// --- GoogleCallback Error Paths ---

func TestAuthHandler_GoogleCallback_ExchangeError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=bad-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	mockProvider.On("Exchange", testifyMock.Anything, "bad-code", testifyMock.Anything).Return(nil, assert.AnError)

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

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(nil, assert.AnError)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

// --- findOrCreateUser: Profile Update Save Error (graceful) ---

func TestAuthHandler_GoogleCallback_UpdateProfileSaveError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

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

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g-update-err").Return(existingUser, nil)
	// SaveUser fails during profile update — should log but NOT fail
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)

	err := h.GoogleCallback(c)

	// Should still succeed (graceful error handling)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}

// --- setSessionAndRedirect: Nil Session ---

func TestAuthHandler_SetSessionAndRedirect_NilSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No session set — session will be nil

	mockRepo := &mock.MockListingRepository{}
	mockProvider := &MockGoogleProvider{}
	h := handler.NewAuthHandler(mockRepo, mockProvider, config.LoadConfig())

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &handler.GoogleUser{ID: "g-no-session", Email: "no-session@test.com", Name: "Test", Picture: "http://pic.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", testifyMock.Anything).Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g-no-session").Return(domain.User{}, assert.AnError)
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error Page")
}

// --- OptionalAuth: No Session ---

func TestAuthHandler_OptionalAuth_NoSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No session set at all

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	handlerFunc := h.OptionalAuth(func(c echo.Context) error {
		u := c.Get("User")
		if u == nil {
			return c.String(http.StatusOK, "No User")
		}
		return c.String(http.StatusOK, "Has User")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, "No User", rec.Body.String())
}

// --- Logout: No Session ---

func TestAuthHandler_Logout_NoSession(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No session set

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAuthHandler(mockRepo, nil, config.LoadConfig())

	err := h.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}
