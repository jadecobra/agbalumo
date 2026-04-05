package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{
		ID:           "1",
		Title:        "Listing 1",
		Type:         domain.Business,
		IsActive:     true,
		Status:       domain.ListingStatusApproved,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
	})
	_ = app.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: "Business", Active: true})
	h := listmod.NewListingHandler(app)
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

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{
		ID:           "1",
		Title:        "Detail View",
		Type:         domain.Business,
		Status:       domain.ListingStatusApproved,
		IsActive:     true,
		Address:      "Lagos",
		ContactEmail: "test@example.com",
		OwnerOrigin:  "Nigeria",
	})
	h := listmod.NewListingHandler(app)
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

	user := domain.User{ID: "u1", Name: "John Doe"}
	c.Set("User", user)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{
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
	h := listmod.NewListingHandler(app)
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

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	for i := 1; i <= 31; i++ {
		_ = app.DB.Save(context.Background(), domain.Listing{
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
	h := listmod.NewListingHandler(app)
	if err := h.HandleFragment(c); err != nil {
		t.Fatal(err)
	}

	// Verify fragment contains results
	assert.Contains(t, rec.Body.String(), "Search Result 1")
	// Verify it contains the OOB swap for pagination
	assert.Contains(t, rec.Body.String(), `hx-swap-oob="true"`)
	assert.Contains(t, rec.Body.String(), `id="pagination-controls"`)
}
