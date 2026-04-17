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

func TestHandleHome_FeaturedPrioritization(t *testing.T) {
	t.Parallel()
	c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/")
	defer env.Cleanup()

	// Seed data (Defaulting to Food for Ada)
	testutil.SaveTestListing(t, env.App.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Food })
	testutil.SaveTestListing(t, env.App.DB, "f2", "Featured 2", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Food })
	testutil.SaveTestListing(t, env.App.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Food })
	testutil.SaveTestListing(t, env.App.DB, "r2", "Regular 2", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Food })

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
	assert.Contains(t, rec.Body.String(), "Featured 2")
	assert.Contains(t, rec.Body.String(), "Regular 1")
	assert.Contains(t, rec.Body.String(), "Regular 2")
}

func TestHandleHome_FeaturedListings_EmptyCategory(t *testing.T) {
	t.Parallel()
	c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/")
	defer env.Cleanup()

	// Seed data: Featured listings across MULTIPLE categories (One Food, one Business, one Event)
	testutil.SaveTestListing(t, env.App.DB, "f1", "Featured Food", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Food })
	testutil.SaveTestListing(t, env.App.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	testutil.SaveTestListing(t, env.App.DB, "r1", "Regular Food", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Food })

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	// Since HandleHome now defaults to Food, only "Featured Food" should appear
	assert.Contains(t, body, "Featured Food")
	assert.NotContains(t, body, "Featured Event")
}

func TestHandleFragment_FeaturedPrioritization(t *testing.T) {
	t.Parallel()
	c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/listings?page=1")
	defer env.Cleanup()

	// Seed data
	testutil.SaveTestListing(t, env.App.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	testutil.SaveTestListing(t, env.App.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
	assert.Contains(t, rec.Body.String(), "Regular 1")
}

func TestHandleFragment_FeaturedPrioritization_Page2(t *testing.T) {
	t.Parallel()
	c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/listings/fragment?page=1")
	defer env.Cleanup()

	testutil.SaveTestListing(t, env.App.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	testutil.SaveTestListing(t, env.App.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
}

func TestHandleFragment_FeaturedListings_CategoryFilter(t *testing.T) {
	t.Parallel()
	// Requesting fragment for 'Business' category, page 1
	c, rec, env, h := setupFeaturedTest(t, http.MethodGet, "/listings/fragment?type=Business&page=1")
	defer env.Cleanup()

	// Seed data: Featured listings across MULTIPLE categories
	testutil.SaveTestListing(t, env.App.DB, "f1", "Featured Business", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Business })
	testutil.SaveTestListing(t, env.App.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	testutil.SaveTestListing(t, env.App.DB, "r1", "Regular Business", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Business })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	body := rec.Body.String()

	// Assert only the featured listing for 'business' category is present
	assert.Contains(t, body, "Featured Business")
	// Assert the featured listing for 'event' category is NOT present
	assert.NotContains(t, body, "Featured Event")
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
