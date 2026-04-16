package listing

import (
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"mime/multipart"
	"net/http"
	"sync"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	App *env.AppEnv
}

func NewListingHandler(app *env.AppEnv) *ListingHandler {
	return &ListingHandler{
		App: app,
	}
}

// RegisterRoutes wires up all HTTP endpoints relating to the Listing domain.
func (h *ListingHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	// Public Routes
	e.GET("/", h.HandleHome)
	e.GET("/listings/fragment", h.HandleFragment)
	e.GET("/listings/:id", h.HandleDetail)
	e.POST("/api/metrics", h.HandleMetricsIngestion)

	// Authenticated Routes
	authGroup := e.Group("", authMw.RequireAuth)
	authGroup.POST(domain.PathListings, h.HandleCreate)
	authGroup.GET("/listings/:id/edit", h.HandleEdit)
	authGroup.PUT("/listings/:id", h.HandleUpdate)
	authGroup.POST("/listings/:id", h.HandleUpdate)
	authGroup.DELETE("/listings/:id", h.HandleDelete)
	authGroup.GET(domain.PathProfile, h.HandleProfile)
	authGroup.POST("/listings/:id/claim", h.HandleClaim)
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
		listings   []domain.Listing
		counts     map[domain.Category]int
		featured   []domain.Listing
		locations  []string
		categories []domain.CategoryData

		listingsErr   error
		countsErr     error
		featuredErr   error
		locationsErr  error
		categoriesErr error

		wg sync.WaitGroup
	)

	var totalCount int
	wg.Add(4)
	go func() {
		defer wg.Done()
		// Default to Food category for the homepage to focus on Ada's primary goal
		listings, totalCount, listingsErr = h.App.DB.FindAll(ctx, string(domain.Food), "", "", "", false, limit, offset)
	}()
	go func() {
		defer wg.Done()
		counts, countsErr = h.App.DB.GetCounts(ctx)
	}()
	go func() {
		defer wg.Done()
		featured, featuredErr = h.App.DB.GetFeaturedListings(ctx, string(domain.Food))
	}()
	go func() {
		defer wg.Done()
		locations, locationsErr = h.App.DB.GetLocations(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		categories, categoriesErr = h.App.CategorizationSvc.GetActiveCategories(ctx)
	}()

	wg.Wait()

	if listingsErr != nil {
		return ui.RespondError(c, listingsErr)
	}
	hasNextPage := offset+len(listings) < totalCount

	h.logError(c, "failed to get listing counts", countsErr)
	h.logError(c, "failed to get featured listings", featuredErr)
	h.logError(c, "failed to get locations", locationsErr)
	h.logError(c, "failed to get categories in HandleHome", categoriesErr)

	strCounts := make(map[string]int)
	for cat, count := range counts {
		strCounts[string(cat)] = count
	}

	u := c.Get("User")

	return h.renderWithBaseContext(c, domain.TemplateIndex, map[string]interface{}{
		"Listings":         listings,
		"Pagination":       Pagination{Page: page, TotalPages: (totalCount + limit - 1) / limit, HasNextPage: hasNextPage, TotalCount: totalCount},
		"FeaturedListings": featured,
		"Counts":           strCounts,
		"Locations":        locations,
		"TotalCount":       totalCount,
		"Categories":       categories,
		"Category":         string(domain.Food),
		"QueryText":        "",
		"User":             u,
		"GoogleMapsApiKey": h.App.Cfg.GoogleMapsAPIKey,
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

	listings, totalCount, err := h.App.DB.FindAll(c.Request().Context(), filterType, queryText, "", "", false, limit, offset)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, err.Error())
	}
	hasNextPage := offset+len(listings) < totalCount

	// For the main feed (no search query), listings already prioritize featured due to SQL optimization.
	finalListings := listings

	data := map[string]interface{}{
		"Listings":   finalListings,
		"Pagination": Pagination{Page: page, TotalPages: (totalCount + limit - 1) / limit, HasNextPage: hasNextPage, TotalCount: totalCount},
		"Category":   filterType,
		"QueryText":  queryText,
		"User":       c.Get("User"),
	}

	// Render the listing list partial (works for both HTMX and full-page requests)
	return c.Render(http.StatusOK, "listing_list", data)
}

// Detail Handler
func (h *ListingHandler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	listing, err := h.App.DB.FindByID(ctx, id)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusNotFound, domain.ErrListingNotFound.Error())
	}

	// Fetch category data to check if claimable
	category, _ := h.App.DB.GetCategory(ctx, string(listing.Type))

	return c.Render(http.StatusOK, "modal_detail", map[string]interface{}{
		"Listing":          listing,
		"Category":         category,
		"User":             c.Get("User"),
		"GoogleMapsApiKey": h.App.Cfg.GoogleMapsAPIKey,
	})
}

// HandleEdit renders the edit modal
func (h *ListingHandler) HandleEdit(c echo.Context) error {
	id := c.Param("id")
	userRaw, ok := user.GetUser(c)
	if !ok {
		return ui.RespondErrorMsg(c, http.StatusUnauthorized, domain.ErrLoginRequired.Error())
	}

	listing, err := h.App.DB.FindByID(c.Request().Context(), id)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusNotFound, domain.ErrListingNotFound.Error())
	}
	// Authorization Check (Owner or Admin)
	if err := h.checkListingAuth(c, listing, userRaw); err != nil {
		return err
	}

	targetID := c.QueryParam("target")
	if targetID == "" {
		targetID = "listing-" + listing.ID
	}
	source := c.QueryParam("source")

	return h.renderWithBaseContext(c, "modal_edit_listing", map[string]interface{}{
		"Listing":          listing,
		"TargetID":         targetID,
		"Source":           source,
		"GoogleMapsApiKey": h.App.Cfg.GoogleMapsAPIKey,
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

	var categories []domain.CategoryData
	var ok bool

	// Check if already provided in data
	if providedCats, exists := data["Categories"]; exists {
		if cats, typeOk := providedCats.([]domain.CategoryData); typeOk {
			categories = cats
			ok = true
		}
	}

	if !ok {
		var err error
		categories, err = h.App.CategorizationSvc.GetActiveCategories(ctx)
		if err != nil {
			c.Logger().Errorf("Failed to retrieve categories: %v", err)
			categories = []domain.CategoryData{}
		}
	}

	data["Categories"] = categories
	data["Env"] = h.App.Cfg.Env
	data["HasGoogleAuth"] = h.App.Cfg.HasGoogleAuth
	return c.Render(http.StatusOK, tmpl, data)
}

func (h *ListingHandler) logError(c echo.Context, msg string, err error) {
	if err != nil {
		c.Logger().Errorf("%s: %v", msg, err)
	}
}
