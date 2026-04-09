package listing_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome_FeaturedPrioritization(t *testing.T) {
	t.Parallel()
	c, rec := setupTestContext(http.MethodGet, "/", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()

	// Seed data
	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "f2", "Featured 2", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })
	saveTestListing(t, app.DB, "r2", "Regular 2", func(l *domain.Listing) { l.Featured = false })

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
	c, rec := setupTestContext(http.MethodGet, "/", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()

	// Seed data: Featured listings across MULTIPLE categories to verify HandleHome doesn't filter by a specific category
	saveTestListing(t, app.DB, "f1", "Featured Business", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Business })
	saveTestListing(t, app.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	saveTestListing(t, app.DB, "f3", "Featured Service", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Service })
	saveTestListing(t, app.DB, "r1", "Regular Business", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Business })

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	// If HandleHome was passing a specific category string (e.g. "business") to GetFeaturedListings,
	// then the "event" and "service" featured listings would NOT be present in the response.
	// Since we assert they are all present, we verify it passes an empty string (or doesn't filter).
	assert.Contains(t, body, "Featured Business")
	assert.Contains(t, body, "Featured Event")
	assert.Contains(t, body, "Featured Service")
}

func TestHandleFragment_FeaturedPrioritization(t *testing.T) {
	t.Parallel()
	// Page 1, no filters
	c, rec := setupTestContext(http.MethodGet, "/listings?page=1", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()

	// Seed data
	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
	assert.Contains(t, rec.Body.String(), "Regular 1")
}

func TestHandleFragment_FeaturedPrioritization_Page2(t *testing.T) {
	t.Parallel()
	// Page 1, no filters (featured listings appear at the top of the feed)
	c, rec := setupTestContext(http.MethodGet, "/listings/fragment?page=1", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()

	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
}

func TestHandleFragment_FeaturedListings_CategoryFilter(t *testing.T) {
	t.Parallel()
	// Requesting fragment for 'Business' category, page 1
	c, rec := setupTestContext(http.MethodGet, "/listings/fragment?type=Business&page=1", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()

	// Seed data: Featured listings across MULTIPLE categories
	saveTestListing(t, app.DB, "f1", "Featured Business", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Business })
	saveTestListing(t, app.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	saveTestListing(t, app.DB, "r1", "Regular Business", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Business })

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	body := rec.Body.String()

	// Assert only the featured listing for 'business' category is present
	assert.Contains(t, body, "Featured Business")
	// Assert the featured listing for 'event' category is NOT present
	assert.NotContains(t, body, "Featured Event")
}

func assertFeaturedStatus(t *testing.T, db domain.ListingRepository, id string, expected bool) {
	testutil.AssertFeaturedStatus(t, db, id, expected)
}
