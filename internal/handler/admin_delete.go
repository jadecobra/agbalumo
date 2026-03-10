package handler

import (
	"fmt"
	"net/http"

	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

// HandleAdminDeleteView renders the double-confirmation page for deleting listings.
func (h *AdminHandler) HandleAdminDeleteView(c echo.Context) error {
	// Parse IDs from query parameters (can be multiple)
	_ = c.Request().ParseForm()
	ids := c.Request().Form["id"]

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	ctx := c.Request().Context()
	// Safe bounded admin action: N+1 here is acceptable because batch sizes are limited
	// by pagination (e.g. 50 items) and SQLite connection overhead is negligible.
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
	_ = c.Request().ParseForm()
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
	// Safe bounded admin action: N+1 here is acceptable because batch sizes are limited
	// by pagination (e.g. 50 items) and SQLite connection overhead is negligible.
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
		_ = sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin/listings")
}
