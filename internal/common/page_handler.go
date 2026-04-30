package common

import (
	"log/slog"
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/labstack/echo/v4"
)

// PageHandler manages rendering for generic pages like About.
type PageHandler struct {
	App *env.AppEnv
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(app *env.AppEnv) *PageHandler {
	return &PageHandler{App: app}
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	return h.renderWithBaseContext(c, "about.html", map[string]interface{}{
		"User": c.Get("User"),
	})
}

func (h *PageHandler) renderWithBaseContext(c echo.Context, tmpl string, data map[string]interface{}) error {
	ctx := c.Request().Context()
	categories, err := h.App.DB.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		slog.Error("Failed to fetch categories", "error", err)
	}

	counts, err := h.App.DB.GetCounts(ctx)
	if err != nil {
		slog.Error("Failed to fetch counts", "error", err)
		data["Counts"] = map[string]int{}
	} else {
		strCounts := make(map[string]int)
		for cat, count := range counts {
			strCounts[string(cat)] = count
		}
		data["Counts"] = strCounts
	}

	data["Categories"] = categories
	data["Config"] = h.App.Cfg
	data["Env"] = h.App.Cfg.Env
	data["DevMode"] = h.App.Cfg.Env == "development"
	data["HasGoogleAuth"] = h.App.Cfg.HasGoogleAuth
	return c.Render(http.StatusOK, tmpl, data)
}
