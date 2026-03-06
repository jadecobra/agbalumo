package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	Repo             domain.ListingRepository
	ImageService     service.ImageService
	ListingSvc       *service.ListingService
	GoogleMapsAPIKey string
}

func NewListingHandler(repo domain.ListingRepository, imageService service.ImageService, opts ...string) *ListingHandler {
	var uploadDir string
	if len(opts) > 0 {
		uploadDir = opts[0]
	}
	if imageService == nil {
		imageService = service.NewLocalImageService(uploadDir)
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

	// P1.3: Run all three queries in parallel
	var (
		listings  []domain.Listing
		counts    map[domain.Category]int
		featured  []domain.Listing
		locations []string

		listingsErr  error
		countsErr    error
		featuredErr  error
		locationsErr error

		wg sync.WaitGroup
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		listings, listingsErr = h.Repo.FindAll(ctx, "", "", "", "", false, limit, offset)
	}()
	go func() {
		defer wg.Done()
		counts, countsErr = h.Repo.GetCounts(ctx)
	}()
	go func() {
		defer wg.Done()
		featured, featuredErr = h.Repo.GetFeaturedListings(ctx)
	}()
	go func() {
		defer wg.Done()
		locations, locationsErr = h.Repo.GetLocations(ctx)
	}()
	wg.Wait()

	if listingsErr != nil {
		return RespondError(c, listingsErr)
	}
	hasNextPage := len(listings) == limit

	if countsErr != nil {
		c.Logger().Errorf("failed to get listing counts: %v", countsErr)
		counts = make(map[domain.Category]int)
	}

	if featuredErr != nil {
		c.Logger().Errorf("failed to get featured listings: %v", featuredErr)
		featured = []domain.Listing{} // Graceful fallback
	}

	strCounts := make(map[string]int)
	totalCount := 0
	for cat, count := range counts {
		strCounts[string(cat)] = count
		totalCount += count
	}

	user := c.Get("User")

	if locationsErr != nil {
		c.Logger().Errorf("failed to get locations: %v", locationsErr)
		locations = []string{}
	}

	return h.renderWithBaseContext(c, "index.html", map[string]interface{}{
		"Listings":         listings,
		"Page":             page,
		"HasNextPage":      hasNextPage,
		"FeaturedListings": featured,
		"Counts":           strCounts,
		"Locations":        locations,
		"TotalCount":       totalCount,
		"User":             user,
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
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

	listings, err := h.Repo.FindAll(c.Request().Context(), filterType, queryText, "", "", false, limit, offset)
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
	return c.Render(http.StatusOK, "listing_list", data)
}

// Detail Handler
func (h *ListingHandler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	listing, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	// Fetch category data to check if claimable
	category, _ := h.Repo.GetCategory(ctx, string(listing.Type))

	return c.Render(http.StatusOK, "modal_detail", map[string]interface{}{
		"Listing":          listing,
		"Category":         category,
		"User":             c.Get("User"),
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
	})
}

// HandleEdit renders the edit modal
func (h *ListingHandler) HandleEdit(c echo.Context) error {
	id := c.Param("id")
	user, ok := GetUser(c)
	if !ok {
		return RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
	}

	listing, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
	}

	// Authorization Check (Owner or Admin)
	if listing.OwnerID != user.ID && user.Role != domain.UserRoleAdmin {
		return RespondError(c, echo.NewHTTPError(http.StatusForbidden, "You are not the owner of this listing"))
	}

	return h.renderWithBaseContext(c, "modal_edit_listing", map[string]interface{}{
		"Listing":          listing,
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
	})
}

// HandleAbout renders the generic about page.
func (h *ListingHandler) HandleAbout(c echo.Context) error {
	return h.renderWithBaseContext(c, "about.html", map[string]interface{}{
		"User": c.Get("User"),
	})
}

// HandleClaim processes a request to claim an unowned listing.
func (h *ListingHandler) HandleClaim(c echo.Context) error {
	user, ok := GetUser(c)
	if !ok {
		return c.Redirect(http.StatusFound, "/auth/google/login")
	}

	id := c.Param("id")

	_, err := h.ListingSvc.ClaimListing(c.Request().Context(), *user, id)
	if err != nil {
		switch err.Error() {
		case "listing not found":
			return RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
		case "listing is already owned", "listing type cannot be claimed":
			return RespondError(c, echo.NewHTTPError(http.StatusForbidden, err.Error()))
		case "you already have a pending claim for this listing":
			return RespondError(c, echo.NewHTTPError(http.StatusConflict, err.Error()))
		default:
			return RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit claim: "+err.Error()))
		}
	}

	// HTMX: replace the claim button with a pending-approval notice
	c.Response().Header().Set("Content-Type", "text/html")
	return c.HTML(http.StatusOK, `
		<div class="flex items-center gap-2 bg-earth-accent/10 border border-earth-accent/20 px-3 py-1.5">
			<span class="material-symbols-outlined text-[14px] text-earth-accent">pending</span>
			<span class="text-earth-accent text-xs font-bold uppercase tracking-widest">Claim Pending Review</span>
		</div>`)
}

type ListingFormRequest struct {
	Title            string `form:"title"`
	Type             string `form:"type"`
	OwnerOrigin      string `form:"owner_origin"`
	Description      string `form:"description"`
	City             string `form:"city"`
	Address          string `form:"address"`
	HoursOfOperation string `form:"hours_of_operation"`
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
	RemoveImage      bool   `form:"remove_image"`
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
	l.WebsiteURL = normalizeURL(req.WebsiteURL)
	l.Skills = req.Skills
	l.JobApplyURL = normalizeURL(req.JobApplyURL)
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

func (h *ListingHandler) HandleProfile(c echo.Context) error {
	user := c.Get("User")
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
	}
	u := user.(domain.User)

	p := GetPagination(c, 50)
	listings, err := h.Repo.FindAllByOwner(c.Request().Context(), u.ID, p.Limit, p.Offset)
	if err != nil {
		return RespondError(c, err)
	}

	data := map[string]interface{}{
		"User":             u,
		"Listings":         listings,
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.renderWithBaseContext(c, "modal_profile", data)
	}

	return h.renderWithBaseContext(c, "profile.html", data)
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
		l.ImageURL = fmt.Sprintf("%s?t=%d", imageURL, time.Now().Unix())
		return nil
	} else if err != nil {
		return err
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

func (h *ListingHandler) getFileHeader(c echo.Context, key string) *multipart.FileHeader {
	file, err := c.FormFile(key)
	if err != nil {
		return nil
	}
	return file
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

func (h *ListingHandler) renderWithBaseContext(c echo.Context, tmpl string, data map[string]interface{}) error {
	ctx := c.Request().Context()
	categories, err := h.Repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		c.Logger().Errorf("Failed to retrieve categories: %v", err)
		categories = []domain.CategoryData{}
	}

	data["Categories"] = categories
	return c.Render(http.StatusOK, tmpl, data)
}
