package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
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

func TestAuthHandler_GoogleCallback_ProfileUpdates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		googleUser   map[string]string
		check        func(t *testing.T, app *env.AppEnv, rec *httptest.ResponseRecorder)
		existingUser domain.User
		name         string
	}{
		{
			name: "Update Name and Avatar",
			existingUser: domain.User{
				ID: "u1", GoogleID: "g1", Email: "test@example.com",
				Name: "Old Name", AvatarURL: "http://old-pic.com",
			},
			googleUser: map[string]string{
				"id": "g1", "email": "test@example.com",
				"name": "New Name", "picture": "http://new-pic.com",
			},
			check: func(t *testing.T, app *env.AppEnv, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
				updatedUser, _ := app.DB.FindUserByGoogleID(context.Background(), "g1")
				assert.Equal(t, "New Name", updatedUser.Name)
			},
		},
		{
			name: "No Changes Needed",
			existingUser: domain.User{
				ID: "u2", GoogleID: "g2", Email: "test2@example.com",
				Name: "Same Name", AvatarURL: "http://same-pic.com",
			},
			googleUser: map[string]string{
				"id": "g2", "email": "test2@example.com",
				"name": "Same Name", "picture": "http://same-pic.com",
			},
			check: func(t *testing.T, app *env.AppEnv, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
			},
		},
		{
			name: "Save Error Resilience",
			existingUser: domain.User{
				ID: "u-err", GoogleID: "g-err", Email: "err@test.com",
			},
			googleUser: map[string]string{
				"id": "g-err", "email": "err@test.com", "name": "Error Person",
			},
			check: func(t *testing.T, app *env.AppEnv, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app, cleanup := testutil.SetupTestAppEnv(t)
			defer cleanup()

			rec := setupExistingUserAndRegister(t, app, tt.existingUser, tt.googleUser)
			tt.check(t, app, rec)
		})
	}
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
