package admin

import (
	"context"
	"fmt"
	"net/http"
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

	category := c.QueryParam("category")
	sortField := c.QueryParam("sort")
	sortOrder := strings.ToUpper(c.QueryParam("order"))
	queryText := strings.TrimSpace(c.QueryParam("q"))

	// Fetch all listings with the given category filter, including inactive ones.
	listings, totalCountRows, err := h.App.DB.FindAll(ctx, category, queryText, sortField, sortOrder, true, pagination.Limit, pagination.Offset)
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
		"User":       c.Get("User"),
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

	featured := c.FormValue("featured") == "true"
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
	featured, err := h.App.DB.GetFeaturedListings(ctx, string(listing.Type))
	if err != nil {
		return err
	}
	if len(featured) >= 3 {
		return fmt.Errorf("Maximum of 3 featured listings per category reached")
	}
	return nil
}
