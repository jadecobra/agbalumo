package listing_test

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFallbackListingStore struct {
	domain.ListingRepository
	mock.Mock
}

func (m *mockFallbackListingStore) FindAll(ctx context.Context, filterType, query, city string, lat, lng, radius float64, sortField, sortOrder string, includeInactive bool, limit, offset int) ([]domain.Listing, int, error) {
	args := m.Called(ctx, filterType, query, city, lat, lng, radius, sortField, sortOrder, includeInactive, limit, offset)
	return args.Get(0).([]domain.Listing), args.Int(1), args.Error(2)
}

func (m *mockFallbackListingStore) GetLocations(ctx context.Context) ([]domain.Location, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Location), args.Error(1)
}

func (m *mockFallbackListingStore) GetFeaturedListings(ctx context.Context, category, city string) ([]domain.Listing, error) {
	args := m.Called(ctx, category, city)
	return args.Get(0).([]domain.Listing), args.Error(1)
}

func (m *mockFallbackListingStore) GetSavedListingIDs(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (m *mockFallbackListingStore) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	return map[domain.Category]int{}, nil
}

func (m *mockFallbackListingStore) GetCategory(ctx context.Context, name string) (domain.CategoryData, error) {
	return domain.CategoryData{}, nil
}

type captureRenderer struct {
	Templates    *template.Template
	CapturedData interface{}
}

func (r *captureRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	r.CapturedData = data
	return r.Templates.ExecuteTemplate(w, name, data)
}

