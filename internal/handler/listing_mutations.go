package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		if IsImageError(err) {
			return h.renderImageErrorToast(c, err)
		}
		return RespondError(c, err)
	}

	u := c.Get("User")
	if u == nil {
		return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required to post listings"))
	}
	user := u.(domain.User)
	l.OwnerID = user.ID

	// Check for duplicate title
	existing, err := h.Repo.FindByTitle(c.Request().Context(), l.Title)
	if err == nil && len(existing) > 0 {
		return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Title already exists. Please choose a different title."))
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
		return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login Required"))
	}
	user := u.(domain.User)

	ctx := c.Request().Context()
	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	// Authorization Check (Owner or Admin)
	if listing.OwnerID != user.ID && user.Role != domain.UserRoleAdmin {
		return RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	// Save original image URL BEFORE bindAndMapListing may modify it
	originalImageURL := listing.ImageURL

	err = h.bindAndMapListing(c, &listing)
	if err != nil {
		if IsImageError(err) {
			return h.renderImageErrorToast(c, err)
		}
		return RespondError(c, err)
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
	existing, fErr := h.Repo.FindByTitle(ctx, listing.Title)
	if fErr == nil {
		// bounded action: title duplicates are usually few
		for _, ext := range existing {
			if ext.ID != listing.ID {
				return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Title already exists. Please choose a different title."))
			}
		}
	}

	return h.processAndSave(c, &listing)
}

func (h *ListingHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	user, ok := GetUser(c)
	if !ok {
		return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
	}

	ctx := c.Request().Context()
	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	if listing.OwnerID != user.ID {
		return RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	if err := h.Repo.Delete(ctx, id); err != nil {
		return RespondError(c, err)
	}

	return c.Redirect(http.StatusSeeOther, "/profile")
}

func (h *ListingHandler) processAndSave(c echo.Context, l *domain.Listing) error {
	// Auto-populate city from address if missing and GeocodingSvc is available
	if l.City == "" && l.Address != "" && h.GeocodingSvc != nil {
		city, err := h.GeocodingSvc.GetCity(c.Request().Context(), l.Address)
		if err == nil && city != "" {
			l.City = city
		} else if err != nil {
			c.Logger().Errorf("Failed to geocode address: %v", err)
		}
	}

	if err := l.Validate(); err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	if err := h.Repo.Save(c.Request().Context(), *l); err != nil {
		return RespondError(c, err)
	}

	var user interface{}
	if u := c.Get("User"); u != nil {
		user = u
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

// normalizeURL ensures the given string url has a 'http://' or 'https://' prefix.
func normalizeURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return ""
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "https://" + u
	}
	return u
}

func IsImageError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "File size exceeds") ||
		strings.Contains(msg, "Invalid file type") ||
		strings.Contains(msg, "Invalid or unsupported image")
}

func (h *ListingHandler) renderImageErrorToast(c echo.Context, err error) error {
	var msg string
	if he, ok := err.(*echo.HTTPError); ok {
		if m, ok := he.Message.(string); ok {
			msg = m
		}
	}

	toastID := uuid.New().String()

	c.Response().Header().Set("HX-Reswap", "none")
	c.Response().Header().Set("Content-Type", "text/html")

	return c.HTML(http.StatusBadRequest, fmt.Sprintf(`
	<div id="toast-%s" 
	     class="fixed top-4 right-4 z-50 max-w-sm w-full bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-xl shadow-lg p-4 flex items-start gap-3 animate-in slide-in-from-top-2 fade-in"
	     role="alert"
	     hx-on::after-transaction="if(event.detail.failed) setTimeout(() => { const t = document.getElementById('toast-%s'); if(t) { t.style.animation = 'fade-out 0.3s ease-out forwards'; setTimeout(() => t.remove(), 300); } }, 5000)">
	    <span class="material-symbols-outlined text-red-500 text-[20px] mt-0.5">error</span>
	    <div class="flex-1 min-w-0">
	        <p class="text-sm font-medium text-red-800 dark:text-red-200">Image Upload Failed</p>
	        <p class="text-sm text-red-600 dark:text-red-300 mt-1">%s</p>
	    </div>
	    <button hx-on:click="this.parentElement.remove()" 
	            class="text-red-400 hover:text-red-600 dark:hover:text-red-200 transition-colors">
	        <span class="material-symbols-outlined text-[18px]">close</span>
	    </div>`, toastID, toastID, msg))
}
