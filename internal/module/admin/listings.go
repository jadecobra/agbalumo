package admin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

const tmplListingTableRow = "admin_listing_table_row"

// HandleAllListings renders the list of all listings for admins, with category filtering.
func (h *AdminHandler) HandleAllListings(c echo.Context) error {
	ctx := c.Request().Context()

	pagination := listing.GetPagination(c, 50)

	category := c.QueryParam(domain.ParamCategory)
	sortField := c.QueryParam(domain.ParamSort)
	sortOrder := strings.ToUpper(c.QueryParam(domain.ParamOrder))
	queryText := strings.TrimSpace(c.QueryParam(domain.ParamQuery))

	// Fetch all listings with the given category filter, including inactive ones.
	listings, totalCountRows, err := h.App.DB.FindAll(ctx, category, queryText, "", sortField, sortOrder, true, pagination.Limit, pagination.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}

	hasNextPage := pagination.Offset+len(listings) < totalCountRows

	counts, err := h.App.DB.GetCounts(ctx)
	if err != nil {
		c.Logger().Errorf("failed to get listing counts: %v", err)
		counts = make(map[domain.Category]int)
	}

	categories, err := h.App.DB.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		c.Logger().Errorf("failed to get categories: %v", err)
		categories = []domain.CategoryData{}
	}

	strCounts, _ := listing.ConvertCounts(counts)

	return c.Render(http.StatusOK, "admin_listings.html", map[string]interface{}{
		"Listings":   listings,
		"Pagination": listing.Pagination{Page: pagination.Page, TotalPages: (totalCountRows + pagination.Limit - 1) / pagination.Limit, HasNextPage: hasNextPage, TotalCount: totalCountRows},
		"Category":   category,
		"SortField":  sortField,
		"SortOrder":  sortOrder,
		"QueryText":  queryText,
		"Counts":     strCounts,
		"Categories": categories,
		"TotalCount": totalCountRows, // Use totalCountRows from FindAll for consistent count
		"User":       c.Get(domain.CtxKeyUser),
	})
}

// HandleListingRow renders a single table row for a listing.
func (h *AdminHandler) HandleListingRow(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	listing, err := h.App.DB.FindByID(ctx, id)
	if err != nil {
		return ui.RespondError(c, err)
	}
	return h.renderListingRow(c, listing)
}

// HandleToggleFeatured toggles the featured status of a listing.
func (h *AdminHandler) HandleToggleFeatured(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return ui.RespondJSONError(c, http.StatusBadRequest, "Listing ID is required")
	}

	featured := c.FormValue(domain.FieldFeatured) == "true"

	ctx := c.Request().Context()

	listing, err := h.App.DB.FindByID(ctx, id)
	if err != nil {
		return ui.RespondError(c, err)
	}

	if featured {
		if err := h.validateFeaturedLimit(ctx, listing); err != nil {
			return ui.RespondJSONError(c, http.StatusBadRequest, err.Error())
		}
	}

	if err := h.App.DB.SetFeatured(ctx, id, featured); err != nil {
		return ui.RespondError(c, err)
	}

	updatedListing, _ := h.App.DB.FindByID(ctx, id)
	return h.renderListingRow(c, updatedListing)
}

func (h *AdminHandler) renderListingRow(c echo.Context, listing domain.Listing) error {
	return c.Render(http.StatusOK, tmplListingTableRow, listing)
}

func (h *AdminHandler) validateFeaturedLimit(ctx context.Context, listing domain.Listing) error {
	if listing.Featured {
		return nil
	}
	featured, err := h.App.DB.GetFeaturedListings(ctx, string(listing.Type), "")
	if err != nil {
		return err
	}
	if len(featured) >= 3 {
		return fmt.Errorf("Maximum of 3 featured listings per category reached")
	}
	return nil
}

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

// HandleExportListings generates and serves a CSV of all listings.
func (h *AdminHandler) HandleExportListings(c echo.Context) error {
	ctx := c.Request().Context()

	// Fetch all listings. Using a large limit for export.
	// In a very large system, we might want to stream this from the DB directly.
	listings, _, err := h.App.DB.FindAll(ctx, "", "", "", domain.FieldCreatedAt, "desc", true, 10000, 0)

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
