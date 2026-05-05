package common

import (
	"log/slog"
	"net/http"

	"github.com/jadecobra/agbalumo/internal/config"
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

// AboutViewModel represents the strongly-typed data passed to the about page template.
type AboutViewModel struct {
	User          interface{}
	Counts        map[string]int
	Config        *config.Config
	Env           string
	Categories    []domain.CategoryData
	DevMode       bool
	HasGoogleAuth bool
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	ctx := c.Request().Context()
	categories, err := h.App.DB.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		slog.Error("Failed to fetch categories", "error", err)
	}

	counts, err := h.App.DB.GetCounts(ctx)
	strCounts := make(map[string]int)
	if err != nil {
		slog.Error("Failed to fetch counts", "error", err)
	} else {
		for cat, count := range counts {
			strCounts[string(cat)] = count
		}
	}

	vm := AboutViewModel{
		User:          c.Get("User"),
		Categories:    categories,
		Counts:        strCounts,
		Config:        h.App.Cfg,
		Env:           h.App.Cfg.Env,
		DevMode:       h.App.Cfg.Env == "development",
		HasGoogleAuth: h.App.Cfg.HasGoogleAuth,
	}

	return c.Render(http.StatusOK, "about.html", vm)
}
