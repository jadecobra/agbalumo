package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleLoginView(t *testing.T) {
	t.Parallel()
	tests := []struct {
		user       interface{}
		name       string
		expectLoc  string
		expectCode int
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
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := admin.NewAdminHandler(env.App)

			c, rec := testutil.SetupModuleContext(http.MethodGet, "/admin/login", nil)
			if tt.user != nil {
		c.Set(domain.CtxKeyUser, tt.user)

			}

			_ = h.HandleLoginView(c)

			assert.Equal(t, tt.expectCode, rec.Code)
			if tt.expectLoc != "" {
				assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
			}
		})
	}
}

func TestAdminHandler_HandleLoginAction(t *testing.T) {
	t.Parallel()
	tests := []loginActionTest{
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
	runAdminLoginActionTests(t, tests)
}

// runAdminLoginActionTests exists solely to reduce the cognitive complexity of the test suite
type loginActionTest struct {
	user       *domain.User
	verifyUser func(*testing.T, *admin.AdminHandler, string)
	name       string
	code       string
	adminCode  string
	expectLoc  string
	expectCode int
}

func runAdminLoginActionTests(t *testing.T, tests []loginActionTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			runSingleAdminLoginActionTest(t, tt)
		})
	}
}

func runSingleAdminLoginActionTest(t *testing.T, tt loginActionTest) {
	formData := url.Values{}
	formData.Set("code", tt.code)
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/admin/login", strings.NewReader(formData.Encode()))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	if tt.user != nil {
		err := env.App.DB.SaveUser(context.Background(), *tt.user)
		assert.NoError(t, err)
		c.Set(domain.CtxKeyUser, *tt.user)

	}

	cfg := env.App.Cfg
	cfg.AdminCode = tt.adminCode

	h := admin.NewAdminHandler(env.App)
	_ = h.HandleLoginAction(c)

	assert.Equal(t, tt.expectCode, rec.Code)
	if tt.expectLoc != "" {
		assert.Equal(t, tt.expectLoc, rec.Header().Get("Location"))
	}

	if tt.name == "ValidCode_PromotesUser" && tt.user != nil {
		updatedUser, err := env.App.DB.FindUserByID(context.Background(), tt.user.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.UserRoleAdmin, updatedUser.Role)
	}
}
