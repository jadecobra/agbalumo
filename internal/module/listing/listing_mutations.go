package listing

import (
	"github.com/jadecobra/agbalumo/internal/handler"

	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// Create Handler
func (h *ListingHandler) HandleCreate(c echo.Context) error {
	l := domain.Listing{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		IsActive:  true,
		Status:    domain.ListingStatusApproved,
	}

	if err := h.bindAndMapListing(c, &l); err != nil {
		if common.IsImageError(err) {
			return common.RenderImageErrorToast(c, err)
		}
		return handler.RespondError(c, err)
	}

	u := c.Get("User")
	if u == nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required to post listings"))
	}
	user := u.(domain.User)
	l.OwnerID = user.ID

	// Check for duplicate title
	existing, err := h.ListingStore.FindByTitle(c.Request().Context(), l.Title)
	if err == nil && len(existing) > 0 {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Title already exists. Please choose a different title."))
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

	u := c.Get("User")
	if u == nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login Required"))
	}
	user := u.(domain.User)

	ctx := c.Request().Context()
	listing, err := h.ListingStore.FindByID(ctx, id)
	if err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	// Authorization Check (Owner or Admin)
	if listing.OwnerID != user.ID && user.Role != domain.UserRoleAdmin {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	// Save original image URL BEFORE bindAndMapListing may modify it
	originalImageURL := listing.ImageURL

	err = h.bindAndMapListing(c, &listing)
	if err != nil {
		if common.IsImageError(err) {
			return common.RenderImageErrorToast(c, err)
		}
		return handler.RespondError(c, err)
	}

	// Handle Image Removal
	var req ListingFormRequest
	_ = c.Bind(&req)
	if originalImageURL != "" && (req.RemoveImage || listing.ImageURL != originalImageURL) {
		err = h.ImageService.DeleteImage(ctx, originalImageURL)
		if err != nil {
			c.Logger().Errorf("Failed to delete image: %v", err)
		}
		if req.RemoveImage {
			listing.ImageURL = ""
		}
	}

	// Check for duplicate title
	existing, fErr := h.ListingStore.FindByTitle(ctx, listing.Title)
	if fErr == nil {
		// bounded action: title duplicates are usually few
		for _, ext := range existing {
			if ext.ID != listing.ID {
				return handler.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Title already exists. Please choose a different title."))
			}
		}
	}

	return h.processAndSave(c, &listing)
}

func (h *ListingHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	user, ok := handler.GetUser(c)
	if !ok {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
	}

	ctx := c.Request().Context()
	listing, err := h.ListingStore.FindByID(ctx, id)
	if err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	if listing.OwnerID != user.ID {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	if err := h.ListingStore.Delete(ctx, id); err != nil {
		return handler.RespondError(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/profile")
}

func (h *ListingHandler) processAndSave(c echo.Context, l *domain.Listing) error {
	// Auto-populate city from address if missing
	if l.City == "" && l.Address != "" {
		if h.GeocodingSvc != nil {
			city, err := h.GeocodingSvc.GetCity(c.Request().Context(), l.Address)
			if err == nil && city != "" {
				l.City = city
			} else if err != nil {
				c.Logger().Errorf("Failed to geocode address: %v", err)
			}
		}

		// Last resort: Try to extract city manually if still empty
		if l.City == "" {
			l.City = domain.ExtractCityFromAddress(l.Address)
		}
	}

	if err := l.Validate(); err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	if err := h.ListingStore.Save(c.Request().Context(), *l); err != nil {
		return handler.RespondError(c, err)
	}

	var user interface{}
	if u := c.Get("User"); u != nil {
		user = u
	}

	// Trigger an HTMX event so other components (like admin table rows) can update themselves
	c.Response().Header().Add("HX-Trigger", fmt.Sprintf("listing-updated-%s", l.ID))

	// If the request came from the admin dashboard, return no content and let the HX-Trigger handle updates
	if c.QueryParam("source") == "admin" {
		return c.NoContent(http.StatusOK)
	}

	return h.renderWithBaseContext(c, "listing_card", map[string]interface{}{
		"Listing": l,
		"User":    user,
	})
}

func (h *ListingHandler) handleImageUpload(c echo.Context, l *domain.Listing) error {
	imageURL, err := h.ImageService.UploadImage(c.Request().Context(), h.getFileHeader(c, "image"), l.ID)
	if err == nil && imageURL != "" {
		l.ImageURL = fmt.Sprintf("%s?t=%d", imageURL, time.Now().Unix())
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
