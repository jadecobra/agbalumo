package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	Repo       domain.ListingRepository
	CSVService domain.CSVService
	Cfg        *config.Config
}

func NewAdminHandler(repo domain.ListingRepository, csvService domain.CSVService, cfg *config.Config) *AdminHandler {
	return &AdminHandler{Repo: repo, CSVService: csvService, Cfg: cfg}
}

// AdminMiddleware checks if the user is an admin.
func (h *AdminHandler) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("User")
		if user == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
		}

		u, ok := user.(domain.User)
		if !ok || u.Role != domain.UserRoleAdmin {
			// Redirect to claim page to enter access code
			return c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		}

		return next(c)
	}
}

// HandleLoginView renders the admin access code form.
func (h *AdminHandler) HandleLoginView(c echo.Context) error {
	// If already admin, redirect to dashboard
	user := c.Get("User")
	if user != nil {
		if u, ok := user.(domain.User); ok && u.Role == domain.UserRoleAdmin {
			return c.Redirect(http.StatusTemporaryRedirect, "/admin")
		}
	}
	// Pass empty map to avoid potential nil pointer issues in template engine
	return c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{})
}

// HandleLoginAction processes the access code and promotes the user.
func (h *AdminHandler) HandleLoginAction(c echo.Context) error {
	code := c.FormValue("code")

	if code != h.Cfg.AdminCode {
		return c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{
			"Error": "Invalid Access Code",
		})
	}

	// Promote User
	user := c.Get("User")
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}

	u, ok := user.(domain.User)
	if !ok {
		// Should not happen if OptionalAuth/RequireAuth are working, but handle safely
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}

	u.Role = domain.UserRoleAdmin
	// SaveUser now handles update via ID efficiently
	if err := h.Repo.SaveUser(c.Request().Context(), u); err != nil {
		return RespondError(c, err)
	}

	return c.Redirect(http.StatusFound, "/admin")
}

// HandleDashboard renders the admin dashboard.
func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	ctx := c.Request().Context()

	claimRequests, err := h.Repo.GetPendingClaimRequests(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	userCount, err := h.Repo.GetUserCount(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	feedbackCounts, err := h.Repo.GetFeedbackCounts(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	listingGrowth, err := h.Repo.GetListingGrowth(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	userGrowth, err := h.Repo.GetUserGrowth(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	feedbacks, err := h.Repo.GetAllFeedback(ctx)
	if err != nil {
		return RespondError(c, err)
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

	counts, _ := h.Repo.GetCounts(ctx)
	listingCount := 0
	for _, c := range counts {
		listingCount += c
	}

	categories, err := h.Repo.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		c.Logger().Errorf("failed to get categories: %v", err)
		categories = []domain.CategoryData{}
	}

	users, err := h.Repo.GetAllUsers(ctx, 10, 0)
	if err != nil {
		c.Logger().Errorf("failed to get users: %v", err)
		users = []domain.User{}
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
	existing, err := h.Repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false})
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

	err = h.Repo.SaveCategory(ctx, cat)
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
	p := GetPagination(c, 50)
	users, err := h.Repo.GetAllUsers(ctx, p.Limit, p.Offset)
	if err != nil {
		return RespondError(c, err)
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
	listings, _, err := h.Repo.FindAll(ctx, "", "", "created_at", "desc", true, 10000, 0)
	if err != nil {
		return RespondError(c, err)
	}

	reader, err := h.CSVService.GenerateCSV(ctx, listings)
	if err != nil {
		return RespondError(c, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="listings.csv"`)
	c.Response().WriteHeader(http.StatusOK)

	_, err = io.Copy(c.Response().Writer, reader)
	return err
}
