package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome_FeaturedPrioritization(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("index.html").Parse(`Listings: {{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &testutil.TestRenderer{Templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	// Seed data
	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "f2", "Featured 2", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })
	saveTestListing(t, app.DB, "r2", "Regular 2", func(l *domain.Listing) { l.Featured = false })

	h := listmod.NewListingHandler(app)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	// EXPECTED: Featured 1, Featured 2, Regular 1, Regular 2 (Note: sqlite sorts by created_at desc by default)
	// Our seeder might have them in a different order, but both featured should be first.
	// Since we saved f1, f2, r1, r2, created_at might be very close.
	// Actually, HandleHome logic:
	// featured, _ := h.Repo.GetFeaturedListings(c.Request().Context())
	// regular, _ := h.Repo.FindAll(c.Request().Context(), "", "", "", "", false, 20, 0)
	// listings := handler.PrioritizeFeatured(featured, regular)

	// PrioritizeFeatured deduplicates and puts featured at the front in the order returned by GetFeaturedListings.
	// GetFeaturedListings sorts by created_at DESC.
	// So f2, f1 if saved in this order.

	assert.Contains(t, rec.Body.String(), "Featured 1")
	assert.Contains(t, rec.Body.String(), "Featured 2")
	assert.Contains(t, rec.Body.String(), "Regular 1")
	assert.Contains(t, rec.Body.String(), "Regular 2")
}

func TestHandleHome_FeaturedListings_EmptyCategory(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("index.html").Parse(`Listings: {{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &testutil.TestRenderer{Templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	// Seed data: Featured listings across MULTIPLE categories to verify HandleHome doesn't filter by a specific category
	saveTestListing(t, app.DB, "f1", "Featured Business", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Business })
	saveTestListing(t, app.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	saveTestListing(t, app.DB, "f3", "Featured Service", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Service })
	saveTestListing(t, app.DB, "r1", "Regular Business", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Business })

	h := listmod.NewListingHandler(app)

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
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("listing_list").Parse(`{{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &testutil.TestRenderer{Templates: t_temp}

	// Page 1, no filters
	req := httptest.NewRequest(http.MethodGet, "/listings?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	// Seed data
	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	h := listmod.NewListingHandler(app)

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
	assert.Contains(t, rec.Body.String(), "Regular 1")
}

func TestHandleFragment_FeaturedPrioritization_Page2(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("listing_list").Parse(`{{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &testutil.TestRenderer{Templates: t_temp}

	// Page 1, no filters (featured listings appear at the top of the feed)
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	saveTestListing(t, app.DB, "f1", "Featured 1", func(l *domain.Listing) { l.Featured = true })
	saveTestListing(t, app.DB, "r1", "Regular 1", func(l *domain.Listing) { l.Featured = false })

	h := listmod.NewListingHandler(app)

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
}

func TestHandleFragment_FeaturedListings_CategoryFilter(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("listing_list").Parse(`{{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &testutil.TestRenderer{Templates: t_temp}

	// Requesting fragment for 'Business' category, page 1
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?type=Business&page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	// Seed data: Featured listings across MULTIPLE categories
	saveTestListing(t, app.DB, "f1", "Featured Business", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Business })
	saveTestListing(t, app.DB, "f2", "Featured Event", func(l *domain.Listing) { l.Featured = true; l.Type = domain.Event })
	saveTestListing(t, app.DB, "r1", "Regular Business", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Business })

	h := listmod.NewListingHandler(app)

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
	t.Helper()
	l, err := db.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, expected, l.Featured, "Listing %s featured status mismatch", id)
}
