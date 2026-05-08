package admin

import (
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

	vm := AdminModalChartsViewModel{
		BaseViewData:  h.PopulateBase(c),
		ListingGrowth: listingGrowth,
		UserGrowth:    userGrowth,
	}

	return h.RenderTyped(c, "admin_modal_charts.html", vm)
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

	vm := AdminModalUsersViewModel{
		BaseViewData: h.PopulateBase(c),
		Users:        users,
		UserCount:    userCount,
	}

	return h.RenderTyped(c, "admin_modal_users.html", vm)
}

// HandleModalBulk renders the bulk actions modal fragment.
func (h *AdminHandler) HandleModalBulk(c echo.Context) error {
	vm := AdminModalBulkViewModel{
		BaseViewData: h.PopulateBase(c),
	}
	return h.RenderTyped(c, "admin_modal_bulk.html", vm)
}

// HandleModalCategory renders the category management modal fragment.
func (h *AdminHandler) HandleModalCategory(c echo.Context) error {
	ctx := c.Request().Context()
	categories, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		return ui.RespondError(c, err)
	}

	vm := AdminModalCategoryViewModel{
		BaseViewData: h.PopulateBase(c),
		Categories:   categories,
	}

	return h.RenderTyped(c, "admin_modal_category.html", vm)
}

// HandleModalModeration renders the claim moderation modal fragment.
func (h *AdminHandler) HandleModalModeration(c echo.Context) error {
	ctx := c.Request().Context()
	claimRequests, err := h.App.DB.GetPendingClaimRequests(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}

	vm := AdminModalModerationViewModel{
		BaseViewData:  h.PopulateBase(c),
		ClaimRequests: claimRequests,
	}

	return h.RenderTyped(c, "admin_modal_moderation.html", vm)
}
