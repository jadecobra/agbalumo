package handler

import (
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	Repo         domain.ListingStore
	ImageService service.ImageService
	ListingSvc   *service.ListingService
}

func NewListingHandler(repo domain.ListingStore, imageService service.ImageService) *ListingHandler {
	if imageService == nil {
		imageService = service.NewLocalImageService()
	}
	listingSvc := service.NewListingService(repo)
	return &ListingHandler{Repo: repo, ImageService: imageService, ListingSvc: listingSvc}
}

// Home Handler
func (h *ListingHandler) HandleHome(c echo.Context) error {
	ctx := c.Request().Context()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	listings, err := h.Repo.FindAll(ctx, "", "", false, limit, offset)
	if err != nil {
		return RespondError(c, err)
	}
	hasNextPage := len(listings) == limit

	counts, err := h.Repo.GetCounts(ctx)
	if err != nil {
		// if just the counts query fails for some reason.
		c.Logger().Errorf("failed to get listing counts: %v", err)
		counts = make(map[domain.Category]int)
	}

	featured, err := h.Repo.GetFeaturedListings(ctx)
	if err != nil {
		c.Logger().Errorf("failed to get featured listings: %v", err)
		featured = []domain.Listing{} // Graceful fallback
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
		"Page":             page,
		"HasNextPage":      hasNextPage,
		"FeaturedListings": featured,
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

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	listings, err := h.Repo.FindAll(c.Request().Context(), filterType, queryText, false, limit, offset)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}
	hasNextPage := len(listings) == limit

	data := map[string]interface{}{
		"Listings":    listings,
		"Page":        page,
		"HasNextPage": hasNextPage,
		"User":        c.Get("User"),
	}

	// If HTMX request, render only the listing list partial
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "listing_list", data)
	}

	// For non-HTMX requests, render the full home page with the filtered listings
	// This might not be the desired behavior if this handler is strictly for fragments.
	// Reverting to original behavior for non-HTMX requests, or consider redirecting.
	// For now, keeping the original behavior of rendering just the list.
	return c.Render(http.StatusOK, "listing_list", data)
}

// Detail Handler
func (h *ListingHandler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	listing, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	var user domain.User
	var ok bool
	u := c.Get("User")
	if u != nil {
		user, ok = u.(domain.User)
	}

	// Check if the current user can claim this listing
	canClaim := listing.OwnerID == "" && domain.ClaimableTypes[listing.Type]

	isOwner := false
	if ok {
		isOwner = listing.OwnerID == user.ID
	}

	return c.Render(http.StatusOK, "modal_detail", map[string]interface{}{
		"Listing":  listing,
		"User":     user,
		"IsOwner":  isOwner,
		"CanClaim": canClaim,
	})
}

