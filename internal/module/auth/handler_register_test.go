package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestAuthHandler_GoogleCallback_SaveUserError(t *testing.T) {
	c, _ := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{ID: "g-err", Email: "err@test.com"}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
}

func TestAuthHandler_GoogleCallback_UpdateProfile(t *testing.T) {
	c, rec := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

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
	_ = app.DB.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	updatedUser, _ := app.DB.FindUserByGoogleID(context.Background(), "g1")
	assert.Equal(t, "New Name", updatedUser.Name)
}

func TestAuthHandler_GoogleCallback_UpdateProfileSaveError(t *testing.T) {
	c, rec := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

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
	_ = app.DB.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_GoogleCallback_UpdateProfile_NoChanges(t *testing.T) {
	c, rec := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

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
	_ = app.DB.SaveUser(context.Background(), existingUser)

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	err := h.GoogleCallback(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_GoogleCallback_CrossSiteCallback(t *testing.T) {
	c, rec := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "random-state"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true
	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

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
