package common

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
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
	Config *config.Config
	module.BaseViewData
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	vm := AboutViewModel{
		BaseViewData: h.PopulateBase(c),
		Config:       h.App.Cfg,
	}

	return h.RenderTyped(c, "about.html", vm)
}

// SandboxViewModel represents the strongly-typed data passed to the sandbox page template.
type SandboxViewModel struct {
	Config *config.Config
	module.BaseViewData
}

// HandleSandbox renders the component sandbox page for non-production environments.
func (h *PageHandler) HandleSandbox(c echo.Context) error {
	if h.App.Cfg.Env == domain.EnvProduction {
		return echo.NewHTTPError(http.StatusNotFound, "Sandbox not available in production")
	}

	vm := SandboxViewModel{
		BaseViewData: h.PopulateBase(c),
		Config:       h.App.Cfg,
	}

	return h.RenderTyped(c, "sandbox.html", vm)
}
