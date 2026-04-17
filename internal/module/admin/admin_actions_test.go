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
	tests := []struct {
		setupData  func(t *testing.T, repo domain.ListingRepository)
		name       string
		id         string
		featured   string
		expectCode int
	}{
		{
			name:     "Success",
			id:       "123",
			featured: "true",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "123", Title: "Test", Featured: false, Type: domain.Business})
			},
			expectCode: http.StatusOK,
		},
		{
			name:       "MissingID",
			id:         "",
			featured:   "true",
			setupData:  func(t *testing.T, repo domain.ListingRepository) {},
			expectCode: http.StatusBadRequest,
		},
		{
			name:     "NotFound",
			id:       "999",
			featured: "true",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
			},
			expectCode: http.StatusInternalServerError, // FindByID fails
		},
		{
			name:     "MaxFeaturedExceeded",
			id:       "999",
			featured: "true",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "999", "New", func(l *domain.Listing) { l.Featured = false })
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name:     "ToggleOffWhenMaxReached",
			id:       "1",
			featured: "false",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
			},
			expectCode: http.StatusOK,
		},
		{
			name:     "FeatureDifferentCategoryAllowed",
			id:       "999",
			featured: "true",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "F1", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "2", "F2", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "3", "F3", func(l *domain.Listing) { l.Featured = true })
				testutil.SaveTestListing(t, repo, "999", "New", func(l *domain.Listing) { l.Featured = false; l.Type = domain.Food })
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := admin.NewAdminHandler(env.App)

			formData := url.Values{}
			formData.Set("featured", tt.featured)
			urlPath := "/admin/listings/" + tt.id + "/featured"
			if tt.id == "" {
				urlPath = "/admin/listings/featured"
			}
			c, rec := testutil.SetupAdminContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			if tt.id != "" {
				c.SetParamNames("id")
				c.SetParamValues(tt.id)
			}

			tt.setupData(t, env.App.DB)

			_ = h.HandleToggleFeatured(c)
			assertFeaturedResponse(t, rec, tt.expectCode, tt.id, env.App.DB)
		})
	}
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
