package listing

import (
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	uRaw, err := user.RequireUser(c)
	if err != nil || uRaw == nil {
		return err
	}

	p := GetPagination(c, 50)
	listings, _, err := h.App.DB.FindAllByOwner(c.Request().Context(), uRaw.ID, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}

	savedIDs := h.getSavedIDs(c)
	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := ProfileViewModel{
		BaseViewData:     h.PopulateBase(c),
		Listings:         listings,
		SavedIDs:         savedMap,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.RenderTyped(c, "modal_profile", vm)
	}

	return h.RenderTyped(c, "profile.html", vm)
}
