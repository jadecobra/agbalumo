package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/labstack/echo/v4"

	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate_GeocodingFallback(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	mockGeocoding := &testutil.MockGeocodingService{}
	mockGeocoding.On("GetCity", context.Background(), "1600 Amphitheatre Parkway, Mountain View, CA").Return("Mountain View", nil)

	env.App.GeocodingSvc = mockGeocoding
	h := listmod.NewListingHandler(env.App)

	// Create context with a user
	body := "title=Google+HQ&type=Business&owner_origin=Nigeria&description=Tech+Giant+HQ&contact_email=info@google.com&address=1600+Amphitheatre+Parkway,+Mountain+View,+CA"
	// NOTE WELL: mapping 'city' is intentionally left empty in the form body to trigger fallback

	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings", strings.NewReader(body))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c.Set("User", domain.User{ID: "test-user-id", Email: "info@google.com"})

	// 2. Execute
	err := h.HandleCreate(c)
	assert.NoError(t, err)

	// 3. Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify that the city was automatically populated in the database
	listings, err := env.App.DB.FindByTitle(context.Background(), "Google HQ")
	assert.NoError(t, err)
	assert.Len(t, listings, 1)
	assert.Equal(t, "Mountain View", listings[0].City, "City should be populated from geocoding fallback")
}
