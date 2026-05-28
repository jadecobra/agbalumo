package listing

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	uRaw, err := user.RequireUserAPI(c)
	if err != nil || uRaw == nil {
		return err
	}

	p := GetPagination(c, 50)
	listings, _, err := h.App.DB.FindAllByOwner(c.Request().Context(), uRaw.ID, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondError(c, err)
	}

	// Load favorited/saved listings for profile "Saved" section (best-effort; do not break profile if saved load has issues)
	savedRecords, err := h.App.DB.GetSavedListings(c.Request().Context(), uRaw.ID)
	if err != nil {
		savedRecords = nil // profile still renders posted
	}
	var savedListings []domain.Listing
	for _, sl := range savedRecords {
		l, err := h.App.DB.FindByID(c.Request().Context(), sl.ListingID)
		if err != nil {
			continue // skip deleted
		}
		savedListings = append(savedListings, l)
	}

	savedIDs := h.getSavedIDs(c)
	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := ProfileViewModel{
		BaseViewData:     h.PopulateBase(c),
		Listings:         listings,
		SavedListings:    savedListings,
		SavedIDs:         savedMap,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.RenderTyped(c, "modal_profile", vm)
	}

	return h.RenderTyped(c, "profile.html", vm)
}
