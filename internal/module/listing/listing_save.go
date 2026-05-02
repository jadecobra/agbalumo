package listing

import (
	"net/http"

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

	return c.Render(http.StatusOK, "save_button", map[string]interface{}{
		"ListingID": id,
		"IsSaved":   !saved,
	})
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
	for _, sl := range savedListings {
		l, err := h.App.DB.FindByID(ctx, sl.ListingID)
		if err != nil {
			continue // skip deleted listings
		}
		listings = append(listings, l)
	}

	data := h.buildListingViewData(c, listings)
	data["Source"] = "saved"
	return c.Render(http.StatusOK, "listing_list", data)
}
