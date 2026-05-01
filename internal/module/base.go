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

	h.injectCategories(c, data)
	h.injectCounts(c, data)

	data["Env"] = h.App.Cfg.Env
	data["DevMode"] = h.App.Cfg.Env == "development"
	data["HasGoogleAuth"] = h.App.Cfg.HasGoogleAuth

	// Add User if present in context
	if u := c.Get("User"); u != nil {
		data["User"] = u
	}

	return c.Render(http.StatusOK, tmpl, data)
}

func (h *BaseHandler) injectCategories(c echo.Context, data map[string]interface{}) {
	if _, exists := data["Categories"]; exists {
		return
	}

	categories, err := h.App.CategorizationSvc.GetActiveCategories(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("Failed to retrieve categories: %v", err)
		data["Categories"] = []interface{}{}
		return
	}
	data["Categories"] = categories
}

func (h *BaseHandler) injectCounts(c echo.Context, data map[string]interface{}) {
	if _, exists := data["Counts"]; exists {
		return
	}

	counts, err := h.App.DB.GetCounts(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("Failed to retrieve counts: %v", err)
		data["Counts"] = map[string]int{}
		return
	}

	strCounts := make(map[string]int)
	for cat, count := range counts {
		strCounts[string(cat)] = count
	}
	data["Counts"] = strCounts
}
