package admin

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

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
