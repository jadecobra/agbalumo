package handler

import (
	"net/http"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

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