// HandleEdit renders the edit modal
func (h *ListingHandler) HandleEdit(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("User").(domain.User)

	listing, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	// Authorization Check
	if listing.OwnerID != user.ID {
		return RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	return c.Render(http.StatusOK, "modal_edit_listing.html", map[string]interface{}{
		"Listing":          listing,
		"GoogleMapsApiKey": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

// HandleAbout renders the generic about page.
func (h *ListingHandler) HandleAbout(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"User": c.Get("User"),
	})
}

// HandleClaim processes a request to claim an unowned listing.
func (h *ListingHandler) HandleClaim(c echo.Context) error {
	user, ok := c.Get("User").(domain.User)
	if !ok {
		// If not logged in, redirect to login
		// For HTMX, this might need a specific header to trigger client-side redirect,
		// but standard redirect often works if HX-Redirect handled or simple link.
		// Given detail modal handles logic, we stick to standard redirect fallback.
		return c.Redirect(http.StatusFound, "/auth/google/login")
	}

	id := c.Param("id")

	// Call the service layer to perform business logic
	_, err := h.ListingSvc.ClaimListing(c.Request().Context(), user.ID, id)
	if err != nil {
		if err.Error() == "listing not found" {
			return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
		}
		if err.Error() == "listing is already owned" || err.Error() == "listing type cannot be claimed" {
			return RespondError(c, echo.NewHTTPError(http.StatusForbidden, err.Error()))
		}
		return RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to claim listing: "+err.Error()))
	}

	// Success - Render the updated detail modal (HTMX Swap)
	// We reuse HandleDetail logic to re-fetch and render the modal
	return h.HandleDetail(c)
}

type ListingFormRequest struct {
	Title            string `form:"title"`
	Type             string `form:"type"`
	OwnerOrigin      string `form:"owner_origin"`
	Description      string `form:"description"`
	City             string `form:"city"`
	Address          string `form:"address"`            // New
	HoursOfOperation string `form:"hours_of_operation"` // New
	ContactEmail     string `form:"contact_email"`
	ContactPhone     string `form:"contact_phone"`
	ContactWhatsApp  string `form:"contact_whatsapp"`
	WebsiteURL       string `form:"website_url"`
	DeadlineDate     string `form:"deadline_date"`
	EventStart       string `form:"event_start"`
	EventEnd         string `form:"event_end"`
	Skills           string `form:"skills"`
	JobStartDate     string `form:"job_start_date"`
	JobApplyURL      string `form:"job_apply_url"`
	Company          string `form:"company"`
	PayRange         string `form:"pay_range"`
}

// ToListing maps the DTO fields directly to the domain Listing and parses dates.
func (req *ListingFormRequest) ToListing(l *domain.Listing) error {
	l.Title = req.Title
	l.Type = domain.Category(req.Type)
	l.OwnerOrigin = req.OwnerOrigin
	l.Description = req.Description
	l.City = req.City
	l.Address = req.Address
	l.HoursOfOperation = req.HoursOfOperation
	l.ContactEmail = req.ContactEmail
	l.ContactPhone = req.ContactPhone
	l.ContactWhatsApp = req.ContactWhatsApp
	l.WebsiteURL = req.WebsiteURL
	l.Skills = req.Skills
	l.JobApplyURL = req.JobApplyURL
	l.Company = req.Company
	l.PayRange = req.PayRange

	if err := parseDeadline(req, l); err != nil {
		return err
	}
	if err := parseEventDates(req, l); err != nil {
		return err
	}
	if err := parseJobStartDate(req, l); err != nil {
		return err
	}

	return nil
}

// Create Handler
func (h *ListingHandler) HandleCreate(c echo.Context) error {
	l := domain.Listing{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		IsActive:  true,                        // ACTIVE immediately (Post-Moderation)
		Status:    domain.ListingStatusPending, // Marked as Pending for Admin review
	}

	if err := h.bindAndMapListing(c, &l); err != nil {
		return RespondError(c, err)
	}

	// Assign Owner (Auth is now required by middleware)
	u := c.Get("User")
	if u == nil {
		return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required to post listings"))
	}
	user := u.(domain.User)
	l.OwnerID = user.ID

	// Check for duplicate title
	existing, err := h.Repo.FindByTitle(c.Request().Context(), l.Title)
	if err != nil {
		return RespondError(c, err)
	}
	if len(existing) > 0 {
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

	// Explicit Auth check
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

	// Authorization Check
	if listing.OwnerID != user.ID {
		return RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	if err := h.bindAndMapListing(c, &listing); err != nil {
		return RespondError(c, err)
	}

	// Check for duplicate title
	existing, err := h.Repo.FindByTitle(ctx, listing.Title)
	if err != nil {
		return RespondError(c, err)
	}
	for _, ext := range existing {
		if ext.ID != listing.ID {
			return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Title already exists. Please choose a different title."))
		}
	}

	return h.processAndSave(c, &listing)
}

func (h *ListingHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	user := c.Get("User").(domain.User)

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

	return c.Render(http.StatusOK, "modal_profile", map[string]interface{}{
		"User":             u,
		"Listings":         listings,
		"GoogleMapsApiKey": os.Getenv("GOOGLE_MAPS_API_KEY"),
	})
}

// Helper methods

func (h *ListingHandler) bindAndMapListing(c echo.Context, l *domain.Listing) error {
	var req ListingFormRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Request")
	}

	if err := req.ToListing(l); err != nil {
		return err
	}

	if err := h.handleImageUpload(c, l); err != nil {
		return err
	}

	return nil
}

func (h *ListingHandler) handleImageUpload(c echo.Context, l *domain.Listing) error {
	imageURL, err := h.ImageService.UploadImage(c.Request().Context(), h.getFileHeader(c, "image"), l.ID)
	if err == nil && imageURL != "" {
		l.ImageURL = imageURL
		return nil
	} else if err != nil {
		return RespondError(c, err)
	}
	return nil
}

func parseDeadline(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type == domain.Request && req.DeadlineDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.DeadlineDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Date Format")
		}
		l.Deadline = parsedTime
	}
	return nil
}

func parseEventDates(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type == domain.Event {
		if req.EventStart != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventStart)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid Start Date Format")
			}
			l.EventStart = parsedTime
		}
		if req.EventEnd != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04", req.EventEnd)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid End Date Format")
			}
			l.EventEnd = parsedTime
		}
	}
	return nil
}

func parseJobStartDate(req *ListingFormRequest, l *domain.Listing) error {
	if l.Type == domain.Job && req.JobStartDate != "" {
		parsedTime, err := time.Parse("2006-01-02T15:04", req.JobStartDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid Job Start Date Format")
		}
		l.JobStartDate = parsedTime
	}
	return nil
}

func (h *ListingHandler) processAndSave(c echo.Context, l *domain.Listing) error {
	// Domain Validation
	if err := l.Validate(); err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	// Save
	if err := h.Repo.Save(c.Request().Context(), *l); err != nil {
		return RespondError(c, err)
	}

	// Get User for template context
	var user interface{}
	if u := c.Get("User"); u != nil {
		user = u
	}

	return c.Render(http.StatusOK, "listing_card", map[string]interface{}{
		"Listing": l,
		"User":    user,
	})
}

// Helper to safely get file header
func (h *ListingHandler) getFileHeader(c echo.Context, key string) *multipart.FileHeader {
	file, err := c.FormFile(key)
	if err != nil {
		return nil
	}
	return file
}
