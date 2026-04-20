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

// HandleBulkUpload processes a CSV file upload.
func (h *AdminHandler) HandleBulkUpload(c echo.Context) error {
	// 1. Get File
	file, err := c.FormFile(domain.ParamCSVFile)
	if err != nil {
		return h.redirectWithFlash(c, "Please select a valid CSV file", domain.PathAdmin)
	}

	src, err := file.Open()
	if err != nil {
		return h.redirectWithFlash(c, "Failed to open file: "+err.Error(), domain.PathAdmin)
	}
	defer func() { _ = src.Close() }()

	// 2. Parse and Import
	result, err := h.App.CSVService.ParseAndImport(c.Request().Context(), src, h.App.DB)
	if err != nil {
		return h.redirectWithFlash(c, "Failed to process CSV: "+err.Error(), domain.PathAdmin)
	}

	// 3. Render Result / Redirect
	msg := fmt.Sprintf("Processed %d items. Success: %d, Failed: %d", result.TotalProcessed, result.SuccessCount, result.FailureCount)
	if len(result.Errors) > 0 {
		if len(result.Errors) > 3 {
			msg += fmt.Sprintf(". Errors: %v ...", result.Errors[:3])
		} else {
			msg += fmt.Sprintf(". Errors: %v", result.Errors)
		}
	}
	return h.redirectWithFlash(c, msg, domain.PathAdmin)
}
