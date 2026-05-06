package common

import (
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/labstack/echo/v4"
)

// PageHandler manages rendering for generic pages like About.
type PageHandler struct {
	module.BaseHandler
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(app *env.AppEnv) *PageHandler {
	return &PageHandler{
		BaseHandler: module.BaseHandler{App: app},
	}
}

// AboutViewModel represents the strongly-typed data passed to the about page template.
type AboutViewModel struct {
	module.BaseViewData
	Config *config.Config
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	vm := AboutViewModel{
		BaseViewData: h.PopulateBase(c),
		Config:       h.App.Cfg,
	}

	return h.RenderTyped(c, "about.html", vm)
}
