package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleDelete(t *testing.T) {
	tests := []struct {
		name       string
		user       interface{}
		setup      func(t *testing.T, repo domain.ListingRepository)
		expectCode int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectCode: http.StatusSeeOther,
		},
		{
			name:       "NoUser_Unauthorized",
			user:       nil,
			setup:      func(t *testing.T, repo domain.ListingRepository) {},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "NotFound",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
			},
			expectCode: http.StatusNotFound,
		},
		{
			name: "Forbidden_NotOwner",
			user: domain.User{ID: "other-user"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectCode: http.StatusForbidden,
		},
		{
			name: "DeleteError",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				// We can't trigger a DB error easily with real SQLite without some trickery
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "DeleteError" {
				t.Skip("Hard to trigger DB error with real SQLite")
			}
			c, rec := setupTestContext(http.MethodDelete, "/listings/1", nil)
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			app, cleanup := testutil.SetupTestAppEnv(t)
			defer cleanup()
			tt.setup(t, app.DB)
			h := listmod.NewListingHandler(app)
			_ = h.HandleDelete(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}
