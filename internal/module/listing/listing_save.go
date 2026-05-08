package listing

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleSaveToggle(c echo.Context) error {
	u, err := user.RequireUserAPI(c)
	if err != nil {
		return err
	}
	id := c.Param("id")
	ctx := c.Request().Context()

	saved, err := h.App.DB.IsListingSaved(ctx, u.ID, id)
	if err != nil {
		return ui.RespondError(c, err)
	}

	if saved {
		err = h.App.DB.UnsaveListing(ctx, u.ID, id)
	} else {
		err = h.App.DB.SaveListing(ctx, u.ID, id)
	}
	if err != nil {
		return ui.RespondError(c, err)
	}

	vm := SaveButtonViewModel{
		ListingID: id,
		IsSaved:   !saved,
	}

	return h.RenderTyped(c, "save_button", vm)
}

func (h *ListingHandler) HandleSavedListings(c echo.Context) error {
	u, err := user.RequireUserAPI(c)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	savedListings, err := h.App.DB.GetSavedListings(ctx, u.ID)
	if err != nil {
		return ui.RespondError(c, err)
	}

	// Fetch full listing objects for saved IDs
	var listings []domain.Listing
	savedIDs := make(map[string]bool)
	for _, sl := range savedListings {
		l, err := h.App.DB.FindByID(ctx, sl.ListingID)
		if err != nil {
			continue // skip deleted listings
		}
		listings = append(listings, l)
		savedIDs[sl.ListingID] = true
	}

	vm := SavedListingsViewModel{
		BaseViewData: h.PopulateBase(c),
		Listings:     listings,
		SavedIDs:     savedIDs,
		Source:       "saved",
		Pagination: Pagination{
			TotalCount: len(listings),
		},
	}

	return h.RenderTyped(c, "listing_list", vm)
}
