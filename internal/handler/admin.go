package handler

import (
	"fmt"
	"net/http"

	"github.com/jadecobra/agbalumo/internal/config"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	Repo       domain.ListingRepository
	CSVService *service.CSVService
	Cfg        *config.Config
}

func NewAdminHandler(repo domain.ListingRepository, csvService *service.CSVService, cfg *config.Config) *AdminHandler {
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

	pendingListings, err := h.Repo.GetPendingListings(ctx)
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

	// Get Flash Messages
	sess := customMiddleware.GetSession(c)
	var flashMsg interface{}
	if sess != nil {
		if flashes := sess.Flashes("message"); len(flashes) > 0 {
			flashMsg = flashes[0]
			sess.Save(c.Request(), c.Response())
		}
	}

	return c.Render(http.StatusOK, "admin_dashboard.html", map[string]interface{}{
		"PendingListings": pendingListings,
		"UserCount":       userCount,
		"FeedbackCounts":  feedbackCounts,
		"ListingGrowth":   listingGrowth,
		"UserGrowth":      userGrowth,
		"User":            c.Get("User"),
		"FlashMessage":    flashMsg,
	})
}

// HandleUsers renders the list of users for admins.
func (h *AdminHandler) HandleUsers(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.Repo.GetAllUsers(ctx)
	if err != nil {
		return RespondError(c, err)
	}

	return c.Render(http.StatusOK, "admin_users.html", map[string]interface{}{
		"Users": users,
		"User":  c.Get("User"),
	})
}

// HandleApprove approves a listing (clears Pending status).
func (h *AdminHandler) HandleApprove(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	listing.Status = domain.ListingStatusApproved
	listing.IsActive = true // Ensure it remains active

	if err := h.Repo.Save(ctx, listing); err != nil {
		return RespondError(c, err)
	}

	// HTMX Partial Update: Return empty 200 to remove the element from the list
	return c.NoContent(http.StatusOK)
}

// HandleReject rejects a listing (hides it and marks Rejected).
func (h *AdminHandler) HandleReject(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	listing.Status = domain.ListingStatusRejected
	listing.IsActive = false

	if err := h.Repo.Save(ctx, listing); err != nil {
		return RespondError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

// HandleBulkUpload processes a CSV file upload.
func (h *AdminHandler) HandleBulkUpload(c echo.Context) error {
	// 1. Get File
	file, err := c.FormFile("csv_file")
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Please select a valid CSV file"))
	}

	src, err := file.Open()
	if err != nil {
		return RespondError(c, err)
	}
	defer src.Close()

	// 2. Parse and Import
	result, err := h.CSVService.ParseAndImport(c.Request().Context(), src, h.Repo)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Failed to process CSV: "+err.Error()))
	}

	// 3. Render Result / Redirect
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		msg := fmt.Sprintf("Processed %d items. Success: %d, Failed: %d", result.TotalProcessed, result.SuccessCount, result.FailureCount)
		if len(result.Errors) > 0 {
			// Truncate errors if too long?
			if len(result.Errors) > 3 {
				msg += fmt.Sprintf(". Errors: %v ...", result.Errors[:3])
			} else {
				msg += fmt.Sprintf(". Errors: %v", result.Errors)
			}
		}
		sess.AddFlash(msg, "message")
		sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin")
}
