package handler

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	Repo             domain.ListingRepository
	ImageService     domain.ImageService
	ListingSvc       domain.ListingService
	GeocodingSvc     domain.GeocodingService
	GoogleMapsAPIKey string
	Cfg              *config.Config
}

func NewListingHandler(repo domain.ListingRepository, imageService domain.ImageService, geocodingSvc domain.GeocodingService, cfg *config.Config, opts ...string) *ListingHandler {
	var uploadDir string
	if len(opts) > 0 {
		uploadDir = opts[0]
	}
	if imageService == nil {
		imageService = service.NewLocalImageService(uploadDir)
	}
	listingSvc := service.NewListingService(repo)
	return &ListingHandler{
		Repo:         repo,
		ImageService: imageService,
		ListingSvc:   listingSvc,
		GeocodingSvc: geocodingSvc,
		Cfg:          cfg,
	}
}

// Home Handler
func (h *ListingHandler) HandleHome(c echo.Context) error {
	ctx := c.Request().Context()
	limit := 30
	p := GetPagination(c, limit)
	page := p.Page
	offset := p.Offset

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
	var totalCount int
	go func() {
		defer wg.Done()
		listings, totalCount, listingsErr = h.Repo.FindAll(ctx, "", "", "", "", false, limit, offset)
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
	hasNextPage := offset+len(listings) < totalCount

	if countsErr != nil {
		c.Logger().Errorf("failed to get listing counts: %v", countsErr)
		counts = make(map[domain.Category]int)
	}

	if featuredErr != nil {
		c.Logger().Errorf("failed to get featured listings: %v", featuredErr)
		featured = []domain.Listing{} // Graceful fallback
	}

	// Prioritize featured listings on the first page
	finalListings := listings
	if page == 1 {
		finalListings = h.mergeFeatured(featured, listings, limit)
	}

	strCounts := make(map[string]int)
	categoryTotal := 0
	for cat, count := range counts {
		strCounts[string(cat)] = count
		categoryTotal += count
	}

	user := c.Get("User")

	if locationsErr != nil {
		c.Logger().Errorf("failed to get locations: %v", locationsErr)
		locations = []string{}
	}

	return h.renderWithBaseContext(c, "index.html", map[string]interface{}{
		"Listings":         finalListings,
		"Pagination":       Pagination{Page: page, TotalPages: (totalCount + limit - 1) / limit, HasNextPage: hasNextPage, TotalCount: totalCount},
		"FeaturedListings": featured,
		"Counts":           strCounts,
		"Locations":        locations,
		"TotalCount":       totalCount,
		"Category":         "",
		"QueryText":        "",
		"User":             user,
		"GoogleMapsApiKey": h.GoogleMapsAPIKey,
	})
}

// Fragment Handler (HTMX)
func (h *ListingHandler) HandleFragment(c echo.Context) error {
	filterType := c.QueryParam("type")
	queryText := c.QueryParam("q")

	p := GetPagination(c, 30)
	page := p.Page
	limit := p.Limit
	offset := p.Offset

	listings, totalCount, err := h.Repo.FindAll(c.Request().Context(), filterType, queryText, "", "", false, limit, offset)
	if err != nil {
		return RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
	}
	hasNextPage := offset+len(listings) < totalCount

	// For the first page of the main feed (no filters/search), include featured listings
	finalListings := listings
	if page == 1 && filterType == "" && queryText == "" {
		featured, err := h.Repo.GetFeaturedListings(c.Request().Context())
		if err == nil {
			finalListings = h.mergeFeatured(featured, listings, limit)
		}
	}

	data := map[string]interface{}{
		"Listings":   finalListings,
		"Pagination": Pagination{Page: page, TotalPages: (totalCount + limit - 1) / limit, HasNextPage: hasNextPage, TotalCount: totalCount},
		"Category":   filterType,
		"QueryText":  queryText,
		"User":       c.Get("User"),
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

// Helper methods

func (h *ListingHandler) getFileHeader(c echo.Context, key string) *multipart.FileHeader {
	file, err := c.FormFile(key)
	if err != nil {
		return nil
	}
	return file
}

func (h *ListingHandler) renderWithBaseContext(c echo.Context, tmpl string, data map[string]interface{}) error {
	ctx := c.Request().Context()
	categories, err := h.Repo.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: true})
	if err != nil {
		c.Logger().Errorf("Failed to retrieve categories: %v", err)
		categories = []domain.CategoryData{}
	}

	data["Categories"] = categories
	data["Env"] = h.Cfg.Env
	data["HasGoogleAuth"] = h.Cfg.HasGoogleAuth
	return c.Render(http.StatusOK, tmpl, data)
}

// mergeFeatured prepends featured listings to the list and removes duplicates, keeping total at limit.
func (h *ListingHandler) mergeFeatured(featured, listings []domain.Listing, limit int) []domain.Listing {
	featuredMap := make(map[string]bool)
	for _, f := range featured {
		featuredMap[f.ID] = true
	}

	var filtered []domain.Listing
	for _, l := range listings {
		if !featuredMap[l.ID] {
			filtered = append(filtered, l)
		}
	}

	merged := append(featured, filtered...)
	if len(merged) > limit {
		merged = merged[:limit]
	}
	return merged
}
