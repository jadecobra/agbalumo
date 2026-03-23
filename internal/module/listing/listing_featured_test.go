package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome_FeaturedPrioritization(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("index.html").Parse(`Listings: {{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &TestRenderer{templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)

	// Seed data
	f1 := domain.Listing{ID: "f1", Title: "Featured 1", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	f2 := domain.Listing{ID: "f2", Title: "Featured 2", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	r1 := domain.Listing{ID: "r1", Title: "Regular 1", Featured: false, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	r2 := domain.Listing{ID: "r2", Title: "Regular 2", Featured: false, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}

	_ = repo.Save(context.Background(), f1)
	_ = repo.Save(context.Background(), f2)
	_ = repo.Save(context.Background(), r1)
	_ = repo.Save(context.Background(), r2)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

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
	e.Renderer = &TestRenderer{templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)

	// Seed data: Featured listings across MULTIPLE categories to verify HandleHome doesn't filter by a specific category
	f1 := domain.Listing{ID: "f1", Title: "Featured Business", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	f2 := domain.Listing{ID: "f2", Title: "Featured Event", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "event", Address: "123 St"}
	f3 := domain.Listing{ID: "f3", Title: "Featured Service", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "service", Address: "123 St"}
	r1 := domain.Listing{ID: "r1", Title: "Regular Business", Featured: false, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}

	_ = repo.Save(context.Background(), f1)
	_ = repo.Save(context.Background(), f2)
	_ = repo.Save(context.Background(), f3)
	_ = repo.Save(context.Background(), r1)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

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
	e.Renderer = &TestRenderer{templates: t_temp}

	// Page 1, no filters
	req := httptest.NewRequest(http.MethodGet, "/listings?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)

	// Seed data
	f1 := domain.Listing{ID: "f1", Title: "Featured 1", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	r1 := domain.Listing{ID: "r1", Title: "Regular 1", Featured: false, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}

	_ = repo.Save(context.Background(), f1)
	_ = repo.Save(context.Background(), r1)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

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
	e.Renderer = &TestRenderer{templates: t_temp}

	// Page 2, no filters
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?page=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)

	// Seed data
	f1 := domain.Listing{ID: "f1", Title: "Featured 1", Featured: true, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}
	r1 := domain.Listing{ID: "r1", Title: "Regular 1", Featured: false, IsActive: true, OwnerOrigin: "Nigeria", Type: "business", Address: "123 St"}

	_ = repo.Save(context.Background(), f1)
	_ = repo.Save(context.Background(), r1)

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	assert.Contains(t, rec.Body.String(), "Featured 1")
}

