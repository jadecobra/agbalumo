package listing_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHandleDelete(t *testing.T) {
	tests := []struct {
		user       interface{}
		setup      func(t *testing.T, repo domain.ListingRepository)
		name       string
		expectCode int
	}{
		{
			name: "Success",
			user: newTestUser("owner-1", domain.UserRoleUser),
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
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
			user: newTestUser("owner-1", domain.UserRoleUser),
			setup: func(t *testing.T, repo domain.ListingRepository) {
			},
			expectCode: http.StatusNotFound,
		},
		{
			name: "Forbidden_NotOwner",
			user: newTestUser("other-user", domain.UserRoleUser),
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
			},
			expectCode: http.StatusForbidden,
		},
		{
			name: "DeleteError",
			user: newTestUser("owner-1", domain.UserRoleUser),
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

			h, app, cleanup := setupListingHandler(t)
			defer cleanup()
			tt.setup(t, app.DB)
			_ = h.HandleDelete(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}
