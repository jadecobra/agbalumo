package admin

import (
	"github.com/jadecobra/agbalumo/internal/module/listing"

	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type AdminHandler struct {
	App *env.AppEnv
}

func NewAdminHandler(app *env.AppEnv) *AdminHandler {
	return &AdminHandler{
		App: app,
	}
}

// RegisterRoutes registers all admin-related routes.
func (h *AdminHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	// Strict rate limiter for sensitive admin login endpoint (5 req/min, burst 10)
	adminAuthLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  rate.Limit(5),
		Burst: 10,
	})

	adminGroup := e.Group(domain.PathAdmin)
	adminGroup.Use(authMw.OptionalAuth)
	adminGroup.GET("/login", h.HandleLoginView)

	// Admin login POST with strict rate limiting
	adminLoginGroup := adminGroup.Group("/login")
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
	code := c.FormValue("code")

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

// HandleDashboard renders the admin dashboard.
func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.loadDashboardData(ctx, c)
	if err != nil {
		return ui.RespondError(c, err)
	}

	// Get Flash Messages
	sess := customMiddleware.GetSession(c)
	var flashMsg interface{}
	if sess != nil {
		if flashes := sess.Flashes("message"); len(flashes) > 0 {
			flashMsg = flashes[0]
			_ = sess.Save(c.Request(), c.Response())
		}
	}

	return c.Render(http.StatusOK, "admin_dashboard.html", map[string]interface{}{
		"ClaimRequests":   data.ClaimRequests,
		"UserCount":       data.UserCount,
		"FeedbackCounts":  data.FeedbackCounts,
		"ListingGrowth":   data.ListingGrowth,
		"UserGrowth":      data.UserGrowth,
		"Feedbacks":       data.Feedbacks,
		"User":            c.Get("User"),
		"FlashMessage":    flashMsg,
		"ListingCount":    data.ListingCount,
		"Categories":      data.Categories,
		"Users":           data.Users,
		"AdaDiscoveryAvg": data.AdaDiscoveryAvg,
	})
}

type dashboardData struct {
	ClaimRequests   []domain.ClaimRequest
	FeedbackCounts  map[domain.FeedbackType]int
	ListingGrowth   []domain.DailyMetric
	UserGrowth      []domain.DailyMetric
	Feedbacks       []domain.Feedback
	Categories      []domain.CategoryData
	Users           []domain.User
	UserCount       int
	ListingCount    int
	AdaDiscoveryAvg float64
}

func (h *AdminHandler) loadDashboardData(ctx context.Context, c echo.Context) (dashboardData, error) {
	var data dashboardData
	var err error

	data.ClaimRequests, err = h.App.DB.GetPendingClaimRequests(ctx)
	if err != nil {
		return data, err
	}

	data.UserCount, err = h.App.DB.GetUserCount(ctx)
	if err != nil {
		return data, err
	}

	data.FeedbackCounts, err = h.App.DB.GetFeedbackCounts(ctx)
	if err != nil {
		return data, err
	}

	data.ListingGrowth, err = h.App.DB.GetListingGrowth(ctx)
	if err != nil {
		return data, err
	}

	data.UserGrowth, err = h.App.DB.GetUserGrowth(ctx)
	if err != nil {
		return data, err
	}

	data.Feedbacks, err = h.App.DB.GetAllFeedback(ctx)
	if err != nil {
		return data, err
	}

	counts, _ := h.App.DB.GetCounts(ctx)
	for _, count := range counts {
		data.ListingCount += count
	}

	data.Categories, err = h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		c.Logger().Errorf("failed to get categories from service: %v", err)
		data.Categories = []domain.CategoryData{}
	}

	data.Users, err = h.App.DB.GetAllUsers(ctx, 10, 0)
	if err != nil {
		c.Logger().Errorf("failed to get users: %v", err)
		data.Users = []domain.User{}
	}

	// Fetch Ada Metrics (Last 24h)
	since := time.Now().Add(-24 * time.Hour)
	data.AdaDiscoveryAvg, err = h.App.DB.GetAverageValue(ctx, "discovery_success", since)
	if err != nil {
		c.Logger().Errorf("failed to get Ada metrics: %v", err)
	}

	return data, nil
}

func (h *AdminHandler) renderLoginView(c echo.Context, errMsg string) error {
	data := map[string]interface{}{}
	if errMsg != "" {
		data["Error"] = errMsg
	}
	return c.Render(http.StatusOK, "admin_login.html", data)
}

// HandleAddCategory processes the form submission to add a new category
func (h *AdminHandler) HandleAddCategory(c echo.Context) error {
	ctx := c.Request().Context()
	name := strings.TrimSpace(c.FormValue("name"))
	if name == "" {
		return c.Redirect(http.StatusFound, domain.PathAdmin)
	}

	if existing, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false}); err == nil {
		if hasDuplicateCategory(existing, name) {
			return flashAndRedirect(c, fmt.Sprintf("Category '%s' already exists!", name), domain.PathAdmin)
		}
	}

	claimable := c.FormValue("claimable") == "true"
	now := time.Now()
	cat := domain.CategoryData{
		ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Name:      name,
		Claimable: claimable,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.App.DB.SaveCategory(ctx, cat); err != nil {
		c.Logger().Errorf("failed to save custom category: %v", err)
	}

	return flashAndRedirect(c, "Category added successfully!", domain.PathAdmin)
}

// HandleUsers renders the list of users for admins.
func (h *AdminHandler) HandleUsers(c echo.Context) error {
	ctx := c.Request().Context()
	p := listing.GetPagination(c, 50)
	users, err := h.App.DB.GetAllUsers(ctx, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}
	p.HasNextPage = len(users) == p.Limit

	return c.Render(http.StatusOK, "admin_users.html", map[string]interface{}{
		"Users":      users,
		"User":       c.Get("User"),
		"Pagination": p,
	})
}

// HandleExportListings generates and serves a CSV of all listings.
func (h *AdminHandler) HandleExportListings(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch all listings. Using a large limit for export.
	// In a very large system, we might want to stream this from the DB directly.
	listings, _, err := h.App.DB.FindAll(ctx, "", "", "created_at", "desc", true, 10000, 0)
	if err != nil {
		return ui.RespondError(c, err)
	}

	reader, err := h.App.CSVService.GenerateCSV(ctx, listings)
	if err != nil {
		return ui.RespondError(c, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="listings.csv"`)
	c.Response().WriteHeader(http.StatusOK)

	_, err = io.Copy(c.Response().Writer, reader)
	return err
}

func hasDuplicateCategory(existing []domain.CategoryData, name string) bool {
	for _, cat := range existing {
		if strings.EqualFold(cat.Name, name) {
			return true
		}
	}
	return false
}

func flashAndRedirect(c echo.Context, msg, url string) error {
	if sess := customMiddleware.GetSession(c); sess != nil {
		sess.AddFlash(msg, "message")
		_ = sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusFound, url)
}
