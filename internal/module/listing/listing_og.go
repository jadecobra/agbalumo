package listing

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/service/ogcard"
	"github.com/labstack/echo/v4"
)

// HandleOGImage returns a 1200x630 branded Open Graph card image for the listing.
// If the listing does not exist, a generic fallback card is returned with HTTP 200.
func (h *ListingHandler) HandleOGImage(c echo.Context) error {
	id := c.Param("id")
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Cache-Control", "public, max-age=86400")

	listing, err := h.App.DB.FindByID(c.Request().Context(), id)
	if err != nil || listing.ID == "" {
		fallback := h.ogSvc.GenerateFallback()
		return c.Blob(http.StatusOK, "image/png", fallback)
	}

	cardData := ogcard.CardData{
		UpdatedAt:   listing.CreatedAt,
		Title:       listing.Title,
		City:        listing.City,
		Category:    string(listing.Type),
		Rating:      listing.Rating,
		ReviewCount: listing.ReviewCount,
	}
	if listing.RatingUpdatedAt != nil && listing.RatingUpdatedAt.After(cardData.UpdatedAt) {
		cardData.UpdatedAt = *listing.RatingUpdatedAt
	}

	pngBytes, err := h.ogSvc.Generate(cardData)
	if err != nil {
		fallback := h.ogSvc.GenerateFallback()
		return c.Blob(http.StatusOK, "image/png", fallback)
	}

	return c.Blob(http.StatusOK, "image/png", pngBytes)
}

// SetOGCardService allows injecting a custom or mock generator for testing.
func (h *ListingHandler) SetOGCardService(s ogcard.Generator) {
	h.ogSvc = s
}
