package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := testutil.SetupTestRepository(t)
	// Seed a listing
	_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Listing"})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
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
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "F1", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "2", Title: "F2", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "3", Title: "F3", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "999", Title: "New", Type: domain.Business, Featured: false, IsActive: true})
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name:     "ToggleOffWhenMaxReached",
			id:       "1",
			featured: "false",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "F1", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "2", Title: "F2", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "3", Title: "F3", Type: domain.Business, Featured: true, IsActive: true})
			},
			expectCode: http.StatusOK,
		},
		{
			name:     "FeatureDifferentCategoryAllowed",
			id:       "999",
			featured: "true",
			setupData: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "F1", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "2", Title: "F2", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "3", Title: "F3", Type: domain.Business, Featured: true, IsActive: true})
				_ = repo.Save(context.Background(), domain.Listing{ID: "999", Title: "New", Type: domain.Food, Featured: false, IsActive: true})
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
			if tt.id != "" {
				c.SetParamNames("id")
				c.SetParamValues(tt.id)
			}
			c.Set("User", domain.User{Role: domain.UserRoleAdmin})

			repo := testutil.SetupTestRepository(t)
			tt.setupData(t, repo)

			h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
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
				l, _ := repo.FindByID(context.Background(), tt.id)
				assert.True(t, l.Featured)
			}
		})
	}
}

func TestAdminHandler_HandleApproveClaim(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/claims/cr1/approve", nil)
	c.SetParamNames("id")
	c.SetParamValues("cr1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := testutil.SetupTestRepository(t)
	// Seed a claim request
	_ = repo.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "cr1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	cr, _ := repo.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
	assert.Equal(t, domain.ClaimStatusApproved, cr.Status)
}

func TestAdminHandler_HandleListingRow(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/1/row", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := testutil.SetupTestRepository(t)
	_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Row Listing"})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: repo, FeedbackStore: repo, AnalyticsStore: repo, CategoryStore: repo, UserStore: repo, ListingStore: repo, ClaimRequestStore: repo, CSVService: nil, Cfg: config.LoadConfig()})
	_ = h.HandleListingRow(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}
