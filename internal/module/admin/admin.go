package admin

import (
	"fmt"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type AdminHandler struct {
	module.BaseHandler
}

func NewAdminHandler(app *env.AppEnv) *AdminHandler {
	return &AdminHandler{
		BaseHandler: module.BaseHandler{App: app},
	}
}

// RegisterRoutes registers all admin-related routes.
func (h *AdminHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	if h == nil {
		fmt.Println("⚠️  CRITICAL: Attempted to register routes on nil AdminHandler")
		return
	}
	// Strict rate limiter for sensitive admin login endpoint (5 req/min, burst 10)
	adminAuthLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  rate.Limit(5),
		Burst: 10,
	})

	adminGroup := e.Group(domain.PathAdmin)
	adminGroup.Use(authMw.OptionalAuth)
	adminGroup.GET(domain.PathLogin, h.HandleLoginView)

	// Admin login POST with strict rate limiting
	adminLoginGroup := adminGroup.Group(domain.PathLogin)
	adminLoginGroup.Use(adminAuthLimiter.Middleware())
	adminLoginGroup.POST("", h.HandleLoginAction)
	adminGroup.Use(h.AdminMiddleware)
	adminGroup.GET("", h.HandleDashboard)
	adminGroup.GET("/users", h.HandleUsers)
	adminGroup.GET(domain.PathListings, h.HandleAllListings)
	adminGroup.POST("/claims/:id/approve", h.HandleApproveClaim)
	adminGroup.POST("/claims/:id/reject", h.HandleRejectClaim)
	adminGroup.POST("/listings/bulk", h.HandleBulkAction)
	adminGroup.GET("/listings/:id/row", h.HandleListingRow)
	adminGroup.GET("/listings/delete-confirm", h.HandleAdminDeleteView)
	adminGroup.POST("/listings/delete", h.HandleAdminDeleteAction)
	adminGroup.POST("/listings/:id/featured", h.HandleToggleFeatured)
	adminGroup.POST("/upload", h.HandleBulkUpload)
	adminGroup.GET("/listings/export", h.HandleExportListings)
	adminGroup.POST("/categories", h.HandleAddCategory)

	// Modal Fragments
	adminGroup.GET("/modal/charts", h.HandleModalCharts)
	adminGroup.GET("/modal/users", h.HandleModalUsers)
	adminGroup.GET("/modal/bulk", h.HandleModalBulk)
	adminGroup.GET("/modal/category", h.HandleModalCategory)
	adminGroup.GET("/modal/moderation", h.HandleModalModeration)
}

// AdminMiddleware checks if the user is an admin.
func (h *AdminHandler) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := user.RequireUser(c)
		if err != nil || u == nil {
			return err
		}

		if u.Role != domain.UserRoleAdmin {
			// Redirect to claim page to enter access code
			return c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		}

		return next(c)
	}
}

// HandleLoginView renders the admin access code form.
func (h *AdminHandler) HandleLoginView(c echo.Context) error {
	// If already admin, redirect to dashboard
	if u, ok := user.GetUser(c); ok && u != nil {
		if u.Role == domain.UserRoleAdmin {
			return c.Redirect(http.StatusTemporaryRedirect, domain.PathAdmin)
		}
	}
	// Pass empty string for no error message
	return h.renderLoginView(c, "")
}

// HandleLoginAction processes the access code and promotes the user.
func (h *AdminHandler) HandleLoginAction(c echo.Context) error {
	code := c.FormValue(domain.FieldCode)

	if code != h.App.Cfg.AdminCode {
		return h.renderLoginView(c, "Invalid Access Code")
	}

	// Promote User
	u, err := user.RequireUser(c)
	if err != nil || u == nil {
		return err
	}

	u.Role = domain.UserRoleAdmin
	// SaveUser now handles update via ID efficiently
	if err := h.App.DB.SaveUser(c.Request().Context(), *u); err != nil {
		return ui.RespondError(c, err)
	}

	return c.Redirect(http.StatusFound, domain.PathAdmin)
}

func (h *AdminHandler) renderLoginView(c echo.Context, errMsg string) error {
	data := map[string]interface{}{}
	if errMsg != "" {
		data["Error"] = errMsg
	}
	return c.Render(http.StatusOK, "admin_login.html", data)
}

func (h *AdminHandler) redirectWithFlash(c echo.Context, msg, targetURL string) error {
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash(msg, domain.FlashMessageKey)
		_ = sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusFound, targetURL)
}
