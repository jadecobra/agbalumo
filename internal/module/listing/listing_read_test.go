package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	repo := handler.SetupTestRepository(t)
	_ = repo.Save(ctx, domain.Listing{
		ID:           "1",
		Title:        "Listing 1",
		Type:         domain.Business,
		IsActive:     true,
		Status:       domain.ListingStatusApproved,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
	})
	_ = repo.SaveCategory(ctx, domain.CategoryData{ID: string(domain.Business), Name: "Business", Active: true})

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
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Listing 1")
}

func TestHandleDetail(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/listings/1", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	ctx := context.Background()

	repo := handler.SetupTestRepository(t)
	_ = repo.Save(ctx, domain.Listing{
		ID:           "1",
		Title:        "Detail View",
		Type:         domain.Business,
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
	})

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
	if err := h.HandleDetail(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "Detail View")
}

func TestHandleProfile(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/profile", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	user := domain.User{ID: "u1", Name: "John Doe"}
	c.Set("User", user)

	repo := handler.SetupTestRepository(t)
	_ = repo.Save(ctx, domain.Listing{
		ID:           "1",
		Title:        "My Listing",
		OwnerID:      "u1",
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
		Type:         domain.Business,
	})

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
	if err := h.HandleProfile(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "John Doe")
}



func TestHandleFragment(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/listings/fragment?q=Search", nil)
	req.Header.Set("HX-Request", "true")
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	repo := handler.SetupTestRepository(t)
	// Seed 31 listings to test pagination limit of 30
	for i := 1; i <= 31; i++ {
		_ = repo.Save(ctx, domain.Listing{
			ID:           strconv.Itoa(i),
			Title:        "Search Result " + strconv.Itoa(i),
			Type:         domain.Business,
			Status:       domain.ListingStatusApproved,
			IsActive:     true,
			Address:      "Lagos",
			ContactEmail: "test@example.com",
			OwnerOrigin:  "Nigeria",
		})
	}

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
		t.Fatal(err)
	}

	// Verify fragment contains results
	assert.Contains(t, rec.Body.String(), "Search Result 1")
	// Verify it contains the OOB swap for pagination
	assert.Contains(t, rec.Body.String(), `hx-swap-oob="true"`)
	assert.Contains(t, rec.Body.String(), `id="pagination-controls"`)
}
