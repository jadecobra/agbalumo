package admin

import (
	"fmt"
	"net/http"

	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

// HandleAdminDeleteView renders the double-confirmation page for deleting listings.
func (h *AdminHandler) HandleAdminDeleteView(c echo.Context) error {
	ids := c.QueryParams()["id"]
	if len(ids) == 0 {
		if id := c.QueryParam("id"); id != "" {
			ids = []string{id}
		}
	}

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	ctx := c.Request().Context()
	// Safe bounded admin action: N+1 here is acceptable because batch sizes are limited
	// by pagination (e.g. 50 items) and SQLite connection overhead is negligible.
	for _, id := range ids {
		if _, err := h.App.DB.FindByID(ctx, id); err != nil {
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

	_ = c.Request().ParseForm()
	ids := c.Request().PostForm["id"]

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, "/admin/listings")
	}

	if adminCode != h.App.Cfg.AdminCode {
		return c.Render(http.StatusOK, "admin_delete_confirm.html", map[string]interface{}{
			"IDs":   ids,
			"Error": "Invalid Admin Code. Deletion aborted.",
			"User":  c.Get("User"),
		})
	}

	ctx := c.Request().Context()
	successCount := 0
	for _, id := range ids {
		if err := h.App.DB.Delete(ctx, id); err == nil {
			successCount++
		} else {
			c.Logger().Errorf("Failed to delete listing %s: %v", id, err)
		}
	}

	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash(fmt.Sprintf("Successfully deleted %d listings", successCount), "message")
		_ = sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin/listings")
}
