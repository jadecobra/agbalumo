package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/moderator"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	Repo domain.ListingRepository
}

func NewListingHandler(repo domain.ListingRepository) *ListingHandler {
	return &ListingHandler{Repo: repo}
}

// Home Handler
func (h *ListingHandler) HandleHome(c echo.Context) error {
	ctx := c.Request().Context()
	listings, err := h.Repo.FindAll(ctx, "", "", false)
	if err != nil {
		return RespondError(c, err)
	}

	counts, err := h.Repo.GetCounts(ctx)
	if err != nil {
		// We log the error but proceed with empty counts to avoid crashing the home page
		// if just the counts query fails for some reason.
		c.Logger().Errorf("failed to get listing counts: %v", err)
		counts = make(map[domain.Category]int)
	}

	strCounts := make(map[string]int)
	totalCount := 0
	for cat, count := range counts {
		strCounts[string(cat)] = count
		totalCount += count
	}

	user := c.Get("User")

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Listings":         listings,
		"Counts":           strCounts,
		"TotalCount":       totalCount,
		"User":             user,
		"GoogleMapsApiKey": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}



// Fragment Handler (HTMX)
func (h *ListingHandler) HandleFragment(c echo.Context) error {
	filterType := c.QueryParam("type")
	queryText := c.QueryParam("q")
	listings, err := h.Repo.FindAll(c.Request().Context(), filterType, queryText, false)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "listing_list.html", map[string]interface{}{
		"Listings": listings,
		"User":     c.Get("User"),
	})
}

// Detail Handler
func (h *ListingHandler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	listing, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	user := c.Get("User")

	return c.Render(http.StatusOK, "modal_detail.html", map[string]interface{}{
		"Listing": listing,
		"User":    user,
	})
}

// HandleEdit renders the edit modal
func (h *ListingHandler) HandleEdit(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("User").(domain.User)

	listing, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	// Authorization Check
	if listing.OwnerID != user.ID {
		return c.String(http.StatusForbidden, "You are not the owner of this listing")
	}

	return c.Render(http.StatusOK, "modal_edit_listing.html", map[string]interface{}{
		"Listing":          listing,
		"GoogleMapsApiKey": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

type ListingFormRequest struct {
	Title           string `form:"title"`
	Type            string `form:"type"`
	OwnerOrigin     string `form:"owner_origin"`
	Description     string `form:"description"`
	City            string `form:"city"`
	Address         string `form:"address"` // New
	HoursOfOperation string `form:"hours_of_operation"` // New
	ContactEmail    string `form:"contact_email"`
	ContactPhone    string `form:"contact_phone"`
	ContactWhatsApp string `form:"contact_whatsapp"`
	WebsiteURL      string `form:"website_url"`
	DeadlineDate    string `form:"deadline_date"`
	EventStart      string `form:"event_start"`
	EventEnd        string `form:"event_end"`
	Skills          string `form:"skills"`
	JobStartDate    string `form:"job_start_date"`
	JobApplyURL     string `form:"job_apply_url"`
	Company         string `form:"company"`
	PayRange        string `form:"pay_range"`
}

// Create Handler
func (h *ListingHandler) HandleCreate(c echo.Context) error {
	var req ListingFormRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Request")
	}

	l := domain.Listing{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	if err := h.populateListingFromRequest(c, &l, req); err != nil {
		return err
	}

	// Assign Owner if authenticated
	if u := c.Get("User"); u != nil {
		if user, ok := u.(domain.User); ok {
			l.OwnerID = user.ID
		}
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

	// Explicit Auth check
	u := c.Get("User")
	if u == nil {
		return c.String(http.StatusUnauthorized, "Login Required")
	}
	user := u.(domain.User)

	ctx := c.Request().Context()
	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	// Authorization Check
	if listing.OwnerID != user.ID {
		return c.String(http.StatusForbidden, "You are not the owner of this listing")
	}

	var req ListingFormRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Request")
	}

	if err := h.populateListingFromRequest(c, &listing, req); err != nil {
		return err
	}

	return h.processAndSave(c, &listing)
}

func (h *ListingHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("User").(domain.User)

	ctx := c.Request().Context()
	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}

	if listing.OwnerID != user.ID {
		return c.String(http.StatusForbidden, "You are not the owner of this listing")
	}

	if err := h.Repo.Delete(ctx, id); err != nil {
		return RespondError(c, err)
	}

	// For HTMX requests, we might want to just remove the element or return a message.
	// If it's a full page reload or generic request, redirecting to profile is safe.
	// Let's assume HTMX usage where we might want to return nothing (empty 200) to delete the element,
	// or redirect if we are on the detail page.
	// Simplest for now: Redirect to profile.
	return c.Redirect(http.StatusSeeOther, "/profile")
}

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}
	u := user.(domain.User)

	listings, err := h.Repo.FindAllByOwner(c.Request().Context(), u.ID)
	if err != nil {
		return RespondError(c, err)
	}

	return c.Render(http.StatusOK, "modal_profile.html", map[string]interface{}{
		"User":             u,
		"Listings":         listings,
		"GoogleMapsApiKey": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

// Helper methods

func (h *ListingHandler) populateListingFromRequest(c echo.Context, l *domain.Listing, req ListingFormRequest) error {
	l.Title = req.Title
	l.Type = domain.Category(req.Type)
	l.OwnerOrigin = req.OwnerOrigin
	l.Description = req.Description
	l.City = req.City
	l.City = req.City
	l.Address = req.Address
	l.HoursOfOperation = req.HoursOfOperation
	l.ContactEmail = req.ContactEmail
	l.ContactPhone = req.ContactPhone
	l.ContactWhatsApp = req.ContactWhatsApp
	l.ContactWhatsApp = req.ContactWhatsApp
	l.WebsiteURL = req.WebsiteURL
	l.Skills = req.Skills
	l.JobApplyURL = req.JobApplyURL
	l.Company = req.Company
	l.PayRange = req.PayRange

	// Handle Image Upload
	if imageURL, err := h.saveUploadedImage(c, l.ID); err == nil && imageURL != "" {
		l.ImageURL = imageURL
	} else if err != nil {
		return RespondError(c, err)
	}

	// Handle Deadline
	if l.Type == domain.Request && req.DeadlineDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.DeadlineDate)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Date Format")
		}
		l.Deadline = parsedTime
	}
	
	// Handle Event Dates
	if l.Type == domain.Event {
		if req.EventStart != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventStart)
			if err != nil {
				return c.String(http.StatusBadRequest, "Invalid Start Date Format")
			}

			l.EventStart = parsedTime
		}
		if req.EventEnd != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventEnd)
			if err != nil {
				return c.String(http.StatusBadRequest, "Invalid End Date Format")
			}
			l.EventEnd = parsedTime
		}
	}

	// Handle Job Start Date
	if l.Type == domain.Job && req.JobStartDate != "" {
		parsedTime, err := time.Parse("2006-01-02T15:04", req.JobStartDate)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Job Start Date Format")
		}
		l.JobStartDate = parsedTime
	}


	return nil
}

