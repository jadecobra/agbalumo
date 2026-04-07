package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings", nil)
	setupAdminAuth(t, c)
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	// Seed a listing
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Listing"})

	_ = h.HandleAllListings(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_HandleToggleFeatured(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		featured   string
		setupData  func(t *testing.T, repo domain.ListingRepository)
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
			formData := url.Values{}
			formData.Set("featured", tt.featured)
			urlPath := "/admin/listings/" + tt.id + "/featured"
			if tt.id == "" {
				urlPath = "/admin/listings/featured"
			}
			c, rec := setupAdminTestContext(http.MethodPost, urlPath, strings.NewReader(formData.Encode()))
			setupAdminAuth(t, c)
			if tt.id != "" {
				c.SetParamNames("id")
				c.SetParamValues(tt.id)
			}

			app, h, cleanup := setupAdminTest(t)
			defer cleanup()
			tt.setupData(t, app.DB)

			_ = h.HandleToggleFeatured(c)
			assert.Equal(t, tt.expectCode, rec.Code)

			if tt.expectCode == http.StatusOK {
				// The response must be the HTML row snippet for HTMX swapping.
				// Not a JSON response.
				htmlResponse := rec.Body.String()
				assert.Contains(t, htmlResponse, "listing-row-")
				assert.NotContains(t, htmlResponse, "{\"featured\":")
			}

			if tt.expectCode == http.StatusOK && tt.id == "123" {
				testutil.AssertFeaturedStatus(t, app.DB, tt.id, true)
			}
		})
	}
}

func TestAdminHandler_HandleApproveClaim(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/cr1/approve", nil)
	setupAdminAuth(t, c)
	c.SetParamNames("id")
	c.SetParamValues("cr1")

	app, h, cleanup := setupAdminTest(t)
	defer cleanup()

	// Seed a claim request
	_ = app.DB.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "cr1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	cr, _ := app.DB.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
	assert.Equal(t, domain.ClaimStatusApproved, cr.Status)
}

func TestAdminHandler_HandleListingRow(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/1/row", nil)
	setupAdminAuth(t, c)
	c.SetParamNames("id")
	c.SetParamValues("1")

	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Row Listing"})

	_ = h.HandleListingRow(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}
