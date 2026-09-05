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
	"github.com/stretchr/testify/require"
)

type metaCaptureRenderer struct {
	CapturedData interface{}
}

func (r *metaCaptureRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	r.CapturedData = data
	return nil
}

func TestHandleHome_DynamicMetaTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		queryURL      string
		expectedTitle string
		expectedDesc  string
		expectedURL   string
	}{
		{
			name:          "home_default",
			queryURL:      "/",
			expectedTitle: "agbalumo - Find African Food in <60s",
			expectedDesc:  "Verified Nigerian and West African food in Dallas and beyond. Authentic dishes, real contact details, and directions in under 60 seconds.",
			expectedURL:   "https://agbalumo.com",
		},
		{
			name:          "filter_by_city",
			queryURL:      "/?city=Dallas",
			expectedTitle: "African Food in Dallas | agbalumo",
			expectedDesc:  "Find top-rated Nigerian and West African restaurants in Dallas. Authentic dining in under 60 seconds.",
			expectedURL:   "https://agbalumo.com/?city=Dallas",
		},
		{
			name:          "search_query",
			queryURL:      "/?q=suya",
			expectedTitle: "Suya - African Food Search | agbalumo",
			expectedDesc:  "Discover authentic suya spots across the diaspora on agbalumo.",
			expectedURL:   "https://agbalumo.com/?q=suya",
		},
		{
			name:          "search_query_and_city",
			queryURL:      "/?q=suya&city=Dallas",
			expectedTitle: "Suya in Dallas | agbalumo",
			expectedDesc:  "Find verified suya in Dallas. Real menus, contact information, and directions in under 60 seconds on agbalumo.",
			expectedURL:   "https://agbalumo.com/?q=suya&city=Dallas",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := testutil.SetupModuleContext(http.MethodGet, tc.queryURL, nil)
			renderer := &metaCaptureRenderer{}
			c.Echo().Renderer = renderer

			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := listing.NewListingHandler(env.App)
			_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

			err := h.HandleHome(c)
			require.NoError(t, err)

			vm, ok := renderer.CapturedData.(listing.HomeViewModel)
			require.True(t, ok, "Captured data must be HomeViewModel")

			assert.Equal(t, tc.expectedTitle, vm.MetaTitle)
			assert.Equal(t, tc.expectedDesc, vm.MetaDescription)
			assert.Equal(t, tc.expectedURL, vm.MetaURL)
		})
	}
}

func TestHandleDetail_DynamicMetaTags(t *testing.T) {
	t.Parallel()

	c, _ := testutil.SetupModuleContext(http.MethodGet, "/listings/detail-meta-1", nil)
	c.SetParamNames("id")
	c.SetParamValues("detail-meta-1")
	renderer := &metaCaptureRenderer{}
	c.Echo().Renderer = renderer

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	testutil.SaveTestListing(t, env.App.DB, "detail-meta-1", "Aria Suya Kitchen", func(l *domain.Listing) {
		l.Type = domain.Food
		l.City = "Arlington"
		l.Description = "Finest smoked suya and jollof in Arlington."
		l.ImageURL = "https://agbalumo.com/static/uploads/aria.jpg"
	})
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	err := h.HandleDetail(c)
	require.NoError(t, err)

	vm, ok := renderer.CapturedData.(listing.DetailViewModel)
	require.True(t, ok, "Captured data must be DetailViewModel")

	assert.Equal(t, "Aria Suya Kitchen | agbalumo", vm.MetaTitle)
	assert.Equal(t, "Finest smoked suya and jollof in Arlington.", vm.MetaDescription)
	assert.Equal(t, "https://agbalumo.com/static/uploads/aria.jpg", vm.MetaImage)
	assert.Equal(t, "https://agbalumo.com/listings/detail-meta-1", vm.MetaURL)
	assert.Equal(t, "restaurant", vm.MetaType)
}
