package listing

import (
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"mime/multipart"
	"net/http"
	"sync"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	module.BaseHandler
}

func NewListingHandler(app *env.AppEnv) *ListingHandler {
	return &ListingHandler{
		BaseHandler: module.BaseHandler{App: app},
	}
}

// RegisterRoutes wires up all HTTP endpoints relating to the Listing domain.
func (h *ListingHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	// Public Routes
	e.GET("/", h.HandleHome)
	e.GET("/listings/fragment", h.HandleFragment)
	e.GET(domain.PathListingID, h.HandleDetail)
	e.POST("/api/metrics", h.HandleMetricsIngestion)

	// Authenticated Routes
	authGroup := e.Group("", authMw.RequireAuth)
	authGroup.POST(domain.PathListings, h.HandleCreate)
	authGroup.GET(domain.PathListingID+"/edit", h.HandleEdit)
	authGroup.PUT(domain.PathListingID, h.HandleUpdate)
	authGroup.POST(domain.PathListingID, h.HandleUpdate)
	authGroup.DELETE(domain.PathListingID, h.HandleDelete)
	authGroup.GET(domain.PathProfile, h.HandleProfile)
	authGroup.POST(domain.PathListingID+"/claim", h.HandleClaim)
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
		listings, totalCount, listingsErr = h.App.DB.FindAll(ctx, string(domain.Food), "", "", "", "", false, limit, offset)
	}()
	go func() {
		defer wg.Done()
		counts, countsErr = h.App.DB.GetCounts(ctx)
	}()
	go func() {
		defer wg.Done()
		featured, featuredErr = h.App.DB.GetFeaturedListings(ctx, string(domain.Food), "")
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

	h.LogError(c, "failed to get listing counts", countsErr)
	h.LogError(c, "failed to get featured listings", featuredErr)
	h.LogError(c, "failed to get locations", locationsErr)
	h.LogError(c, "failed to get categories in HandleHome", categoriesErr)

	strCounts := make(map[string]int)
	for cat, count := range counts {
		strCounts[string(cat)] = count
	}

	u := c.Get("User")

	return h.RenderWithBaseContext(c, domain.TemplateIndex, map[string]interface{}{
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
	city := c.QueryParam("city")

	p := GetPagination(c, 30)
	page := p.Page
	limit := p.Limit
	offset := p.Offset

	// Ada focus: If location is picked but no category, default to Food
	if city != "" && filterType == "" {
		filterType = string(domain.Food)
	}

	listings, totalCount, err := h.App.DB.FindAll(c.Request().Context(), filterType, queryText, city, "", "", false, limit, offset)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, err.Error())
	}
	hasNextPage := offset+len(listings) < totalCount

	// Fetch featured listings for the selected city and category to support Ada's discovery flow
	featured, _ := h.App.DB.GetFeaturedListings(c.Request().Context(), filterType, city)

	data := map[string]interface{}{
		"Listings":         listings,
		"Pagination":       Pagination{Page: page, TotalPages: (totalCount + limit - 1) / limit, HasNextPage: hasNextPage, TotalCount: totalCount},
		"FeaturedListings": featured,
		"Category":         filterType,
		"QueryText":        queryText,
		"User":             c.Get("User"),
	}

	// Render the listing list partial (works for both HTMX and full-page requests)
	return c.Render(http.StatusOK, "listing_list", data)
}

// Detail Handler
func (h *ListingHandler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	listing, err := h.findListing(c, id)
	if err != nil {
		return err
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
	listing, _, err := h.findAndAuthListing(c, id)
	if err != nil {
		return err
	}

	targetID := c.QueryParam("target")
	if targetID == "" {
		targetID = "listing-" + listing.ID
	}
	source := c.QueryParam("source")

	return h.RenderWithBaseContext(c, "modal_edit_listing", map[string]interface{}{
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

// findListing fetches a listing by ID from the database.
// If the listing does not exist, it writes a 404 response to c and returns echo.ErrNotFound.
// Callers must return the sentinel immediately; the response is already committed.
func (h *ListingHandler) findListing(c echo.Context, id string) (domain.Listing, error) {
	listing, err := h.App.DB.FindByID(c.Request().Context(), id)
	if err != nil {
		_ = ui.RespondErrorMsg(c, http.StatusNotFound, (domain.ErrListingNotFound).Error())
		return domain.Listing{}, echo.ErrNotFound
	}
	return listing, nil
}

// findAndAuthListing combines user requirement, listing retrieval, and authorization check.
func (h *ListingHandler) findAndAuthListing(c echo.Context, id string) (domain.Listing, *domain.User, error) {
	uRaw, err := user.RequireUserAPI(c)
	if err != nil {
		return domain.Listing{}, nil, err
	}
	listing, err := h.findListing(c, id)
	if err != nil {
		return domain.Listing{}, nil, err
	}
	if err := h.checkListingAuth(c, listing, uRaw); err != nil {
		return domain.Listing{}, nil, err
	}
	return listing, uRaw, nil
}

