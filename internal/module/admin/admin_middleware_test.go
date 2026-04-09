package admin_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware_NoUser(t *testing.T) {
	t.Parallel()
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	called := false
	mdw := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})

	_ = mdw(c)
	assert.False(t, called)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestAdminMiddleware_AdminUser(t *testing.T) {
	t.Parallel()
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	called := false
	mdw := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return c.String(http.StatusOK, "ok")
	})

	_ = mdw(c)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
}
