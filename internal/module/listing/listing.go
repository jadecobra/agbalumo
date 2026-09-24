package listing

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
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

type homeData struct {
	listings   []domain.Listing
	featured   []domain.Listing
	locations  []domain.Location
	savedIDs   []string
	totalCount int
}

func (h *ListingHandler) fetchHomeData(ctx context.Context, c echo.Context, params *queryParams, lat, lng float64, limit, offset int) (homeData, error, error, error) {
	var (
		data         homeData
		listingsErr  error
		featuredErr  error
		locationsErr error
		wg           sync.WaitGroup
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		data.listings, data.totalCount, listingsErr = h.App.DB.FindAll(ctx, params.Type, params.Query, params.City, lat, lng, params.Radius, "", "", false, limit, offset)
	}()
	go func() {
		defer wg.Done()
		data.featured, featuredErr = h.App.DB.GetFeaturedListings(ctx, string(domain.Food), "")
	}()
	go func() {
		defer wg.Done()
		data.locations, locationsErr = h.App.DB.GetLocations(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data.savedIDs = h.getSavedIDs(c)
	}()

	wg.Wait()
	return data, listingsErr, featuredErr, locationsErr
}

func resolveRadius(lat, lng float64, city string, requestedRadius float64) float64 {
	if (lat != 0 && lng != 0) || city != "" {
		return requestedRadius
	}
	return 0
}

func (h *ListingHandler) resolveHomeListings(ctx context.Context, category string, query string, city string, lat, lng float64, radius float64, limit, offset int, initialListings []domain.Listing, initialCount int) ([]domain.Listing, int, string, error) {
	if initialCount == 0 && category != "" {
		listings, totalCount, err := h.App.DB.FindAll(ctx, "", query, city, lat, lng, radius, "", "", false, limit, offset)
		return listings, totalCount, "", err
	}
	return initialListings, initialCount, category, nil
}

// Home Handler
func (h *ListingHandler) HandleHome(c echo.Context) error {
	ctx := c.Request().Context()
	p := GetPagination(c, 30)
	limit := p.Limit

	params := h.parseQueryParams(c)
	lat, lng := h.resolveCoordinates(ctx, &params)

	data, listingsErr, featuredErr, locationsErr := h.fetchHomeData(ctx, c, &params, lat, lng, limit, p.Offset)
	if listingsErr != nil {
		return ui.RespondError(c, listingsErr)
	}

	listings, totalCount, newType, fallbackErr := h.resolveHomeListings(ctx, params.Type, params.Query, params.City, lat, lng, params.Radius, limit, p.Offset, data.listings, data.totalCount)
	if fallbackErr != nil {
		return ui.RespondError(c, fallbackErr)
	}
	params.Type = newType

	h.LogError(c, "failed to get featured listings", featuredErr)
	h.LogError(c, "failed to get locations", locationsErr)

	h.processListings(listings)
	h.processListings(data.featured)

	var fallbackCity string
	if locationsErr == nil {
		listings, totalCount, fallbackCity = h.resolveFallback(ctx, totalCount, lat, lng, data.locations, listings, &params)
	}

	savedMap := make(map[string]bool)
	for _, id := range data.savedIDs {
		savedMap[id] = true
	}

	vm := HomeViewModel{
		BaseViewData: h.PopulateBase(c),
		Listings:     listings,
		Featured:     data.featured,
		SavedIDs:     savedMap,
		City:         params.City,
		FilterType:   params.Type,
		TotalCount:   totalCount,
		Pagination: domain.Pagination{
			Page:        p.Page,
			Limit:       limit,
			TotalPages:  (totalCount + limit - 1) / limit,
			HasNextPage: p.Offset+len(listings) < totalCount,
			TotalCount:  totalCount,
		},
		Locations:        data.locations,
		Radius:           resolveRadius(lat, lng, params.City, params.Radius),
		UserLat:          lat,
		UserLng:          lng,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		Source:           c.QueryParam("source"),
		Query:            params.Query,
		FallbackCity:     fallbackCity,
	}

	if params.Query != "" && params.City != "" {
		vm.MetaTitle = fmt.Sprintf("%s in %s | agbalumo", capitalize(params.Query), capitalize(params.City))
		vm.MetaDescription = fmt.Sprintf("Find verified %s in %s. Real menus, contact information, and directions in under 60 seconds on agbalumo.", params.Query, params.City)
		vm.MetaURL = fmt.Sprintf("https://agbalumo.com/?q=%s&city=%s", url.QueryEscape(params.Query), url.QueryEscape(params.City))
	} else if params.City != "" {
		vm.MetaTitle = fmt.Sprintf("African Food in %s | agbalumo", capitalize(params.City))
		vm.MetaDescription = fmt.Sprintf("Find top-rated Nigerian and West African restaurants in %s. Authentic dining in under 60 seconds.", params.City)
		vm.MetaURL = fmt.Sprintf("https://agbalumo.com/?city=%s", url.QueryEscape(params.City))
	} else if params.Query != "" {
		vm.MetaTitle = fmt.Sprintf("%s - African Food Search | agbalumo", capitalize(params.Query))
		vm.MetaDescription = fmt.Sprintf("Discover authentic %s spots across the diaspora on agbalumo.", params.Query)
		vm.MetaURL = fmt.Sprintf("https://agbalumo.com/?q=%s", url.QueryEscape(params.Query))
	}

	if startTS := c.QueryParam("start_ts"); startTS != "" {
		h.App.Logger.Info("Search latency metric", "start_ts", startTS, "now", time.Now().UnixMilli())
	}

	h.logZeroResults(ctx, "home", totalCount, params.Query, params.City, lat, lng)

	return h.RenderTyped(c, domain.TemplateIndex, vm)
}

func (h *ListingHandler) logZeroResults(ctx context.Context, source string, totalCount int, query, city string, lat, lng float64) {
	if totalCount == 0 && (query != "" || city != "") {
		h.App.Logger.InfoContext(ctx, "search_zero_results",
			"query", query,
			"city", city,
			"lat", lat,
			"lng", lng,
			"source", source,
		)
	}
}

func (h *ListingHandler) fetchFragmentData(ctx context.Context, c echo.Context, params *queryParams, limit, offset int, includeFeatured bool) (homeData, error) {
	var (
		data         homeData
		listingsErr  error
		featuredErr  error
		locationsErr error
		wg           sync.WaitGroup
	)

	featuredType := params.Type
	featuredCity := params.City

	wg.Add(4)
	go func() {
		defer wg.Done()
		data.listings, data.totalCount, listingsErr = h.App.DB.FindAll(ctx, params.Type, params.Query, params.City, params.Lat, params.Lng, params.Radius, "", "", false, limit, offset)
		if listingsErr == nil && data.totalCount == 0 && params.Type != "" {
			data.listings, data.totalCount, listingsErr = h.App.DB.FindAll(ctx, "", params.Query, params.City, params.Lat, params.Lng, params.Radius, "", "", false, limit, offset)
			if listingsErr == nil {
				params.Type = ""
			}
		}
	}()
	go func() {
		defer wg.Done()
		if includeFeatured {
			data.featured, featuredErr = h.App.DB.GetFeaturedListings(ctx, featuredType, featuredCity)
		}
	}()
	go func() {
		defer wg.Done()
		data.locations, locationsErr = h.App.DB.GetLocations(ctx)
	}()
	go func() {
		defer wg.Done()
		data.savedIDs = h.getSavedIDs(c)
	}()

	wg.Wait()

	_ = featuredErr
	_ = locationsErr

	if listingsErr != nil {
		return homeData{}, listingsErr
	}

	return data, nil
}

// Fragment Handler (HTMX)
func (h *ListingHandler) HandleFragment(c echo.Context) error {
	params := h.parseQueryParams(c)
	p := GetPagination(c, 30)

	lat, lng := h.resolveCoordinates(c.Request().Context(), &params)
	params.Lat = lat
	params.Lng = lng

	data, err := h.fetchFragmentData(c.Request().Context(), c, &params, p.Limit, p.Offset, p.Page == 1)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, err.Error())
	}

	listings := data.listings
	totalCount := data.totalCount

	h.processListings(listings)
	if len(data.featured) > 0 {
		h.processListings(data.featured)
	}

	var fallbackCity string
	if len(data.locations) > 0 {
		listings, totalCount, fallbackCity = h.resolveFallback(c.Request().Context(), totalCount, lat, lng, data.locations, listings, &params)
	}

	savedMap := make(map[string]bool)
	for _, id := range data.savedIDs {
		savedMap[id] = true
	}

	vm := ListingFragmentViewModel{
		Listings:   listings,
		Featured:   data.featured,
		SavedIDs:   savedMap,
		Query:      params.Query,
		City:       params.City,
		FilterType: params.Type,
		Radius:     resolveRadius(lat, lng, params.City, params.Radius),
		UserLat:    lat,
		UserLng:    lng,
		Pagination: domain.Pagination{
			Page:        p.Page,
			Limit:       p.Limit,
			TotalPages:  (totalCount + p.Limit - 1) / p.Limit,
			HasNextPage: p.Offset+len(listings) < totalCount,
			TotalCount:  totalCount,
		},
		User:         c.Get(domain.CtxKeyUser),
		Source:       c.QueryParam("source"),
		FallbackCity: fallbackCity,
	}

	if startTS := c.QueryParam("start_ts"); startTS != "" {
		h.App.Logger.Info("Search latency metric (fragment)", "start_ts", startTS, "now", time.Now().UnixMilli())
	}

	h.logZeroResults(c.Request().Context(), "fragment", totalCount, params.Query, params.City, lat, lng)

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
	vm.MetaTitle = fmt.Sprintf("%s | agbalumo", listing.Title)
	if listing.Description != "" {
		vm.MetaDescription = listing.Description
	} else {
		vm.MetaDescription = fmt.Sprintf("View authentic details, location, and contact information for %s on agbalumo.", listing.Title)
	}
	if listing.ImageURL != "" {
		vm.MetaImage = listing.ImageURL
	}
	vm.MetaURL = fmt.Sprintf("https://agbalumo.com/listings/%s", listing.ID)
	vm.MetaType = "restaurant"

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

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMiles = 3959.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMiles * c
}

