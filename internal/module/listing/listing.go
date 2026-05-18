package listing

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"mime/multipart"
	"net/http"
	"sync"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
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
	authGroup.POST(domain.PathListingID+"/save", h.HandleSaveToggle)
	authGroup.GET("/saved", h.HandleSavedListings)
}

// Helper to resolve lat/lng coordinates for geo queries
func (h *ListingHandler) resolveCoordinates(ctx context.Context, params *queryParams) (float64, float64) {
	if params.Lat != 0 && params.Lng != 0 {
		return params.Lat, params.Lng
	}
	if params.City != "" && params.Radius > 0 {
		lat, lng, _ := h.App.GeocodingSvc.Geocode(ctx, params.City)
		return lat, lng
	}
	return 0, 0
}

// Home Handler
func (h *ListingHandler) HandleHome(c echo.Context) error {
	ctx := c.Request().Context()
	p := GetPagination(c, 30)
	limit := p.Limit

	params := h.parseQueryParams(c)
	lat, lng := h.resolveCoordinates(ctx, &params)

	// Auto-filter for "Nigerian" if geolocated and no query provided
	if (lat != 0 || lng != 0) && params.Query == "" {
		params.Query = "Nigerian"
	}

	var (
		listings  []domain.Listing
		featured  []domain.Listing
		locations []domain.Location
		savedIDs  []string

		listingsErr  error
		featuredErr  error
		locationsErr error

		wg sync.WaitGroup
	)

	var totalCount int
	wg.Add(3)
	go func() {
		defer wg.Done()
		listings, totalCount, listingsErr = h.App.DB.FindAll(ctx, params.Type, params.Query, params.City, lat, lng, params.Radius, "", "", false, limit, p.Offset)
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
		savedIDs = h.getSavedIDs(c)
	}()

	wg.Wait()

	if listingsErr != nil {
		return ui.RespondError(c, listingsErr)
	}

	h.LogError(c, "failed to get featured listings", featuredErr)
	h.LogError(c, "failed to get locations", locationsErr)

	h.processListings(listings)
	h.processListings(featured)

	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := HomeViewModel{
		BaseViewData: h.PopulateBase(c),
		Listings:     listings,
		Featured:     featured,
		SavedIDs:     savedMap,
		City:         params.City,
		FilterType:   params.Type,
		TotalCount:   totalCount,
		Pagination: Pagination{
			Page:        p.Page,
			Limit:       limit,
			TotalPages:  (totalCount + limit - 1) / limit,
			HasNextPage: p.Offset+len(listings) < totalCount,
			TotalCount:  totalCount,
		},
		Locations:        locations,
		Radius:           params.Radius,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		Source:           c.QueryParam("source"),
		Query:            params.Query,
	}

	if startTS := c.QueryParam("start_ts"); startTS != "" {
		h.App.Logger.Info("Search latency metric", "start_ts", startTS, "now", time.Now().UnixMilli())
	}

	return h.RenderTyped(c, domain.TemplateIndex, vm)
}

// Fragment Handler (HTMX)
func (h *ListingHandler) HandleFragment(c echo.Context) error {
	params := h.parseQueryParams(c)
	p := GetPagination(c, 30)

	lat, lng := h.resolveCoordinates(c.Request().Context(), &params)

	// Auto-filter for "Nigerian" if geolocated and no query provided
	if (lat != 0 || lng != 0) && params.Query == "" {
		params.Query = "Nigerian"
	}

	listings, totalCount, err := h.App.DB.FindAll(c.Request().Context(), params.Type, params.Query, params.City, lat, lng, params.Radius, "", "", false, p.Limit, p.Offset)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, err.Error())
	}

	var featured []domain.Listing
	if p.Page == 1 {
		featured, _ = h.App.DB.GetFeaturedListings(c.Request().Context(), params.Type, params.City)
	}
	h.processListings(listings)
	if len(featured) > 0 {
		h.processListings(featured)
	}

	savedIDs := h.getSavedIDs(c)
	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := ListingFragmentViewModel{
		Listings:   listings,
		Featured:   featured,
		SavedIDs:   savedMap,
		Query:      params.Query,
		City:       params.City,
		FilterType: params.Type,
		Radius:     params.Radius,
		Pagination: Pagination{
			Page:        p.Page,
			Limit:       p.Limit,
			TotalPages:  (totalCount + p.Limit - 1) / p.Limit,
			HasNextPage: p.Offset+len(listings) < totalCount,
			TotalCount:  totalCount,
		},
		User:   c.Get(domain.CtxKeyUser),
		Source: c.QueryParam("source"),
	}

	if startTS := c.QueryParam("start_ts"); startTS != "" {
		h.App.Logger.Info("Search latency metric (fragment)", "start_ts", startTS, "now", time.Now().UnixMilli())
	}

	return h.RenderTyped(c, "listing_list", vm)
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

	canClaim := category.Claimable && listing.OwnerID == ""
	savedIDs := h.getSavedIDs(c)
	savedMap := make(map[string]bool)
	for _, id := range savedIDs {
		savedMap[id] = true
	}

	vm := DetailViewModel{
		BaseViewData:     h.PopulateBase(c),
		Listing:          listing,
		Category:         category,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		CanClaim:         canClaim,
		SavedIDs:         savedMap,
	}

	return h.RenderTyped(c, "modal_detail", vm)
}

// HandleEdit renders the edit modal
func (h *ListingHandler) HandleEdit(c echo.Context) error {
	id := c.Param("id")
	listing, _, err := h.findAndAuthListing(c, id)
	if err != nil {
		return err
	}

	targetID := c.QueryParam(domain.ParamTarget)
	if targetID == "" {
		targetID = "listing-" + listing.ID
	}
	source := c.QueryParam(domain.ParamSource)

	vm := EditViewModel{
		BaseViewData:     h.PopulateBase(c),
		Listing:          listing,
		TargetID:         targetID,
		Source:           source,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
	}

	return h.RenderTyped(c, "modal_edit_listing", vm)
}

// Helper methods

type queryParams struct {
	Type   string
	Query  string
	City   string
	Radius float64
	Lat    float64
	Lng    float64
}

func (h *ListingHandler) parseQueryParams(c echo.Context) queryParams {
	filterType := c.QueryParam(domain.FieldType)
	if filterType == "All" {
		filterType = ""
	} else if filterType == "" {
		filterType = string(domain.Food)
	}

	radius, _ := strconv.ParseFloat(c.QueryParam("radius"), 64)
	lat, _ := strconv.ParseFloat(c.QueryParam("lat"), 64)
	lng, _ := strconv.ParseFloat(c.QueryParam("lng"), 64)

	return queryParams{
		Type:   filterType,
		Query:  c.QueryParam(domain.ParamQuery),
		City:   c.QueryParam(domain.FieldCity),
		Radius: radius,
		Lat:    lat,
		Lng:    lng,
	}
}

func (h *ListingHandler) processListings(listings []domain.Listing) {
	// No-op for now as operational status display is removed
}

func (h *ListingHandler) getSavedIDs(c echo.Context) []string {
	if u := c.Get(domain.CtxKeyUser); u != nil {
		if user, ok := u.(*domain.User); ok {
			savedIDs, _ := h.App.DB.GetSavedListingIDs(c.Request().Context(), user.ID)
			return savedIDs
		}
	}
	return nil
}

func (h *ListingHandler) mapCounts(counts map[domain.Category]int) map[string]int {
	strCounts := make(map[string]int)
	for cat, count := range counts {
		strCounts[string(cat)] = count
	}
	return strCounts
}

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
