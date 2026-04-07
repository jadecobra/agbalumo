package admin

import (
	"github.com/jadecobra/agbalumo/internal/module/listing"

	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
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

	adminGroup := e.Group("/admin")
	adminGroup.Use(authMw.OptionalAuth)
	adminGroup.GET("/login", h.HandleLoginView)

	// Admin login POST with strict rate limiting
	adminLoginGroup := adminGroup.Group("/login")
	adminLoginGroup.Use(adminAuthLimiter.Middleware())
	adminLoginGroup.POST("", h.HandleLoginAction)
	adminGroup.Use(h.AdminMiddleware)
	adminGroup.GET("", h.HandleDashboard)
	adminGroup.GET("/users", h.HandleUsers)
	adminGroup.GET("/listings", h.HandleAllListings)
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
		u, ok := user.GetUser(c)
		if !ok || u == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
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
			return c.Redirect(http.StatusTemporaryRedirect, "/admin")
		}
	}
	// Pass empty map to avoid potential nil pointer issues in template engine
	return c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{})
}

// HandleLoginAction processes the access code and promotes the user.
func (h *AdminHandler) HandleLoginAction(c echo.Context) error {
	code := c.FormValue("code")

	if code != h.App.Cfg.AdminCode {
		return c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{
			"Error": "Invalid Access Code",
		})
	}

	// Promote User
	u, ok := user.GetUser(c)
	if !ok || u == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}

	u.Role = domain.UserRoleAdmin
	// SaveUser now handles update via ID efficiently
	if err := h.App.DB.SaveUser(c.Request().Context(), *u); err != nil {
		return ui.RespondError(c, err)
	}

	return c.Redirect(http.StatusFound, "/admin")
}

// HandleDashboard renders the admin dashboard.
func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	g, ctx := errgroup.WithContext(ctx)

	var (
		claimRequests  []domain.ClaimRequest
		userCount      int
		feedbackCounts map[domain.FeedbackType]int
		listingGrowth  []domain.DailyMetric
		userGrowth     []domain.DailyMetric
		feedbacks      []domain.Feedback
		counts         map[domain.Category]int
		categories     []domain.CategoryData
		users          []domain.User
	)

	g.Go(func() error {
		var err error
		claimRequests, err = h.App.DB.GetPendingClaimRequests(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		userCount, err = h.App.DB.GetUserCount(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		feedbackCounts, err = h.App.DB.GetFeedbackCounts(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		listingGrowth, err = h.App.DB.GetListingGrowth(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		userGrowth, err = h.App.DB.GetUserGrowth(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		feedbacks, err = h.App.DB.GetAllFeedback(ctx)
		return err
	})

	g.Go(func() error {
		// No error return expected for GetCounts as per original code
		counts, _ = h.App.DB.GetCounts(ctx)
		return nil
	})

	g.Go(func() error {
		var err error
		categories, err = h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{})
		if err != nil {
			c.Logger().Errorf("failed to get categories from service: %v", err)
			categories = []domain.CategoryData{}
		}
		return nil // Don't fail the whole dashboard if categories fail
	})

	g.Go(func() error {
		var err error
		users, err = h.App.DB.GetAllUsers(ctx, 10, 0)
		if err != nil {
			c.Logger().Errorf("failed to get users: %v", err)
			users = []domain.User{}
		}
		return nil // Don't fail the whole dashboard if users fail
	})

	if err := g.Wait(); err != nil {
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

	listingCount := 0
	for _, count := range counts {
		listingCount += count
	}

	return c.Render(http.StatusOK, "admin_dashboard.html", map[string]interface{}{
		"ClaimRequests":  claimRequests,
		"UserCount":      userCount,
		"FeedbackCounts": feedbackCounts,
		"ListingGrowth":  listingGrowth,
		"UserGrowth":     userGrowth,
		"Feedbacks":      feedbacks,
		"User":           c.Get("User"),
		"FlashMessage":   flashMsg,
		"ListingCount":   listingCount,
		"Categories":     categories,
		"Users":          users,
	})
}

// HandleAddCategory processes the form submission to add a new category
func (h *AdminHandler) HandleAddCategory(c echo.Context) error {
	ctx := c.Request().Context()

	name := strings.TrimSpace(c.FormValue("name"))
	if name == "" {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Case-insensitive check for existing category
	existing, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
	if err == nil {
		for _, cat := range existing {
			if strings.EqualFold(cat.Name, name) {
				sess := customMiddleware.GetSession(c)
				if sess != nil {
					sess.AddFlash(fmt.Sprintf("Category '%s' already exists!", cat.Name), "message")
					_ = sess.Save(c.Request(), c.Response())
				}
				return c.Redirect(http.StatusFound, "/admin")
			}
		}
	}

	claimableStr := c.FormValue("claimable")
	claimable := claimableStr == "true"

	now := time.Now()
	cat := domain.CategoryData{
		ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Name:      name,
		Claimable: claimable,
		IsSystem:  false,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = h.App.DB.SaveCategory(ctx, cat)
	if err != nil {
		c.Logger().Errorf("failed to save custom category: %v", err)
	}

	// Add success flash message
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash("Category added successfully!", "message")
		_ = sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin")
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