func TestGeolocatedListing_Fallback(t *testing.T) {
	// Active locations to return
	locations := []domain.Location{
		{City: "Dallas", Latitude: 32.7767, Longitude: -96.7970},
		{City: "Houston", Latitude: 29.7604, Longitude: -95.3698},
	}

	dallasListings := []domain.Listing{
		{ID: "dallas-1", Title: "Dallas Naija Grill", City: "Dallas", Latitude: 32.7767, Longitude: -96.7970},
	}

	t.Run("HandleHome: Coordinates close to Dallas -> No Fallback", func(t *testing.T) {
		e := echo.New()
		renderer := &captureRenderer{Templates: testutil.NewMainTemplate()}
		e.Renderer = renderer

		mockStore := new(mockFallbackListingStore)
		mockCatSvc := new(mockCategorizationService)
		app := &env.AppEnv{
			DB:                mockStore,
			Cfg:               &config.Config{Env: "test"},
			CategorizationSvc: mockCatSvc,
		}
		h := listing.NewListingHandler(app)

		// Coords close to Dallas: Lat: 32.77, Lng: -96.79
		lat, lng := 32.77, -96.79
		radius := 10.0

		// FindAll with coordinates yields Dallas listing directly
		mockStore.On("FindAll", mock.Anything, "Food", "", "", lat, lng, radius, "", "", false, 30, 0).
			Return(dallasListings, 1, nil)

		// Mock concurrent calls in HandleHome
		mockStore.On("GetLocations", mock.Anything).Return(locations, nil)
		mockStore.On("GetFeaturedListings", mock.Anything, "Food", "").Return([]domain.Listing{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/?lat=32.77&lng=-96.79&radius=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.HandleHome(c)
		assert.NoError(t, err)

		// Assert on captured data
		vm, ok := renderer.CapturedData.(listing.HomeViewModel)
		if assert.True(t, ok, "Captured data should be HomeViewModel") {
			assert.Equal(t, "", vm.FallbackCity)
			assert.Len(t, vm.Listings, 1)
			assert.Equal(t, "dallas-1", vm.Listings[0].ID)
		}

		mockStore.AssertExpectations(t)
	})

	t.Run("HandleHome: Coordinates far away -> Fallback to Dallas", func(t *testing.T) {
		e := echo.New()
		renderer := &captureRenderer{Templates: testutil.NewMainTemplate()}
		e.Renderer = renderer

		mockStore := new(mockFallbackListingStore)
		mockCatSvc := new(mockCategorizationService)
		app := &env.AppEnv{
			DB:                mockStore,
			Cfg:               &config.Config{Env: "test"},
			CategorizationSvc: mockCatSvc,
		}
		h := listing.NewListingHandler(app)

		// Coords far away: Seattle Lat: 47.6062, Lng: -122.3321
		lat, lng := 47.6062, -122.3321
		radius := 10.0

		// First query (Seattle radius search) returns 0 results
		mockStore.On("FindAll", mock.Anything, "Food", "", "", lat, lng, radius, "", "", false, 30, 0).
			Return([]domain.Listing{}, 0, nil).Once()

		// Fallback query (Seattle radius search with all types) returns 0 results
		mockStore.On("FindAll", mock.Anything, "", "", "", lat, lng, radius, "", "", false, 30, 0).
			Return([]domain.Listing{}, 0, nil).Once()

		// Mock concurrent calls in HandleHome:
		// Reuses locations fetched here without extra trip
		mockStore.On("GetLocations", mock.Anything).Return(locations, nil).Once()
		mockStore.On("GetFeaturedListings", mock.Anything, "Food", "").Return([]domain.Listing{}, nil)

		// Re-query FindAll for Dallas, limit 6 (using fallback type "")
		mockStore.On("FindAll", mock.Anything, "", "", "Dallas", 0.0, 0.0, 0.0, "", "", false, 6, 0).
			Return(dallasListings, 1, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/?lat=47.6062&lng=-122.3321&radius=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.HandleHome(c)
		assert.NoError(t, err)

		// Assert on captured data
		vm, ok := renderer.CapturedData.(listing.HomeViewModel)
		if assert.True(t, ok, "Captured data should be HomeViewModel") {
			assert.Equal(t, "Dallas", vm.FallbackCity)
			assert.Len(t, vm.Listings, 1)
			assert.Equal(t, "dallas-1", vm.Listings[0].ID)
		}

		mockStore.AssertExpectations(t)
	})

	t.Run("HandleFragment: Coordinates far away -> Fallback to Dallas", func(t *testing.T) {
		e := echo.New()
		renderer := &captureRenderer{Templates: testutil.NewMainTemplate()}
		e.Renderer = renderer

		mockStore := new(mockFallbackListingStore)
		mockCatSvc := new(mockCategorizationService)
		app := &env.AppEnv{
			DB:                mockStore,
			Cfg:               &config.Config{Env: "test"},
			CategorizationSvc: mockCatSvc,
		}
		h := listing.NewListingHandler(app)

		// Coords far away: Seattle Lat: 47.6062, Lng: -122.3321
		lat, lng := 47.6062, -122.3321
		radius := 10.0

		// First query (Seattle radius search) returns 0 results
		mockStore.On("FindAll", mock.Anything, "Food", "", "", lat, lng, radius, "", "", false, 30, 0).
			Return([]domain.Listing{}, 0, nil).Once()

		// Fallback query (Seattle radius search with all types) returns 0 results
		mockStore.On("FindAll", mock.Anything, "", "", "", lat, lng, radius, "", "", false, 30, 0).
			Return([]domain.Listing{}, 0, nil).Once()

		// GetFeaturedListings is called in HandleFragment on page 1 (with fallback type "")
		mockStore.On("GetFeaturedListings", mock.Anything, "", "").Return([]domain.Listing{}, nil)

		// GetLocations gets called in fragment empty state check since it's not pre-fetched
		mockStore.On("GetLocations", mock.Anything).Return(locations, nil).Once()

		// Re-query FindAll for Dallas, limit 6 (using fallback type "")
		mockStore.On("FindAll", mock.Anything, "", "", "Dallas", 0.0, 0.0, 0.0, "", "", false, 6, 0).
			Return(dallasListings, 1, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/listings/fragment?lat=47.6062&lng=-122.3321&radius=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.HandleFragment(c)
		assert.NoError(t, err)

		// Assert on captured data
		vm, ok := renderer.CapturedData.(listing.ListingFragmentViewModel)
		if assert.True(t, ok, "Captured data should be ListingFragmentViewModel") {
			assert.Equal(t, "Dallas", vm.FallbackCity)
			assert.Len(t, vm.Listings, 1)
			assert.Equal(t, "dallas-1", vm.Listings[0].ID)
		}

		mockStore.AssertExpectations(t)
	})
}
