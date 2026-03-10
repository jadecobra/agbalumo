package handler

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}
	u := user.(domain.User)

	p := GetPagination(c, 50)
	listings, err := h.Repo.FindAllByOwner(c.Request().Context(), u.ID, p.Limit, p.Offset)
	if err != nil {
		return RespondError(c, err)
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