func (h *ListingHandler) processAndSave(c echo.Context, l *domain.Listing) error {
	// Domain Validation
	if err := l.Validate(); err != nil {
		return c.String(http.StatusBadRequest, "Validation Error: "+err.Error())
	}

	// Save (Optimistic)
	if err := h.Repo.Save(c.Request().Context(), *l); err != nil {
		return RespondError(c, err)
	}

	// Async Moderation
	go func(listing domain.Listing) {
		// Use a detached context since the request context will be cancelled
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mod, err := moderator.NewGeminiModerator(ctx)
		if err != nil {
			// Log error, fail open (allow listing)
			return
		}

		if err := mod.CheckListing(ctx, listing); err != nil {
			// Violation confirmed. Mark as inactive.
			listing.IsActive = false
			if saveErr := h.Repo.Save(ctx, listing); saveErr != nil {
				// We use stdlib log here as Echo context might be invalid
				// Ideally we'd have a logger injected into the handler
			}
		}
	}(*l)

	// Get User for template context
	var user interface{}
	if u := c.Get("User"); u != nil {
		user = u
	}

	return c.Render(http.StatusOK, "listing_card.html", map[string]interface{}{
		"Listing": l,
		"User":    user,
	})
}

func (h *ListingHandler) saveUploadedImage(c echo.Context, listingID string) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", nil // No file uploaded
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Ensure directory exists
	// Note: In production this path should be configurable/absolute
	uploadDir := "ui/static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// We use a simple strategy of listingID.jpg for now, or timestamped for updates if needed.
	// To keep it simple and consistent:
	// For updates, we might start accumulating garbage if we timestamp.
	// But browser caching is real.
	// Let's stick to a simple filename + timestamp param in URL if needed, or just overwrite.
	// To match previous update logic: timestamp was separate.
	// Let's simplify: listingID.jpg. If it exists, overwrite.
	// Browser cache busting can be handled by the frontend adding ?v=...
	filename := listingID + ".jpg"
	dstPath := uploadDir + "/" + filename
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return "/static/uploads/" + filename, nil
}
