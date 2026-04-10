package listing_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome(t *testing.T) {
	t.Parallel()
	c, rec := setupTestContext(http.MethodGet, "/", nil)
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "Listing 1", func(l *domain.Listing) { l.Type = domain.Food })
	_ = app.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Listing 1")
}

func TestHandleDetail(t *testing.T) {
	t.Parallel()
	c, rec := setupTestContext(http.MethodGet, "/listings/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")

	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "Detail View")

	if err := h.HandleDetail(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "Detail View")
}

func TestHandleProfile(t *testing.T) {
	t.Parallel()
	c, rec := setupTestContext(http.MethodGet, "/profile", nil)
	user := newTestUser("u1", domain.UserRoleUser)
	user.Name = "John Doe"
	c.Set("User", user)

	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "My Listing", func(l *domain.Listing) { l.OwnerID = "u1" })

	if err := h.HandleProfile(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "John Doe")
}

func TestHandleFragment(t *testing.T) {
	t.Parallel()
	c, rec := setupTestContext(http.MethodGet, "/listings/fragment?q=Search", nil)
	c.Request().Header.Set("HX-Request", "true")

	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	for i := 1; i <= 31; i++ {
		saveTestListing(t, app.DB, strconv.Itoa(i), "Search Result "+strconv.Itoa(i))
	}

	if err := h.HandleFragment(c); err != nil {
		t.Fatal(err)
	}

	// Verify fragment contains results
	assert.Contains(t, rec.Body.String(), "Search Result 1")
	// Verify it contains the OOB swap for pagination
	assertContainsPagination(t, rec.Body.String())
}
