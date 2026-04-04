package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestAuthHandler_GoogleCallback_SaveUserError(t *testing.T) {
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
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{ID: "g-err", Email: "err@test.com"}

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

	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := testutil.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{
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

func TestAuthHandler_GoogleCallback_UpdateProfileSaveError(t *testing.T) {
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
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{
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

func TestAuthHandler_GoogleCallback_UpdateProfile_NoChanges(t *testing.T) {
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
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: config.LoadConfig()})

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{
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

func TestAuthHandler_GoogleCallback_CrossSiteCallback(t *testing.T) {
	// Simulate Google callback where the main session cookie is dropped due to SameSite=StrictMode
	// but the custom 'oauth_state' cookie is preserved because it's SameSite=LaxMode.
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code", nil)

	// ONLY set the oauth_state cookie, NOT the main session cookie
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Inject a fresh, empty session to simulate what SessionMiddleware does when the strict cookie is missing cross-site.
	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := testutil.SetupTestRepository(t)
	mockProvider := &MockGoogleProvider{}
	cfg := config.LoadConfig()
	cfg.HasGoogleAuth = true
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, GoogleProvider: mockProvider, Config: cfg})

	token := &oauth2.Token{AccessToken: "access-token"}
	gUser := &auth.GoogleUser{
		ID:      "google-cross-site",
		Email:   "cross@example.com",
		Name:    "Cross Site",
		Picture: "http://pic.com",
	}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}
