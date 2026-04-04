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
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleLoginView(t *testing.T) {
	tests := []struct {
		name       string
		user       interface{}
		expectCode int
		expectLoc  string
	}{
		{
			name:       "AlreadyAdmin_RedirectsToDashboard",
			user:       domain.User{ID: "u1", Role: domain.UserRoleAdmin},
			expectCode: http.StatusTemporaryRedirect,
			expectLoc:  "/admin",
		},
		{
			name:       "NonAdmin_RendersLoginForm",
			user:       domain.User{ID: "u2", Role: "user"},
			expectCode: http.StatusOK,
		},
		{
			name:       "NoUser_RendersLoginForm",
			user:       nil,
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupAdminTestContext(http.MethodGet, "/admin/login", nil)
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			repo := testutil.SetupTestRepository(t)
			h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
			_ = h.HandleLoginView(c)

			assert.Equal(t, tt.expectCode, rec.Code)
			if tt.expectLoc != "" {
				assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
			}
		})
	}
}

func TestAdminHandler_HandleLoginAction(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		adminCode  string
		user       *domain.User
		expectCode int
		expectLoc  string
		verifyUser func(*testing.T, *admin.AdminHandler, string)
	}{
		{
			name:       "WrongCode_RendersError",
			code:       "wrong",
			adminCode:  "correct",
			user:       &domain.User{ID: "u1", Email: "test@example.com", Role: "user"},
			expectCode: http.StatusOK,
		},
		{
			name:       "NoUser_RedirectsToLogin",
			code:       "correct",
			adminCode:  "correct",
			user:       nil,
			expectCode: http.StatusTemporaryRedirect,
			expectLoc:  "/auth/google/login",
		},
		{
			name:       "ValidCode_PromotesUser",
			code:       "secret",
			adminCode:  "secret",
			user:       &domain.User{ID: "u1", Email: "admin@example.com", Role: "user"},
			expectCode: http.StatusFound,
			expectLoc:  "/admin",
			verifyUser: func(t *testing.T, h *admin.AdminHandler, userID string) {
				// We don't have a direct way to get user from handler easily without repo
				// but we can check the repo we passed in.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Set("code", tt.code)
			c, rec := setupAdminTestContext(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))

			repo := testutil.SetupTestRepository(t)
			if tt.user != nil {
				err := repo.SaveUser(context.Background(), *tt.user)
				assert.NoError(t, err)
				c.Set("User", *tt.user)
			}

			cfg := config.LoadConfig()
			cfg.AdminCode = tt.adminCode

			h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: cfg})
			_ = h.HandleLoginAction(c)

			assert.Equal(t, tt.expectCode, rec.Code)
			if tt.expectLoc != "" {
				assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
			}

			if tt.name == "ValidCode_PromotesUser" && tt.user != nil {
				updatedUser, err := repo.FindUserByID(context.Background(), tt.user.ID)
				assert.NoError(t, err)
				assert.Equal(t, domain.UserRoleAdmin, updatedUser.Role)
			}
		})
	}
}
