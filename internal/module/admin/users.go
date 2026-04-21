package admin

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// HandleUsers renders the list of users for admins.
func (h *AdminHandler) HandleUsers(c echo.Context) error {
	ctx := c.Request().Context()
	p := listing.GetPagination(c, 50)
	users, err := h.App.DB.GetAllUsers(ctx, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}
	p.HasNextPage = len(users) == p.Limit

	return c.Render(http.StatusOK, "admin_users.html", map[string]interface{}{
		"Users": users,
		"User":  c.Get(domain.CtxKeyUser),

		"Pagination": p,
	})
}
