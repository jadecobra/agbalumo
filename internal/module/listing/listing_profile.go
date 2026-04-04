package listing

import (
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	u, ok := user.GetUser(c)
	if !ok || u == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}

	p := GetPagination(c, 50)
	listings, _, err := h.ListingStore.FindAllByOwner(c.Request().Context(), u.ID, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}

	data := map[string]interface{}{
		"User":             u,
		"Listings":         listings,
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.renderWithBaseContext(c, "modal_profile", data)
	}

	return h.renderWithBaseContext(c, "profile.html", data)
}
