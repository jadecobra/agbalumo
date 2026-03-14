package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
)

// HandleBulkAction processes bulk approvals, rejections, and deletions.
func (h *AdminHandler) HandleBulkAction(c echo.Context) error {
	action := c.FormValue("action")
	selectedIDs := c.Request().PostForm["selectedListings"]
	ctx := c.Request().Context()

	if len(selectedIDs) == 0 {
		sess := customMiddleware.GetSession(c)
		if sess != nil {
			sess.AddFlash("No listings selected", "message")
			_ = sess.Save(c.Request(), c.Response())
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

	newCategory := c.FormValue("new_category")

	successCount := 0
	// Safe bounded admin action: N+1 here is acceptable because batch sizes are limited
	// by pagination (e.g. 50 items) and SQLite connection overhead is negligible.
	for _, id := range selectedIDs {
		listing, err := h.Repo.FindByID(ctx, id)
		if err != nil {
			continue // Skip if not found
		}

		switch action {
		case "approve":
			listing.Status = domain.ListingStatusApproved
			listing.IsActive = true
		case "reject":
			listing.Status = domain.ListingStatusRejected
			listing.IsActive = false
		case "change_category":
			if newCategory != "" {
				listing.Type = domain.Category(newCategory)
			} else {
				continue
			}
		default:
			continue // Unknown action
		}

		if err := h.Repo.Save(ctx, listing); err == nil {
			successCount++
		}
	}

	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.AddFlash(fmt.Sprintf("Successfully processed %d listings", successCount), "message")
		_ = sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin/listings")
}

// HandleBulkUpload processes a CSV file upload.
func (h *AdminHandler) HandleBulkUpload(c echo.Context) error {
	handleError := func(msg string) error {
		sess := customMiddleware.GetSession(c)
		if sess != nil {
			sess.AddFlash(msg, "message")
			_ = sess.Save(c.Request(), c.Response())
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
	defer func() { _ = src.Close() }()

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
		_ = sess.Save(c.Request(), c.Response())
	}

	return c.Redirect(http.StatusFound, "/admin")
}
