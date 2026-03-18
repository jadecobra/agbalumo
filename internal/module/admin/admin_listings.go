package admin

import (
	"github.com/jadecobra/agbalumo/internal/module/listing"

	"github.com/jadecobra/agbalumo/internal/handler"

	"net/http"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// HandleAllListings renders the list of all listings for admins, with category filtering.
func (h *AdminHandler) HandleAllListings(c echo.Context) error {
	ctx := c.Request().Context()

	pagination := listing.GetPagination(c, 50)

	category := c.QueryParam("category")
	sortField := c.QueryParam("sort")
	sortOrder := strings.ToUpper(c.QueryParam("order"))
	queryText := strings.TrimSpace(c.QueryParam("q"))

	// Fetch all listings with the given category filter, including inactive ones.
	listings, totalCountRows, err := h.ListingStore.FindAll(ctx, category, queryText, sortField, sortOrder, true, pagination.Limit, pagination.Offset)
	if err != nil {
		return handler.RespondError(c, err)
	}

	hasNextPage := pagination.Offset+len(listings) < totalCountRows

	counts, err := h.ListingStore.GetCounts(ctx)
	if err != nil {
		c.Logger().Errorf("failed to get listing counts: %v", err)
		counts = make(map[domain.Category]int)
	}

	categories, err := h.CategoryStore.GetCategories(ctx, domain.CategoryFilter{})
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

// HandleToggleFeatured toggles the featured status of a listing.
func (h *AdminHandler) HandleToggleFeatured(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Listing ID is required"})
	}

	featured := c.FormValue("featured") == "true"
	ctx := c.Request().Context()

	if err := h.ListingStore.SetFeatured(ctx, id, featured); err != nil {
		return handler.RespondError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       id,
		"featured": featured,
	})
}