func (h *ListingHandler) resolveFallback(ctx context.Context, totalCount int, lat, lng float64, locations []domain.Location, listings []domain.Listing, params *queryParams) ([]domain.Listing, int, string) {
	if totalCount > 0 || lat == 0.0 || lng == 0.0 || len(locations) == 0 {
		return listings, totalCount, ""
	}

	closest, found := findClosestLocation(lat, lng, locations)
	if !found || closest.City == "" {
		return listings, totalCount, ""
	}

	fallbackListings, fallbackCount, fallbackErr := h.App.DB.FindAll(ctx, params.Type, params.Query, closest.City, 0, 0, 0, "", "", false, 6, 0)
	if fallbackErr != nil {
		return listings, totalCount, ""
	}

	h.processListings(fallbackListings)
	return fallbackListings, fallbackCount, closest.City
}

func findClosestLocation(lat, lng float64, locations []domain.Location) (domain.Location, bool) {
	var closest domain.Location
	minDist := -1.0
	for _, loc := range locations {
		if loc.Latitude == 0.0 || loc.Longitude == 0.0 {
			continue
		}
		dist := haversineDistance(lat, lng, loc.Latitude, loc.Longitude)
		if minDist < 0 || dist < minDist {
			minDist = dist
			closest = loc
		}
	}
	return closest, minDist >= 0
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
