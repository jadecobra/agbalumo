package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setup          func(t *testing.T, repo domain.ListingRepository)
		expectedStatus int
		expectedBody   string
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
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			app, cleanup := testutil.SetupTestAppEnv(t)
			defer cleanup()
			tt.setup(t, app.DB)
			h := listmod.NewListingHandler(app)
			c.Set("User", domain.User{ID: "test-user-id"})

			_ = h.HandleCreate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
func TestHandleCreate_NoUser(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	c, rec := setupTestContext(http.MethodPost, "/listings", nil)
	h := listmod.NewListingHandler(app)
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestHandleCreate_DuplicateTitle(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "x", Title: "Existing", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
	body := "title=Existing&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+Street"
	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))
	c.Set("User", domain.User{ID: "test-user-id"})
	h := listmod.NewListingHandler(app)
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
