package admin

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// HandleModalCharts renders the growth charts modal fragment.
func (h *AdminHandler) HandleModalCharts(c echo.Context) error {
	ctx := c.Request().Context()
	listingGrowth, err := h.App.DB.GetListingGrowth(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}
	userGrowth, err := h.App.DB.GetUserGrowth(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}

	return c.Render(http.StatusOK, "admin_modal_charts.html", map[string]interface{}{
		"ListingGrowth": listingGrowth,
		"UserGrowth":    userGrowth,
	})
}

// HandleModalUsers renders the user management modal fragment.
func (h *AdminHandler) HandleModalUsers(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.App.DB.GetAllUsers(ctx, 10, 0)
	if err != nil {
		return ui.RespondError(c, err)
	}
	userCount, err := h.App.DB.GetUserCount(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}

	return c.Render(http.StatusOK, "admin_modal_users.html", map[string]interface{}{
		"Users":     users,
		"UserCount": userCount,
	})
}

// HandleModalBulk renders the bulk actions modal fragment.
func (h *AdminHandler) HandleModalBulk(c echo.Context) error {
	return c.Render(http.StatusOK, "admin_modal_bulk.html", nil)
}

// HandleModalCategory renders the category management modal fragment.
func (h *AdminHandler) HandleModalCategory(c echo.Context) error {
	ctx := c.Request().Context()
	categories, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		return ui.RespondError(c, err)
	}

	return c.Render(http.StatusOK, "admin_modal_category.html", map[string]interface{}{
		"Categories": categories,
	})
}

// HandleModalModeration renders the claim moderation modal fragment.
func (h *AdminHandler) HandleModalModeration(c echo.Context) error {
	ctx := c.Request().Context()
	claimRequests, err := h.App.DB.GetPendingClaimRequests(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}

	return c.Render(http.StatusOK, "admin_modal_moderation.html", map[string]interface{}{
		"ClaimRequests": claimRequests,
	})
}
