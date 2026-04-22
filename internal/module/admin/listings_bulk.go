package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// HandleBulkAction processes bulk approvals, rejections, and deletions.
func (h *AdminHandler) HandleBulkAction(c echo.Context) error {
	action := c.FormValue(domain.FieldAction)
	selectedIDs := c.Request().PostForm[domain.ParamListingIDs]
	ctx := c.Request().Context()

	if len(selectedIDs) == 0 {
		return h.redirectWithFlash(c, "No listings selected", domain.PathAdminListings)
	}

	if action == "delete" {
		return h.redirectToBulkDeleteConfirm(c, selectedIDs)
	}

	newCategory := c.FormValue(domain.FieldNewCategory)
	successCount := h.processBulkListings(ctx, selectedIDs, action, newCategory)

	return h.redirectWithFlash(c, fmt.Sprintf("Successfully processed %d listings", successCount), domain.PathAdminListings)
}

func (h *AdminHandler) redirectToBulkDeleteConfirm(c echo.Context, ids []string) error {
	query := url.Values{}
	for _, id := range ids {
		query.Add(domain.ParamID, id)
	}
	return c.Redirect(http.StatusFound, domain.PathAdminListings+"/delete-confirm?"+query.Encode())
}

func (h *AdminHandler) processBulkListings(ctx context.Context, ids []string, action, newCategory string) int {
	successCount := 0
	for _, id := range ids {
		if err := h.applyActionToListing(ctx, id, action, newCategory); err == nil {
			successCount++
		}
	}
	return successCount
}

func (h *AdminHandler) applyActionToListing(ctx context.Context, id, action, newCategory string) error {
	listing, err := h.App.DB.FindByID(ctx, id)
	if err != nil {
		return err
	}

	switch action {
	case "approve":
		listing.Status = domain.ListingStatusApproved
		listing.IsActive = true
	case "reject":
		listing.Status = domain.ListingStatusRejected
		listing.IsActive = false
		if newCategory != "" {
			listing.Type = domain.Category(newCategory)
		}
	case "change_category":
		if newCategory == "" {
			return fmt.Errorf("category required")
		}
		listing.Type = domain.Category(newCategory)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return h.App.DB.Save(ctx, listing)
}

// HandleAdminDeleteView renders the double-confirmation page for deleting listings.
func (h *AdminHandler) HandleAdminDeleteView(c echo.Context) error {
	ids := c.QueryParams()[domain.ParamID]
	if len(ids) == 0 {
		if id := c.QueryParam(domain.ParamID); id != "" {
			ids = []string{id}
		}
	}

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, domain.PathAdminListings)
	}

	ctx := c.Request().Context()
	// Safe bounded admin action: N+1 here is acceptable because batch sizes are limited
	// by pagination (e.g. 50 items) and SQLite connection overhead is negligible.
	for _, id := range ids {
		if _, err := h.App.DB.FindByID(ctx, id); err != nil {
			return c.String(http.StatusNotFound, domain.ErrListingNotFound.Error())
		}
	}

	return c.Render(http.StatusOK, "admin_delete_confirm.html", map[string]interface{}{
		"IDs":  ids,
		"User": c.Get(domain.CtxKeyUser),
	})
}

// HandleAdminDeleteAction processes explicit admin deletions after password confirmation.
func (h *AdminHandler) HandleAdminDeleteAction(c echo.Context) error {
	adminCode := c.FormValue(domain.FieldAdminCode)

	_ = c.Request().ParseForm()
	ids := c.Request().PostForm[domain.ParamID]

	if len(ids) == 0 {
		return c.Redirect(http.StatusFound, domain.PathAdminListings)
	}

	if adminCode != h.App.Cfg.AdminCode {
		return c.Render(http.StatusOK, "admin_delete_confirm.html", map[string]interface{}{
			"IDs":   ids,
			"Error": "Invalid Admin Code. Deletion aborted.",
			"User":  c.Get(domain.CtxKeyUser),
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

	return h.redirectWithFlash(c, fmt.Sprintf("Successfully deleted %d listings", successCount), domain.PathAdminListings)
}
