package module

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/labstack/echo/v4"
	"net/http"
)

// BaseViewData represents the common data required by all page templates.
type BaseViewData struct {
	User          interface{}
	Counts        map[string]int
	Env           string
	CSRF          string
	Categories    []domain.CategoryData
	DevMode       bool
	HasGoogleAuth bool
}

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

// PopulateBase extracts common data (Categories, Env, etc.) from the context
// and returns a BaseViewData struct.
func (h *BaseHandler) PopulateBase(c echo.Context) BaseViewData {
	data := BaseViewData{
		Env:           h.App.Cfg.Env,
		DevMode:       h.App.Cfg.Env == "development",
		HasGoogleAuth: h.App.Cfg.HasGoogleAuth,
	}

	categories, err := h.App.CategorizationSvc.GetActiveCategories(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("Failed to retrieve categories: %v", err)
	} else {
		data.Categories = categories
	}

	counts, err := h.App.DB.GetCounts(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("Failed to retrieve counts: %v", err)
		data.Counts = map[string]int{}
	} else {
		strCounts := make(map[string]int)
		for cat, count := range counts {
			strCounts[string(cat)] = count
		}
		data.Counts = strCounts
	}

	if u := c.Get("User"); u != nil {
		data.User = u
	}

	if csrf := c.Get("csrf"); csrf != nil {
		if s, ok := csrf.(string); ok {
			data.CSRF = s
		}
	}

	return data
}

// RenderTyped is the preferred way to render templates with typed ViewModels.
func (h *BaseHandler) RenderTyped(c echo.Context, tmpl string, data interface{}) error {
	return c.Render(http.StatusOK, tmpl, data)
}

// RenderWithBaseContext is a shared helper that injects common data (Categories, Env, etc.)
// into the data map before rendering.
// Deprecated: Use PopulateBase() + typed ViewModel structs.
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

	if csrf := c.Get("csrf"); csrf != nil {
		data["CSRF"] = csrf
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
