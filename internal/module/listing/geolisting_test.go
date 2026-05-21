package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockListingStore struct {
	domain.ListingRepository
	mock.Mock
}

func (m *mockListingStore) FindAll(ctx context.Context, filterType, query, city string, lat, lng, radius float64, sortField, sortOrder string, includeInactive bool, limit, offset int) ([]domain.Listing, int, error) {
	args := m.Called(ctx, filterType, query, city, lat, lng, radius, sortField, sortOrder, includeInactive, limit, offset)
	return args.Get(0).([]domain.Listing), args.Int(1), args.Error(2)
}

func (m *mockListingStore) GetFeaturedListings(ctx context.Context, category, city string) ([]domain.Listing, error) {
	return []domain.Listing{}, nil
}

func (m *mockListingStore) GetLocations(ctx context.Context) ([]domain.Location, error) {
	return []domain.Location{}, nil
}

func (m *mockListingStore) GetSavedListingIDs(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (m *mockListingStore) GetCounts(ctx context.Context) (map[domain.Category]int, error) {
	return map[domain.Category]int{}, nil
}

type mockCategorizationService struct {
	mock.Mock
}

func (m *mockCategorizationService) GetActiveCategories(ctx context.Context) ([]domain.CategoryData, error) {
	return []domain.CategoryData{}, nil
}

func (m *mockCategorizationService) GetCategories(ctx context.Context, filter domain.CategoryFilter) ([]domain.CategoryData, error) {
	return []domain.CategoryData{}, nil
}

func TestHandleHome_Geolocation(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	mockStore := new(mockListingStore)
	mockCatSvc := new(mockCategorizationService)
	app := &env.AppEnv{
		DB:                mockStore,
		Cfg:               &config.Config{Env: "test"},
		CategorizationSvc: mockCatSvc,
	}
	h := listing.NewListingHandler(app)

	// Dallas coords
	lat, lng := 32.7767, -96.7970
	radius := 10.0

	mockStore.On("FindAll", mock.Anything, "Food", "Nigerian", "", lat, lng, radius, "", "", false, 30, 0).
		Return([]domain.Listing{{ID: "1", Title: "Naija Kitchen"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/?lat=32.7767&lng=-96.7970&radius=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// We expect this to fail initially because the handler doesn't yet parse lat/lng
	// or inject "Nigerian" query for geolocated requests.
	err := h.HandleHome(c)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}

func TestHandleFragment_NigerianFirstNotExclusive(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	mockStore := new(mockListingStore)
	mockCatSvc := new(mockCategorizationService)
	app := &env.AppEnv{
		DB:                mockStore,
		Cfg:               &config.Config{Env: "test"},
		CategorizationSvc: mockCatSvc,
	}
	h := listing.NewListingHandler(app)

	// Dallas coords
	lat, lng := 32.7767, -96.7970
	radius := 10.0

	mockStore.On("FindAll", mock.Anything, "Food", "", "", lat, lng, radius, "", "", false, 30, 0).
		Return([]domain.Listing{{ID: "1", Title: "Naija Kitchen"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?lat=32.7767&lng=-96.7970&radius=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleFragment(c)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}
