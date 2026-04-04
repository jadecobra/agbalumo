package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAddCategory_Success(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "Music")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	repo := testutil.SetupTestRepository(t)

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify database state
	cats, err := repo.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.Equal(t, "Music", cats[0].Name)
}

func TestAdminHandler_HandleAddCategory_EmptyName_Redirects(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "  ")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})

	repo := testutil.SetupTestRepository(t)
	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify no category added
	cats, _ := repo.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Empty(t, cats)
}

func TestAdminHandler_HandleAddCategory_DuplicateName_Redirects(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "Music")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	repo := testutil.SetupTestRepository(t)
	// Seed existing category
	_ = repo.SaveCategory(context.Background(), domain.CategoryData{ID: "music", Name: "Music"})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))

	// Verify still only one category
	cats, _ := repo.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Len(t, cats, 1)
}

func TestAdminHandler_HandleAddCategory_Claimable(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "Services")
	formData.Set("claimable", "true")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	repo := testutil.SetupTestRepository(t)

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify database state
	cats, _ := repo.GetCategories(context.Background(), domain.CategoryFilter{})
	assert.Len(t, cats, 1)
	assert.Equal(t, "Services", cats[0].Name)
	assert.True(t, cats[0].Claimable)
}
