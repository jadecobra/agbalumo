package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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
			sess.Save(c.Request(), c.Response())
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

	name := c.FormValue("name")
	if name == "" {
		return c.Redirect(http.StatusFound, "/admin")
	}

	claimableStr := c.FormValue("claimable")
	claimable := claimableStr == "true"

	now := time.Now()
	cat := domain.CategoryData{
		ID:        strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-")),
		Name:      name,
		Claimable: claimable,
		IsSystem:  false,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := h.Repo.SaveCategory(ctx, cat)
	if err != nil {
		c.Logger().Errorf("failed to save custom category: %v", err)
	}

	// Add success flash message (optional but good practice)
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash("Category added successfully!", "message")
		sess.Save(c.Request(), c.Response())
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

// HandleAllListings renders the list of all listings for admins, with category filtering.
func (h *AdminHandler) HandleAllListings(c echo.Context) error {
	ctx := c.Request().Context()

	pagination := GetPagination(c, 50)

	category := c.QueryParam("category")
	sortField := c.QueryParam("sort")
	sortOrder := strings.ToUpper(c.QueryParam("order"))

	// Fetch all listings with the given category filter, including inactive ones.
	listings, err := h.Repo.FindAll(ctx, category, "", sortField, sortOrder, true, pagination.Limit, pagination.Offset)
	if err != nil {
		return RespondError(c, err)
	}

	hasNextPage := len(listings) == pagination.Limit

	counts, err := h.Repo.GetCounts(ctx)
	if err != nil {
		c.Logger().Errorf("failed to get listing counts: %v", err)
		counts = make(map[domain.Category]int)
	}

	categories, err := h.Repo.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		c.Logger().Errorf("failed to get categories: %v", err)
		categories = []domain.CategoryData{}
	}

	strCounts, totalCount := ConvertCounts(counts)

	return c.Render(http.StatusOK, "admin_listings.html", map[string]interface{}{
		"Listings":    listings,
		"Page":        pagination.Page,
		"HasNextPage": hasNextPage,
		"Category":    category,
		"SortField":   sortField,
		"SortOrder":   sortOrder,
		"Counts":      strCounts,
		"Categories":  categories,
		"TotalCount":  totalCount,
		"User":        c.Get("User"),
	})
}

// HandleApproveClaim approves a user's claim request and transfers listing ownership.
func (h *AdminHandler) HandleApproveClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.Repo.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusApproved); err != nil {
		return c.String(http.StatusNotFound, "Claim request not found")
	}

	return c.NoContent(http.StatusOK)
}

// HandleRejectClaim rejects a user's claim request.
func (h *AdminHandler) HandleRejectClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.Repo.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusRejected); err != nil {
		return c.String(http.StatusNotFound, "Claim request not found")
	}

	return c.NoContent(http.StatusOK)
}

// HandleBulkAction processes bulk approvals, rejections, and deletions.
func (h *AdminHandler) HandleBulkAction(c echo.Context) error {
	action := c.FormValue("action")
	selectedIDs := c.Request().PostForm["selectedListings"]
	ctx := c.Request().Context()

	if len(selectedIDs) == 0 {
		sess := customMiddleware.GetSession(c)
		if sess != nil {
			sess.AddFlash("No listings selected", "message")
			sess.Save(c.Request(), c.Response())
		}
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	if action == "delete" {
		// Pass IDs as query parameters to the confirmation page
		query := url.Values{}
		for _, id := range selectedIDs {
			query.Add("id", id)
		}
		return c.Redirect(http.StatusFound, "/admin/listings/delete-confirm?"+query.Encode())
	}

	successCount := 0
	for _, id := range selectedIDs {
		listing, err := h.Repo.FindByID(ctx, id)
		if err != nil {
			continue // Skip if not found
		}

		if action == "approve" {
			listing.Status = domain.ListingStatusApproved
			listing.IsActive = true
		} else if action == "reject" {
			listing.Status = domain.ListingStatusRejected
			listing.IsActive = false
		} else {
			continue // Unknown action
		}

		if err := h.Repo.Save(ctx, listing); err == nil {
			successCount++
		}
	}

	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash(fmt.Sprintf("Successfully processed %d listings", successCount), "message")
		sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin/listings")
}

// HandleAdminDeleteView renders the double-confirmation page for deleting listings.
func (h *AdminHandler) HandleAdminDeleteView(c echo.Context) error {
	// Parse IDs from query parameters (can be multiple)
	c.Request().ParseForm()
	ids := c.Request().Form["id"]

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	ctx := c.Request().Context()
	for _, id := range ids {
		if _, err := h.Repo.FindByID(ctx, id); err != nil {
			return c.String(http.StatusNotFound, "Listing not found")
		}
	}

	return c.Render(http.StatusOK, "admin_delete_confirm.html", map[string]interface{}{
		"IDs":  ids,
		"User": c.Get("User"),
	})
}

// HandleAdminDeleteAction processes explicit admin deletions after password confirmation.
func (h *AdminHandler) HandleAdminDeleteAction(c echo.Context) error {
	adminCode := c.FormValue("admin_code")

	// Parse IDs (can be multiple)
	c.Request().ParseForm()
	ids := c.Request().PostForm["id"]

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	// 1. Password (Admin Code) Verification
	if adminCode != h.Cfg.AdminCode {
		return c.Render(http.StatusOK, "admin_delete_confirm.html", map[string]interface{}{
			"IDs":   ids,
			"Error": "Invalid Admin Code. Deletion aborted.",
			"User":  c.Get("User"),
		})
	}

	// 2. Perform Deletions
	ctx := c.Request().Context()
	successCount := 0
	for _, id := range ids {
		if err := h.Repo.Delete(ctx, id); err == nil {
			successCount++
		} else {
			c.Logger().Errorf("Failed to delete listing %s: %v", id, err)
		}
	}

	// 3. Feedback
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash(fmt.Sprintf("Successfully deleted %d listings", successCount), "message")
		sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin/listings")
}

// HandleBulkUpload processes a CSV file upload.
func (h *AdminHandler) HandleBulkUpload(c echo.Context) error {
	handleError := func(msg string) error {
		sess := customMiddleware.GetSession(c)
		if sess != nil {
			sess.AddFlash(msg, "message")
			sess.Save(c.Request(), c.Response())
		}
		return c.Redirect(http.StatusFound, "/admin")
	}

	// 1. Get File
	file, err := c.FormFile("csv_file")
	if err != nil {
		return handleError("Please select a valid CSV file")
	}

	src, err := file.Open()
	if err != nil {
		return handleError("Failed to open file: " + err.Error())
	}
	defer src.Close()

	// 2. Parse and Import
	result, err := h.CSVService.ParseAndImport(c.Request().Context(), src, h.Repo)
	if err != nil {
		return handleError("Failed to process CSV: " + err.Error())
	}

	// 3. Render Result / Redirect
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		msg := fmt.Sprintf("Processed %d items. Success: %d, Failed: %d", result.TotalProcessed, result.SuccessCount, result.FailureCount)
		if len(result.Errors) > 0 {
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

// HandleToggleFeatured toggles the featured status of a listing.
func (h *AdminHandler) HandleToggleFeatured(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Listing ID is required"})
	}

	featured := c.FormValue("featured") == "true"
	ctx := c.Request().Context()

	if err := h.Repo.SetFeatured(ctx, id, featured); err != nil {
		return RespondError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       id,
		"featured": featured,
	})
}
