package listing_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	cityHouston = "Houston"
)

type filterFallbackCaptureRenderer struct {
	CapturedData interface{}
}

func (r *filterFallbackCaptureRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	r.CapturedData = data
	return nil
}

func TestHandleHome_FilterFallback(t *testing.T) {
	t.Parallel()

	// 1. Setup real test database env
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	// 2. Save a Service listing and Service/Food categories
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: "Service", Name: "Services", Active: true})
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	testutil.SaveTestListing(t, env.App.DB, "service-1", "Houston Hair braiding", func(l *domain.Listing) {
		l.Type = domain.Service
		l.City = cityHouston
		l.Status = domain.ListingStatusApproved
		l.IsActive = true
	})

	// 3. Create context with type=Food (which returns 0 results initially)
	c, _ := testutil.SetupModuleContext(http.MethodGet, "/?type=Food", nil)

	// Set custom capture renderer
	renderer := &filterFallbackCaptureRenderer{}
	c.Echo().Renderer = renderer

	// 4. Execute handler
	err := h.HandleHome(c)
	assert.NoError(t, err)

	// 5. Assert fallback occurred: listings contains the service listing, and FilterType is cleared to ""
	vm, ok := renderer.CapturedData.(listing.HomeViewModel)
	if assert.True(t, ok, "Captured data should be HomeViewModel") {
		assert.Equal(t, "", vm.FilterType, "FilterType should be cleared to All ('') because Food yielded 0 results")
		assert.Len(t, vm.Listings, 1, "Should have 1 listing due to fallback to all categories")
		assert.Equal(t, "service-1", vm.Listings[0].ID)
	}
}

func TestHandleFragment_FilterFallback(t *testing.T) {
	t.Parallel()

	// 1. Setup real test database env
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	// 2. Save a Service listing and Service/Food categories
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: "Service", Name: "Services", Active: true})
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	testutil.SaveTestListing(t, env.App.DB, "service-1", "Houston Hair braiding", func(l *domain.Listing) {
		l.Type = domain.Service
		l.City = cityHouston
		l.Status = domain.ListingStatusApproved
		l.IsActive = true
	})

	// 3. Create fragment request with type=Food (returns 0 results initially)
	c, _ := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?type=Food&page=1", nil)
	c.Request().Header.Set("HX-Request", "true")

	// Set custom capture renderer
	renderer := &filterFallbackCaptureRenderer{}
	c.Echo().Renderer = renderer

	// 4. Execute handler
	err := h.HandleFragment(c)
	assert.NoError(t, err)

	// 5. Assert fallback occurred: listings contains the service listing, and FilterType is cleared to ""
	vm, ok := renderer.CapturedData.(listing.ListingFragmentViewModel)
	if assert.True(t, ok, "Captured data should be ListingFragmentViewModel") {
		assert.Equal(t, "", vm.FilterType, "FilterType should be cleared to All ('') because Food yielded 0 results")
		assert.Len(t, vm.Listings, 1, "Should have 1 listing due to fallback to all categories")
		assert.Equal(t, "service-1", vm.Listings[0].ID)
	}
}
