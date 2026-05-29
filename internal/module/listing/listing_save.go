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

	clickedPrefix := c.QueryParam("id_prefix")

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	prefixes := []string{"", "fallback-", "modal-", "modal-saved-", "detail-"}
	for _, prefix := range prefixes {
		var classes, textColorClass string
		if prefix == "detail-" {
			classes = "relative transition-all duration-200"
			textColorClass = "text-white hover:text-red-400"
		} else {
			classes = "absolute top-2 right-2 z-30 transition-all duration-200 rounded-none"
			textColorClass = "text-white hover:text-red-400"
		}

		isClicked := prefix == clickedPrefix
		vm := SaveButtonViewModel{
			ListingID:      id,
			IsSaved:        !saved,
			Classes:        classes,
			TextColorClass: textColorClass,
			IDPrefix:       prefix,
			OOB:            !isClicked,
		}

		if err := c.Echo().Renderer.Render(c.Response().Writer, "save_button", vm, c); err != nil {
			return err
		}
	}
	return nil
}

func (h *ListingHandler) HandleSavedListings(c echo.Context) error {
	u, err := user.RequireUserAPI(c)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	params := h.parseQueryParams(c)
	lat, lng := h.resolveCoordinates(ctx, &params)

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
		UserLat: lat,
		UserLng: lng,
	}

	return h.RenderTyped(c, "listing_list", vm)
}
