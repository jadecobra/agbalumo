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

func TestHandleCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		setup          func(t *testing.T, repo domain.ListingRepository)
		name           string
		body           string
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "Success",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&hours_of_operation=Mon-Fri+9-5&city=Lagos&address=123+Street",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				// No extra setup needed
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ValidationError",
			body:           "title=Test+Title&type=Business",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Page",
		},
		{
			name:           "RequestWithoutDeadline",
			body:           "title=Req&type=Request&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+St",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "BusinessWithNoPrefixURL",
			body:           "title=Biz&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+Street&website_url=example.com",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()
			h := listing.NewListingHandler(env.App)
			tt.setup(t, env.App.DB)
			c.Set("User", domain.User{ID: "test-user-id", Role: domain.UserRoleUser})

			_ = h.HandleCreate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
func TestHandleCreate_NoUser(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings", nil)
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestHandleCreate_DuplicateTitle(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "x", "Existing")
	body := "title=Existing&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+Street"
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings", strings.NewReader(body))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c.Set("User", domain.User{ID: "test-user-id", Role: domain.UserRoleUser})
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
