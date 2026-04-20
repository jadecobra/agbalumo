package listing_test

import (
	"fmt"
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
			name: "Test Title",
			body: fmt.Sprintf("%s=Test+Title&%s=%s&%s=Nigeria&%s=Cool&%s=test@example.com&%s=Lagos&%s=123+Street",
				domain.FieldTitle, domain.FieldType, domain.Business, domain.FieldOwnerOrigin, domain.FieldDescription, domain.FieldContactEmail, domain.FieldCity, domain.FieldAddress),
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
		{
			name: "ValidationError",
			body: fmt.Sprintf("%s=Test+Title&%s=%s", domain.FieldTitle, domain.FieldType, domain.Business),
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Page",
		},
		{
			name: "Req",
			body: fmt.Sprintf("%s=Req&%s=%s&%s=Nigeria&%s=Cool&%s=test@example.com&%s=Lagos&%s=123+St",
				domain.FieldTitle, domain.FieldType, domain.Request, domain.FieldOwnerOrigin, domain.FieldDescription, domain.FieldContactEmail, domain.FieldCity, domain.FieldAddress),
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Biz",
			body: fmt.Sprintf("%s=Biz&%s=%s&%s=Nigeria&%s=Cool&%s=test@example.com&%s=Lagos&%s=123+Street&%s=example.com",
				domain.FieldTitle, domain.FieldType, domain.Business, domain.FieldOwnerOrigin, domain.FieldDescription, domain.FieldContactEmail, domain.FieldCity, domain.FieldAddress, domain.FieldWebsiteURL),
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := testutil.SetupTestModuleEnv(t)
			defer env.Cleanup()

			c, rec := testutil.SetupModuleContext(http.MethodPost, domain.PathListings, strings.NewReader(tt.body))
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c.Set(domain.CtxKeyUser, &domain.User{ID: "test-user-id", Role: domain.UserRoleUser})

			h := listing.NewListingHandler(env.App)
			tt.setup(t, env.App.DB)

			_ = h.HandleCreate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				testutil.AssertListingExists(t, env.App.DB, tt.name) // Use tt.name as title to match test setup
			}
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
