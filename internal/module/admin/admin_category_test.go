package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAddCategory_Success(t *testing.T) {
	t.Parallel()
	formData := url.Values{}
	formData.Set("name", "Music")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := admin.NewAdminHandler(app)
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify database state
	cats, err := app.DB.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.Equal(t, "Music", cats[0].Name)
}

func TestAdminHandler_HandleAddCategory_EmptyName_Redirects(t *testing.T) {
	t.Parallel()
	formData := url.Values{}
	formData.Set("name", "  ")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify no category added
	cats, _ := app.DB.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Empty(t, cats)
}

func TestAdminHandler_HandleAddCategory_DuplicateName_Redirects(t *testing.T) {
	t.Parallel()
	formData := url.Values{}
	formData.Set("name", "Music")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	// Seed existing category
	_ = app.DB.SaveCategory(context.Background(), domain.CategoryData{ID: "music", Name: "Music"})

	h := admin.NewAdminHandler(app)
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify still only one category
	cats, _ := app.DB.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Len(t, cats, 1)
}

func TestAdminHandler_HandleAddCategory_Claimable(t *testing.T) {
	t.Parallel()
	formData := url.Values{}
	formData.Set("name", "Services")
	formData.Set("claimable", "true")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := admin.NewAdminHandler(app)
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify database state
	cats, _ := app.DB.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Len(t, cats, 1)
	assert.Equal(t, "Services", cats[0].Name)
	assert.True(t, cats[0].Claimable)
}
