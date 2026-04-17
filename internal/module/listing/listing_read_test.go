package listing_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/", nil)
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Listing 1", func(l *domain.Listing) { l.Type = domain.Food })
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Listing 1")
}

func TestHandleDetail(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Detail View")

	if err := h.HandleDetail(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "Detail View")
}

func TestHandleProfile(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/profile", nil)
	user := domain.User{ID: "u1", Name: "John Doe", Role: domain.UserRoleUser}
	c.Set("User", user)

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "My Listing", func(l *domain.Listing) { l.OwnerID = "u1" })

	if err := h.HandleProfile(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "John Doe")
}

func TestHandleFragment(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?q=Search", nil)
	c.Request().Header.Set("HX-Request", "true")

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	for i := 1; i <= 31; i++ {
		testutil.SaveTestListing(t, env.App.DB, strconv.Itoa(i), "Search Result "+strconv.Itoa(i))
	}

	if err := h.HandleFragment(c); err != nil {
		t.Fatal(err)
	}

	// Verify fragment contains results
	assert.Contains(t, rec.Body.String(), "Search Result 1")
	// Verify it contains the OOB swap for pagination
	testutil.AssertContainsPagination(t, rec.Body.String())
}

func TestHandleDetail_NotFound(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/nonexistent", nil)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	_ = h.HandleDetail(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleFragment_AdaDefaulting(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	// Seed: One food in Houston, one service in Houston
	testutil.SaveTestListing(t, env.App.DB, "1", "Houston Jollof", func(l *domain.Listing) {
		l.Type = domain.Food
		l.City = "Houston"
		l.Status = domain.ListingStatusApproved
		l.IsActive = true
	})
	testutil.SaveTestListing(t, env.App.DB, "2", "Houston Hair", func(l *domain.Listing) {
		l.Type = domain.Service
		l.City = "Houston"
		l.Status = domain.ListingStatusApproved
		l.IsActive = true
	})

	// Case 1: Filter by City only (no type) -> should default to Food
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?city=Houston", nil)
	c.Request().Header.Set("HX-Request", "true")
	if err := h.HandleFragment(c); err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, rec.Body.String(), "Houston Jollof")
	assert.NotContains(t, rec.Body.String(), "Houston Hair")

	// Case 2: Filter by City AND specify Type -> should respect Type
	c2, rec2 := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?city=Houston&type=Service", nil)
	c2.Request().Header.Set("HX-Request", "true")
	if err := h.HandleFragment(c2); err != nil {
		t.Fatal(err)
	}
	assert.NotContains(t, rec2.Body.String(), "Houston Jollof")
	assert.Contains(t, rec2.Body.String(), "Houston Hair")
}
