package admin_test

import (
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAllListings(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := handler.SetupTestRepository(t)
	// Seed a listing
	_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "Test Listing"})

	h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
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
				_ = repo.Save(context.Background(), domain.Listing{ID: "123", Title: "Test", Featured: false})
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
			expectCode: http.StatusOK, // SQLite UPDATE is no-op, repo returns nil error
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

			repo := handler.SetupTestRepository(t)
			tt.setupData(t, repo)

			h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
			_ = h.HandleToggleFeatured(c)
			assert.Equal(t, tt.expectCode, rec.Code)

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

	repo := handler.SetupTestRepository(t)
	// Seed a claim request
	_ = repo.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "cr1", UserID: "u1", ListingID: "l1", Status: domain.ClaimStatusPending})

	h := admin.NewAdminHandler(repo, repo, repo, repo, repo, repo, repo, nil, config.LoadConfig())
	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	cr, _ := repo.GetClaimRequestByUserAndListing(context.Background(), "u1", "l1")
	assert.Equal(t, domain.ClaimStatusApproved, cr.Status)
}
