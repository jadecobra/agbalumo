package common

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// PageHandler manages rendering for generic pages like About.
type PageHandler struct {
	CategoryStore domain.CategoryStore
	Cfg           *config.Config
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(categoryStore domain.CategoryStore, cfg *config.Config) *PageHandler {
	return &PageHandler{
		CategoryStore: categoryStore,
		Cfg:           cfg,
	}
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	return h.renderWithBaseContext(c, "about.html", map[string]interface{}{
		"User": c.Get("User"),
	})
}

func (h *PageHandler) renderWithBaseContext(c echo.Context, tmpl string, data map[string]interface{}) error {
	ctx := c.Request().Context()
	categories, err := h.CategoryStore.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		c.Logger().Errorf("Failed to retrieve categories: %v", err)
		categories = []domain.CategoryData{}
	}

	data["Categories"] = categories
	data["Env"] = h.Cfg.Env
	data["HasGoogleAuth"] = h.Cfg.HasGoogleAuth
	return c.Render(http.StatusOK, tmpl, data)
}
