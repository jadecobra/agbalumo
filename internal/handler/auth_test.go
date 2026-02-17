package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
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
	// We need a store to prevent "Session Store Missing" error
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Manually inject session (Middleware doesn't run on direct handler call)
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{
		FindUserByGoogleIDFn: func(ctx context.Context, googleID string) (domain.User, error) {
			return domain.User{}, errors.New("not found") // Force create
		},
		SaveUserFn: func(ctx context.Context, u domain.User) error {
			return nil
		},
	}

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

	// Middleware Setup
	handler := h.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Protected")
	})

	// Run
	// Directly calling handler(c) means middleware logic runs.
	// But RequireAuth needs session from context.
	// So we must manually inject session or run session middleware.
	// Let's manually inject.
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

	mockRepo := &mock.MockListingRepository{
		FindUserByIDFn: func(ctx context.Context, id string) (domain.User, error) {
			if id == "user123" {
				return domain.User{ID: "user123", Name: "Test"}, nil
			}
			return domain.User{}, errors.New("not found")
		},
	}
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
}

func TestAuthHandler_DevLogin(t *testing.T) {
	// Set Env to development
	os.Setenv("AGBALUMO_ENV", "development")
	defer os.Unsetenv("AGBALUMO_ENV")

	e := echo.New()
	store := customMiddleware.NewTestSessionStore()
	e.Use(customMiddleware.SessionMiddleware(store))

	mockRepo := &mock.MockListingRepository{
		FindUserByGoogleIDFn: func(ctx context.Context, id string) (domain.User, error) {
			return domain.User{}, errors.New("not found")
		},
		SaveUserFn: func(ctx context.Context, u domain.User) error {
			return nil
		},
	}
	h := &AuthHandler{Repo: mockRepo}

	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=dev@test.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Manually inject session
	session, _ := store.Get(req, "auth_session")
	c.Set("session", session)

	if err := h.DevLogin(c); err != nil {
		t.Fatalf("Handler failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
}
