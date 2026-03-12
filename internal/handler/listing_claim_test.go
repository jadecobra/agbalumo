package handler_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestListingHandler_HandleClaim(t *testing.T) {
	tests := []struct {
		name       string
		user       *domain.User
		listingID  string
		setup      func(t *testing.T, repo domain.ListingRepository)
		expectCode int
	}{
		{
			name:       "NoUser_RedirectsToLogin",
			user:       nil,
			listingID:  "listing1",
			expectCode: http.StatusFound,
		},
		{
			name:      "Success_ReturnsHTML",
			user:      &domain.User{ID: "u1", Name: "Test", Email: "u1@e.com"},
			listingID: "listing1",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "listing1", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "t@e.com", Address: "123 St"})
				_ = repo.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})
			},
			expectCode: http.StatusOK,
		},
		{
			name:      "NotFound_Returns404",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				// No listing saved
			},
			expectCode: http.StatusNotFound,
		},
		{
			name:      "AlreadyOwned_Returns403",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "listing1", OwnerID: "other", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "t@e.com", Address: "123 St"})
			},
			expectCode: http.StatusForbidden,
		},
		{
			name:      "NotClaimable_Returns403",
			user:      &domain.User{ID: "u1", Name: "Test"},
			listingID: "listing1",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "listing1", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "t@e.com", Address: "123 St"})
				_ = repo.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: false, Active: true})
			},
			expectCode: http.StatusForbidden,
		},
		{
			name:      "DuplicateClaim_Returns409",
			user:      &domain.User{ID: "u1", Name: "Test", Email: "u1@e.com"},
			listingID: "listing1",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "listing1", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "t@e.com", Address: "123 St"})
				_ = repo.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})
				_ = repo.SaveClaimRequest(context.Background(), domain.ClaimRequest{ID: "cr1", ListingID: "listing1", UserID: "u1", Status: domain.ClaimStatusPending, CreatedAt: time.Now()})
			},
			expectCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := handler.SetupTestRepository(t)
			if tt.setup != nil {
				tt.setup(t, repo)
			}

			h := handler.NewListingHandler(repo, nil, &handler.MockGeocodingService{})

			c, rec := setupTestContext(http.MethodPost, "/listings/listing1/claim", nil)
			c.SetParamNames("id")
			c.SetParamValues(tt.listingID)

			if tt.user != nil {
				c.Set("User", *tt.user)
			}

			_ = h.HandleClaim(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}
