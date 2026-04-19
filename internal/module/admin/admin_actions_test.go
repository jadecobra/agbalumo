package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings", nil)

	// Seed a listing
	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Listing"})

	_ = h.HandleAllListings(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleToggleFeatured(t *testing.T) {
	t.Parallel()
	runToggleTest := func(name, id, featured string, expectCode int, setup func(t *testing.T, repo domain.ListingRepository)) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := admin.NewAdminHandler(env.App)

			formData := url.Values{}
			formData.Set(domain.FieldFeatured, featured)

			urlPath := "/admin/listings/" + id + "/featured"
			if id == "" {
				urlPath = "/admin/listings/featured"
			}
			c, rec := testutil.SetupAdminContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			if id != "" {
				c.SetParamNames("id")
				c.SetParamValues(id)
			}

			setup(t, env.App.DB)

			_ = h.HandleToggleFeatured(c)
			assertFeaturedResponse(t, rec, expectCode, id, env.App.DB)
		})
	}

	runToggleTest("Success", "123", "true", http.StatusOK, func(t *testing.T, repo domain.ListingRepository) {
		_ = repo.Save(context.Background(), domain.Listing{ID: "123", Title: "Test", Featured: false, Type: domain.Business})
	})
	runToggleTest("MissingID", "", "true", http.StatusBadRequest, func(t *testing.T, repo domain.ListingRepository) {})
	runToggleTest("NotFound", "999", "true", http.StatusInternalServerError, func(t *testing.T, repo domain.ListingRepository) {})
	runToggleTest("MaxFeaturedExceeded", "999", "true", http.StatusBadRequest, func(t *testing.T, repo domain.ListingRepository) {
		testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "999", "New", func(l *domain.Listing) { l.Featured = false })
	})
	runToggleTest("ToggleOffWhenMaxReached", "1", "false", http.StatusOK, func(t *testing.T, repo domain.ListingRepository) {
		testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
	})
	runToggleTest("FeatureDifferentCategoryAllowed", "999", "true", http.StatusOK, func(t *testing.T, repo domain.ListingRepository) {
		testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
		testutil.SaveTestListing(t, repo, "999", "New", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Food })
	})
}

func assertFeaturedResponse(t *testing.T, rec *httptest.ResponseRecorder, expectCode int, id string, repo domain.ListingRepository) {
	t.Helper()
	assert.Equal(t, expectCode, rec.Code)

	if expectCode == http.StatusOK {
		// The response must be the HTML row snippet for HTMX swapping.
		// Not a JSON response.
		htmlResponse := rec.Body.String()
		assert.Contains(t, htmlResponse, "listing-row-")
		assert.NotContains(t, htmlResponse, "{\"featured\":")
	}

	if expectCode == http.StatusOK && id == "123" {
		testutil.AssertFeaturedStatus(t, repo, id, true)
	}
}

func TestAdminHandler_HandleApproveClaim(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/claims/cr1/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("cr1")

	// Seed a claim request
	_ = env.App.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "cr1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	cr, _ := env.App.DB.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
	assert.Equal(t, domain.ClaimStatusApproved, cr.Status)
}

func TestAdminHandler_HandleListingRow(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings/1/row", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")

	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Row Listing"})

	_ = h.HandleListingRow(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}
