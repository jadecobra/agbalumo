package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// -- Mock Provider for Handler Tests --
type MockGoogleProvider struct {
	GetAuthCodeURLFn func(state string, host string) string
	ExchangeFn       func(ctx context.Context, code string, host string) (*oauth2.Token, error)
	GetUserInfoFn    func(ctx context.Context, token *oauth2.Token) (*GoogleUser, error)
}

func (m *MockGoogleProvider) GetAuthCodeURL(state string, host string) string {
	return m.GetAuthCodeURLFn(state, host)
}

func (m *MockGoogleProvider) Exchange(ctx context.Context, code string, host string) (*oauth2.Token, error) {
	return m.ExchangeFn(ctx, code, host)
}

func (m *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	return m.GetUserInfoFn(ctx, token)
}

func TestRealGoogleProvider_GetUserInfo(t *testing.T) {
	// Mock the Google API Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("access_token")
		if token != "valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user := GoogleUser{
			ID:    "123",
			Email: "test@example.com",
			Name:  "Test User",
		}
		json.NewEncoder(w).Encode(user)
	}))
	defer ts.Close()

	// Override URL
	oldURL := googleUserInfoURL
	googleUserInfoURL = ts.URL
	defer func() { googleUserInfoURL = oldURL }()

	p := NewRealGoogleProvider()

	// Case 1: Success
	user, err := p.GetUserInfo(context.Background(), &oauth2.Token{AccessToken: "valid-token"})
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	// Case 2: Failure
	_, err = p.GetUserInfo(context.Background(), &oauth2.Token{AccessToken: "invalid-token"})
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	e := echo.New()

	// Setup Middleware for Session
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Manually inject session
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	// Expect lookup - return error to simulate not found
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, testifyMock.Anything).Return(domain.User{}, errors.New("not found"))
	// Expect create
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(nil)

	mockProvider := &MockGoogleProvider{
		ExchangeFn: func(ctx context.Context, code string, host string) (*oauth2.Token, error) {
			if code == "valid-code" {
				return &oauth2.Token{AccessToken: "tok"}, nil
			}
			return nil, errors.New("invalid code")
		},
		GetUserInfoFn: func(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
			return &GoogleUser{
				ID:    "g-123",
				Email: "new@example.com",
				Name:  "New User",
			}, nil
		},
	}

	h := &AuthHandler{
		Repo:           mockRepo,
		GoogleProvider: mockProvider,
	}

	// Run Handler
	if err := h.GoogleCallback(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockProvider := &MockGoogleProvider{
		GetAuthCodeURLFn: func(state string, host string) string {
			return "http://google.com/auth"
		},
	}
	h := &AuthHandler{GoogleProvider: mockProvider}

	if err := h.GoogleLogin(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "http://google.com/auth" {
		t.Errorf("Expected Location 'http://google.com/auth', got %s", loc)
	}
}

func TestRequireAuth(t *testing.T) {
	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	h := &AuthHandler{} // Repo not needed for redirect check

	// Case 1: No Session -> Redirect
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := h.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Protected")
	})

	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}

	// Case 2: Authenticated -> Next
	req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	session2, _ := store.Get(req2, "auth_session")
	session2.Values["user_id"] = "user123"
	c2.Set("session", session2)

	if err := handler(c2); err != nil {
		t.Fatal(err)
	}
	if rec2.Code != http.StatusOK {
		t.Errorf("Expected OK 200, got %d", rec2.Code)
	}
	if rec2.Body.String() != "Protected" {
		t.Errorf("Expected 'Protected', got %s", rec2.Body.String())
	}
}

