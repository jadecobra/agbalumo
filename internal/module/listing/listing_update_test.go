package listing_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleEdit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		setup          func(t *testing.T, repo domain.ListingRepository)
		user           domain.User
		name           string
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1", Role: domain.UserRoleUser},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user", Role: domain.UserRoleUser},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "Title", func(l *domain.Listing) { l.OwnerID = "owner-1" })
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/1/edit", nil)
			c.SetPath("/listings/:id/edit")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := listing.NewListingHandler(env.App)
			tt.setup(t, env.App.DB)

			_ = h.HandleEdit(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleEdit_NoUser(t *testing.T) {
	t.Parallel()
	c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/1/edit", nil)
	c.SetPath("/listings/:id/edit")
	c.SetParamNames("id")
	c.SetParamValues("1")
	// no user in context — must return 401

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	_ = h.HandleEdit(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleUpdate(t *testing.T) {
	t.Parallel()
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
				testutil.SaveTestListing(t, repo, "1", "Old Title", func(l *domain.Listing) { l.OwnerID = "user1" })
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "user2", Role: domain.UserRoleUser},
			body: "",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				testutil.SaveTestListing(t, repo, "1", "Old Title", func(l *domain.Listing) { l.OwnerID = "user1" })
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1", strings.NewReader(tt.body))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := listing.NewListingHandler(env.App)
			tt.setup(t, env.App.DB)

			_ = h.HandleUpdate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
func TestHandleUpdate_NotFound(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1", Role: domain.UserRoleUser})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestHandleUpdate_NoUser(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestHandleUpdate_DuplicateTitle(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "1", "Old", func(l *domain.Listing) { l.OwnerID = "user1" })
	testutil.SaveTestListing(t, env.App.DB, "2", "Taken Title", func(l *domain.Listing) { l.OwnerID = "user2"; l.Address = "456 St" })
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/1", strings.NewReader("title=Taken+Title&type=Business&owner_origin=Ghana&description=Desc&contact_email=t@e.com&city=Kumasi&address=123+St"))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1", Role: domain.UserRoleUser})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
