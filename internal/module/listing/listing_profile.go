package listing

import (
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	u, err := user.RequireUser(c)
	if err != nil || u == nil {
		return err
	}

	p := GetPagination(c, 50)
	listings, _, err := h.App.DB.FindAllByOwner(c.Request().Context(), u.ID, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}

	data := map[string]interface{}{
		"User":             u,
		"Listings":         listings,
		"GoogleMapsApiKey": h.App.Cfg.GoogleMapsAPIKey,
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.renderWithBaseContext(c, "modal_profile", data)
	}

	return h.renderWithBaseContext(c, "profile.html", data)
}
