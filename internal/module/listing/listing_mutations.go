package listing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

const errTitleExists = "Title already exists. Please choose a different title."
const tmplListingCard = "listing_card"

type ListingCardViewModel struct {
	Listing   *domain.Listing
	SavedIDs  map[string]bool
	IDPrefix  string
	GridClass string
	module.BaseViewData
	UserLat float64
	UserLng float64
}

// Create Handler
func (h *ListingHandler) HandleCreate(c echo.Context) error {
	l := domain.Listing{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		IsActive:  true,
		Status:    domain.ListingStatusApproved,
	}

	if err := h.bindAndMapWithImageCheck(c, &l); err != nil {
		return err
	}

	uRaw, ok := user.GetUser(c)
	if !ok || uRaw == nil {
		return ui.RespondErrorMsg(c, http.StatusUnauthorized, "Authentication required to post listings")
	}
	l.OwnerID = uRaw.ID

	if err := h.checkDuplicateTitle(c.Request().Context(), l.Title, ""); err != nil {
		return ui.RespondError(c, err)
	}

	// Default deadline for requests if not provided
	if l.Type == domain.Request && l.Deadline.IsZero() {
		l.Deadline = l.CreatedAt.Add(90 * 24 * time.Hour).Add(-time.Minute)
	}

	return h.processAndSave(c, &l)
}

// HandleUpdate updates the listing
func (h *ListingHandler) HandleUpdate(c echo.Context) error {
	id := c.Param("id")

	listing, _, err := h.findAndAuthListing(c, id)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	// Save original image URL BEFORE bindAndMapListing may modify it
	originalImageURL := listing.ImageURL

	if err := h.bindAndMapWithImageCheck(c, &listing); err != nil {
		return err
	}

	h.handleImageRemoval(c, &listing, originalImageURL)

	if err := h.checkDuplicateTitle(ctx, listing.Title, listing.ID); err != nil {
		return ui.RespondError(c, err)
	}

	return h.processAndSave(c, &listing)
}

func (h *ListingHandler) checkListingAuth(c echo.Context, listing domain.Listing, uRaw *domain.User) error {
	if listing.OwnerID != uRaw.ID && uRaw.Role != domain.UserRoleAdmin {
		return ui.RespondErrorMsg(c, http.StatusForbidden, "You are not the owner of this listing")
	}
	return nil
}

func (h *ListingHandler) handleImageRemoval(c echo.Context, listing *domain.Listing, originalURL string) {
	var req ListingFormRequest
	_ = c.Bind(&req)
	if originalURL != "" && (req.RemoveImage || listing.ImageURL != originalURL) {
		if err := h.App.ImageSvc.DeleteImage(c.Request().Context(), originalURL); err != nil {
			c.Logger().Errorf("Failed to delete image: %v", err)
		}
		if req.RemoveImage {
			listing.ImageURL = ""
		}
	}
}

func (h *ListingHandler) checkDuplicateTitle(ctx context.Context, title string, currentID string) error {
	existing, err := h.App.DB.FindByTitle(ctx, title)
	if err != nil {
		return nil
	}
	for _, ext := range existing {
		if ext.ID != currentID {
			return echo.NewHTTPError(http.StatusBadRequest, errTitleExists)
		}
	}
	return nil
}

func (h *ListingHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	_, _, err := h.findAndAuthListing(c, id)
	if err != nil {
		return err
	}

	if err := h.App.DB.Delete(c.Request().Context(), id); err != nil {
		return ui.RespondError(c, err)
	}

	return c.Redirect(http.StatusSeeOther, domain.PathProfile)
}

func (h *ListingHandler) bindAndMapWithImageCheck(c echo.Context, l *domain.Listing) error {
	if err := h.bindAndMapListing(c, l); err != nil {
		if common.IsImageError(err) {
			return common.RenderImageErrorToast(c, err)
		}
		return ui.RespondError(c, err)
	}
	return nil
}

func (h *ListingHandler) processAndSave(c echo.Context, l *domain.Listing) error {
	h.autoPopulateLocation(c.Request().Context(), l)

	if err := l.Validate(); err != nil {
		return ui.RespondErrorMsg(c, http.StatusBadRequest, "Validation Error: "+err.Error())
	}

	if err := h.App.DB.Save(c.Request().Context(), *l); err != nil {
		return ui.RespondError(c, err)
	}

	// Trigger an HTMX event so other components (like admin table rows) can update themselves
	c.Response().Header().Add(domain.HeaderHXTrigger, fmt.Sprintf("%s%s", domain.TriggerListingUpdatedPrefix, l.ID))

	// If the request came from the admin dashboard, return no content and let the HX-Trigger handle updates
	if c.QueryParam(domain.ParamSource) == domain.ParamSourceAdmin {
		return c.NoContent(http.StatusOK)
	}

	savedIDs := h.getSavedIDs(c)
	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := ListingCardViewModel{
		BaseViewData: h.PopulateBase(c),
		Listing:      l,
		SavedIDs:     savedMap,
		IDPrefix:     "", // Default for home/profile cards
	}

	return h.RenderTyped(c, tmplListingCard, vm)
}

func (h *ListingHandler) autoPopulateLocation(ctx context.Context, l *domain.Listing) {
	if l.Address == "" {
		return
	}

	// 1. Try Google Geocoding if available
	if h.App.GeocodingSvc != nil && l.City == "" {
		city, err := h.App.GeocodingSvc.GetCity(ctx, l.Address)
		if err == nil && city != "" {
			l.City = city
		}
	}

	// 2. Fetch coordinates via GeocodingSvc if available
	if h.App.GeocodingSvc != nil && (l.Latitude == 0.0 || l.Longitude == 0.0) {
		lat, lng, err := h.App.GeocodingSvc.Geocode(ctx, l.Address)
		if err == nil && lat != 0.0 && lng != 0.0 {
			l.Latitude = lat
			l.Longitude = lng
		} else if l.City != "" {
			lat, lng, err = h.App.GeocodingSvc.Geocode(ctx, l.City)
			if err == nil && lat != 0.0 && lng != 0.0 {
				l.Latitude = lat
				l.Longitude = lng
			}
		}
	}

	// 3. Fallbacks using local extraction
	if l.City == "" {
		l.City = domain.ExtractCityFromAddress(l.Address)
	}
	if l.State == "" {
		l.State = domain.ExtractStateFromAddress(l.Address)
	}
	if l.Country == "" {
		l.Country = domain.ExtractCountryFromAddress(l.Address)
	}
}

func (h *ListingHandler) handleImageUpload(c echo.Context, l *domain.Listing) error {
	imageURL, err := h.App.ImageSvc.UploadImage(c.Request().Context(), h.getFileHeader(c, "image"), l.ID)
	if err == nil && imageURL != "" {
		l.ImageURL = fmt.Sprintf("%s?t=%d", imageURL, time.Now().Unix())
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
