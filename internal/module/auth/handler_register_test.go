package auth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_GoogleCallback_SaveUserError(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	performRegistration(t, app, map[string]string{
		"id":    "g-err",
		"email": "err@test.com",
	})
}

func TestAuthHandler_GoogleCallback_UpdateProfile(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	rec := setupExistingUserAndRegister(t, app, domain.User{
		ID:        "u1",
		GoogleID:  "g1",
		Email:     "test@example.com",
		Name:      "Old Name",
		AvatarURL: "http://old-pic.com",
	}, map[string]string{
		"id":      "g1",
		"email":   "test@example.com",
		"name":    "New Name",
		"picture": "http://new-pic.com",
	})

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)

	updatedUser, _ := app.DB.FindUserByGoogleID(context.Background(), "g1")
	assert.Equal(t, "New Name", updatedUser.Name)
}

func TestAuthHandler_GoogleCallback_UpdateProfileSaveError(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	rec := setupExistingUserAndRegister(t, app, domain.User{
		ID:        "u-update-err",
		GoogleID:  "g-update-err",
		Email:     "user@test.com",
		Name:      "Old Name",
		AvatarURL: "http://old-pic.com",
	}, map[string]string{
		"id":      "g-update-err",
		"email":   "user@test.com",
		"name":    "New Name",
		"picture": "http://new-pic.com",
	})

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_GoogleCallback_UpdateProfile_NoChanges(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	rec := setupExistingUserAndRegister(t, app, domain.User{
		ID:        "u1",
		GoogleID:  "g1",
		Email:     "test@example.com",
		Name:      "Same Name",
		AvatarURL: "http://same-pic.com",
	}, map[string]string{
		"id":      "g1",
		"email":   "test@example.com",
		"name":    "Same Name",
		"picture": "http://same-pic.com",
	})
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAuthHandler_GoogleCallback_CrossSiteCallback(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.HasGoogleAuth = true

	rec := performRegistration(t, app, map[string]string{
		"id":      "google-cross-site",
		"email":   "cross@example.com",
		"name":    "Cross Site",
		"picture": "http://pic.com",
	})

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/", rec.Header().Get("Location"))
}

