package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome_Featured(t *testing.T) {
	t.Parallel()
	tests := []struct {
		seed       func(t *testing.T, db domain.ListingRepository)
		assertions func(t *testing.T, body string)
		name       string
	}{
		{
			name: "Prioritization",
			seed: func(t *testing.T, db domain.ListingRepository) {
				for _, s := range []struct {
					id       string
					title    string
					cat      domain.Category
					featured bool
				}{
					{"f1", "Featured 1", domain.Food, true},
					{"f2", "Featured 2", domain.Food, true},
					{"r1", "Regular 1", domain.Food, false},
				} {
					testutil.SaveTestListing(t, db, s.id, s.title, func(l *domain.Listing) {
						l.Type = s.cat
						l.Featured = s.featured
					})
				}
			},
			assertions: func(t *testing.T, body string) {
				assert.Contains(t, body, "Featured 1")
				assert.Contains(t, body, "Featured 2")
				assert.Contains(t, body, "Regular 1")
			},
		},
		{
			name: "EmptyCategory_DefaultsToFood",
			seed: func(t *testing.T, db domain.ListingRepository) {
				for _, s := range []struct {
					cat   domain.Category
					id    string
					title string
				}{
					{cat: domain.Food, id: "f1", title: "Featured Food"},
					{cat: domain.Event, id: "f2", title: "Featured Event"},
				} {
					testutil.SaveTestListing(t, db, s.id, s.title, func(l *domain.Listing) {
						l.Featured = true
						l.Type = s.cat
					})
				}
			},
			assertions: func(t *testing.T, body string) {
				assert.Contains(t, body, "Featured Food")
				assert.NotContains(t, body, "Featured Event")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/")
			defer env.Cleanup()
			tt.seed(t, env.App.DB)
			if err := h.HandleHome(c); err != nil {
				t.Fatalf("HandleHome failed: %v", err)
			}
			tt.assertions(t, rec.Body.String())
		})
	}
}

func TestHandleFragment_Featured(t *testing.T) {
	t.Parallel()
	tests := []struct {
		seed       func(t *testing.T, db domain.ListingRepository)
		assertions func(t *testing.T, body string)
		name       string
		target     string
	}{
		{
			name:   "BasicPrioritization",
			target: "/listings/fragment?page=1&type=All",
			seed: func(t *testing.T, db domain.ListingRepository) {
				for _, s := range []struct {
					id       string
					title    string
					featured bool
				}{
					{"f1", "Featured 1", true},
					{"r1", "Regular 1", false},
				} {
					testutil.SaveTestListing(t, db, s.id, s.title, func(l *domain.Listing) {
						l.Featured = s.featured
					})
				}
			},
			assertions: func(t *testing.T, body string) {
				assert.Contains(t, body, "Featured 1")
				assert.Contains(t, body, "Regular 1")
				assert.Contains(t, body, `id="featured-section" hx-swap-oob="true"`)
			},
		},
		{
			name:   "CategoryFilter",
			target: "/listings/fragment?type=Business&page=1",
			seed: func(t *testing.T, db domain.ListingRepository) {
				for _, s := range []struct {
					id    string
					title string
					cat   domain.Category
				}{
					{"f1", "Featured Business", domain.Business},
					{"f2", "Featured Event", domain.Event},
				} {
					testutil.SaveTestListing(t, db, s.id, s.title, func(l *domain.Listing) {
						l.Featured = true
						l.Type = s.cat
					})
				}
			},
			assertions: func(t *testing.T, body string) {
				assert.Contains(t, body, "Featured Business")
				assert.NotContains(t, body, "Featured Event")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec, env, h := setupFeaturedTest(t, http.MethodGet, tt.target)
			defer env.Cleanup()
			tt.seed(t, env.App.DB)
			if err := h.HandleFragment(c); err != nil {
				t.Fatalf("HandleFragment failed: %v", err)
			}
			tt.assertions(t, rec.Body.String())
		})
	}
}

func setupFeaturedTest(t *testing.T, method, target string) (echo.Context, *httptest.ResponseRecorder, testutil.ModuleTestEnv, *listing.ListingHandler) {
	c, rec := testutil.SetupModuleContext(method, target, nil)
	env := testutil.SetupTestModuleEnv(t)
	h := listing.NewListingHandler(env.App)
	return c, rec, env, h
}

func assertFeaturedStatus(t *testing.T, db domain.ListingRepository, id string, expected bool) {
	testutil.AssertFeaturedStatus(t, db, id, expected)
}
