package admin

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type AdminUsersViewModel struct {
	Users []domain.User
	module.BaseViewData
	Pagination domain.Pagination
}

// HandleUsers renders the list of users for admins.
func (h *AdminHandler) HandleUsers(c echo.Context) error {
	ctx := c.Request().Context()
	p := listing.GetPagination(c, 50)
	users, err := h.App.DB.GetAllUsers(ctx, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}
	p.HasNextPage = len(users) == p.Limit

	vm := AdminUsersViewModel{
		BaseViewData: h.PopulateBase(c),
		Users:        users,
		Pagination:   p,
	}

	return h.RenderTyped(c, "admin_users.html", vm)
}