func TestOptionalAuth(t *testing.T) {
	e := echo.New()
	store := customMiddleware.NewTestSessionStore()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindUserByID", testifyMock.Anything, "user123").Return(domain.User{ID: "user123", Name: "Test"}, nil)

	h := &AuthHandler{Repo: mockRepo}

	// Case 1: Authenticated -> User in Context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	session, _ := store.Get(req, "auth_session")
	session.Values["user_id"] = "user123"
	c.Set("session", session)

	handler := h.OptionalAuth(func(c echo.Context) error {
		user, ok := c.Get("User").(domain.User)
		if !ok || user.ID != "user123" {
			return c.String(http.StatusInternalServerError, "User not in context")
		}
		return c.String(http.StatusOK, "OK")
	})

	if err := handler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestAuthHandler_DevLogin(t *testing.T) {
	// Set Env to development
	os.Setenv("AGBALUMO_ENV", "development")
	defer os.Unsetenv("AGBALUMO_ENV")

	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, testifyMock.Anything).Return(domain.User{}, errors.New("not found"))
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := &AuthHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=dev@test.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	if err := h.DevLogin(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestAuthHandler_Logout(t *testing.T) {
	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	h := &AuthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set a session
	session, _ := store.Get(req, "auth_session")
	session.Values["user_id"] = "user123"
	session.Save(req, rec) // Setup cookie
	c.Set("session", session)

	if err := h.Logout(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}

	// Verify MaxAge is -1 (expired)
	if session.Options.MaxAge != -1 {
		t.Errorf("Expected MaxAge -1, got %d", session.Options.MaxAge)
	}
}

func TestRealGoogleProvider_GetAuthCodeURL(t *testing.T) {
	p := NewRealGoogleProvider()

	// Test localhost
	url1 := p.GetAuthCodeURL("state", "localhost:8443")
	if !strings.Contains(url1, "redirect_uri") {
		t.Error("Expected redirect_uri in URL")
	}

	// Test production env logic via Setenv if needed, but simple call is fine for coverage of the method body.
}

func TestRealGoogleProvider_Exchange(t *testing.T) {
	// Setup Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock token endpoint behavior
		// OAuth2 exchange makes a POST request
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			// Return a dummy token
			w.Write([]byte(`{"access_token": "mock-token", "token_type": "Bearer", "expires_in": 3600}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	// Override Endpoint
	// We need to modify google.Endpoint which is a global var in the library?
	// No, google.Endpoint is a constant/var in golang.org/x/oauth2/google
	// We can't easily modify it if it's not a var in OUR code.
	// But `NewRealGoogleProvider` uses `google.Endpoint`.
	// Let's check `auth.go`.
	// `Endpoint: google.Endpoint,`

	// Since we can't modify `google.Endpoint` easily without race conditions or if it's const,
	// maybe we can just construct a provider with custom config for this test?
	// But `NewRealGoogleProvider` is what we want to test.

	// WORKAROUND: We will test the Exchange logic by creating a RealGoogleProvider manually
	// or just accept we tested the wrapper logic.
	// `Exchange` method calls `p.config.Exchange`.
	// If `p.config` has the real endpoint, it will fail network call.
	// We can modify `p.config.Endpoint` AFTER creating it?

	p := NewRealGoogleProvider()
	p.config.Endpoint = oauth2.Endpoint{
		AuthURL:  ts.URL + "/auth",
		TokenURL: ts.URL + "/token",
	}

	// Case 1: Success (Network call to our mock server)
	// Exchange code for token
	token, err := p.Exchange(context.Background(), "any-code", "localhost")
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if token.AccessToken != "mock-token" {
		t.Errorf("Expected token mock-token, got %s", token.AccessToken)
	}
}

func TestAuthHandler_GoogleCallback_UserUpdate(t *testing.T) {
	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	existingUser := domain.User{
		ID:        "user-123",
		GoogleID:  "g-123",
		Email:     "test@example.com",
		Name:      "Old Name",
		AvatarURL: "http://old.pic",
		CreatedAt: time.Now(),
	}

	mockRepo := &mock.MockListingRepository{}
	// Return existing user
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g-123").Return(existingUser, nil)

	// Expect SaveUser to be called with UPDATED details
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.MatchedBy(func(u domain.User) bool {
		return u.Name == "New Name" && u.AvatarURL == "http://new.pic"
	})).Return(nil)

	mockProvider := &MockGoogleProvider{
		ExchangeFn: func(ctx context.Context, code string, host string) (*oauth2.Token, error) {
			return &oauth2.Token{AccessToken: "tok"}, nil
		},
		GetUserInfoFn: func(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
			return &GoogleUser{
				ID:      "g-123",
				Email:   "test@example.com",
				Name:    "New Name",
				Picture: "http://new.pic",
			}, nil
		},
	}

	h := &AuthHandler{
		Repo:           mockRepo,
		GoogleProvider: mockProvider,
	}

	if err := h.GoogleCallback(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestAuthHandler_GoogleCallback_CreateError(t *testing.T) {
	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindUserByGoogleID", testifyMock.Anything, "g-123").Return(domain.User{}, errors.New("not found"))
	// Simulate DB error on save
	mockRepo.On("SaveUser", testifyMock.Anything, testifyMock.Anything).Return(errors.New("db error"))

	mockProvider := &MockGoogleProvider{
		ExchangeFn: func(ctx context.Context, code string, host string) (*oauth2.Token, error) {
			return &oauth2.Token{AccessToken: "tok"}, nil
		},
		GetUserInfoFn: func(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
			return &GoogleUser{
				ID:    "g-123",
				Email: "new@example.com",
				Name:  "New User",
			}, nil
		},
	}

	h := &AuthHandler{
		Repo:           mockRepo,
		GoogleProvider: mockProvider,
	}

	// Echo handler returns nil error when c.String successfully writes response
	_ = h.GoogleCallback(c)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rec.Code)
	}

	mockRepo.AssertExpectations(t)
}
