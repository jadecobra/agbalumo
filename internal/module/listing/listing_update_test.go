package listing_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHandleEdit(t *testing.T) {
	tests := []struct {
		setup          func(t *testing.T, repo domain.ListingRepository)
		user           domain.User
		name           string
		expectedStatus int
	}{
		{
			name: "Success",
			user: newTestUser("owner-1", domain.UserRoleUser),
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: newTestUser("other-user", domain.UserRoleUser),
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodGet, "/listings/1/edit", nil)
			c.SetPath("/listings/:id/edit")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			h, app, cleanup := setupListingHandler(t)
			defer cleanup()
			tt.setup(t, app.DB)

			_ = h.HandleEdit(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
func TestHandleUpdate(t *testing.T) {
	tests := []struct {
		setup          func(t *testing.T, repo domain.ListingRepository)
		user           domain.User
		name           string
		body           string
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&city=Accra&address=123+St",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Old Title", func(l *domain.Listing) { l.OwnerID = "user1" })
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: newTestUser("user2", domain.UserRoleUser),
			body: "",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				saveTestListing(t, repo, "1", "Old Title", func(l *domain.Listing) { l.OwnerID = "user1" })
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(tt.body))
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			h, app, cleanup := setupListingHandler(t)
			defer cleanup()
			tt.setup(t, app.DB)

			_ = h.HandleUpdate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
func TestHandleUpdate_NotFound(t *testing.T) {
	h, _, cleanup := setupListingHandler(t)
	defer cleanup()
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", newTestUser("user1", domain.UserRoleUser))
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestHandleUpdate_NoUser(t *testing.T) {
	h, _, cleanup := setupListingHandler(t)
	defer cleanup()
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestHandleUpdate_DuplicateTitle(t *testing.T) {
	h, app, cleanup := setupListingHandler(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "Old", func(l *domain.Listing) { l.OwnerID = "user1" })
	saveTestListing(t, app.DB, "2", "Taken Title", func(l *domain.Listing) { l.OwnerID = "user2"; l.Address = "456 St" })
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader("title=Taken+Title&type=Business&owner_origin=Ghana&description=Desc&contact_email=t@e.com&city=Kumasi&address=123+St"))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", newTestUser("user1", domain.UserRoleUser))
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
