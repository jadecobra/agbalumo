package module

import (
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/labstack/echo/v4"
	"net/http"
)

// BaseHandler provides shared dependencies and utilities for all module handlers.
type BaseHandler struct {
	App *env.AppEnv
}

// LogError is a shared helper to log errors through the echo context logger.
func (h *BaseHandler) LogError(c echo.Context, msg string, err error) {
	if err != nil {
		c.Logger().Errorf("%s: %v", msg, err)
	}
}

// RenderWithBaseContext is a shared helper that injects common data (Categories, Env, etc.) 
// into the data map before rendering.
func (h *BaseHandler) RenderWithBaseContext(c echo.Context, tmpl string, data map[string]interface{}) error {
	ctx := c.Request().Context()

	// Fetch categories if not already provided
	if _, exists := data["Categories"]; !exists {
		categories, err := h.App.CategorizationSvc.GetActiveCategories(ctx)
		if err != nil {
			c.Logger().Errorf("Failed to retrieve categories: %v", err)
			data["Categories"] = []interface{}{}
		} else {
			data["Categories"] = categories
		}
	}

	data["Env"] = h.App.Cfg.Env
	data["HasGoogleAuth"] = h.App.Cfg.HasGoogleAuth
	
	// Add User if present in context
	if u := c.Get("User"); u != nil {
		data["User"] = u
	}

	return c.Render(http.StatusOK, tmpl, data)
}
