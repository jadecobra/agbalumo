package handler_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware_NoUser(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	h := handler.NewAdminHandler(nil, nil, nil, nil, nil, nil, nil, nil, config.LoadConfig())

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
	c, rec := setupAdminTestContext(http.MethodGet, "/admin", nil)
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
	h := handler.NewAdminHandler(nil, nil, nil, nil, nil, nil, nil, nil, config.LoadConfig())

	called := false
	mdw := h.AdminMiddleware(func(c echo.Context) error {
		called = true
		return c.String(http.StatusOK, "ok")
	})

	_ = mdw(c)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)
}
