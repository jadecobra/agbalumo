package listing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleDetail_Routing(t *testing.T) {
	t.Parallel()

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := listing.NewListingHandler(env.App)

	testutil.SaveTestListing(t, env.App.DB, "test-detail-route", "Lagos Kitchen", func(l *domain.Listing) {
		l.Type = domain.Food
		l.City = "Dallas"
		l.Description = "Delicious Yoruba delicacies."
	})
	_ = env.App.DB.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Food), Name: "Food", Active: true})

	t.Run("htmx_request_returns_modal_partial", func(t *testing.T) {
		c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/test-detail-route", nil)
		c.SetParamNames("id")
		c.SetParamValues("test-detail-route")
		c.Request().Header.Set("HX-Request", "true")

		err := h.HandleDetail(c)
		require.NoError(t, err)

		body := rec.Body.String()
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, body, "detail-modal-test-detail-route")
		assert.Contains(t, body, "Lagos Kitchen")
		assert.NotContains(t, body, "<!DOCTYPE html>")
		assert.NotContains(t, body, "<html")
	})

	t.Run("direct_visit_returns_full_html_page_with_meta_and_close_redirect", func(t *testing.T) {
		c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/test-detail-route", nil)
		c.SetParamNames("id")
		c.SetParamValues("test-detail-route")

		err := h.HandleDetail(c)
		require.NoError(t, err)

		body := rec.Body.String()
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, body, "<!DOCTYPE html>")
		assert.Contains(t, body, "<html")
		assert.Contains(t, body, "<title>Lagos Kitchen | agbalumo</title>")
		assert.Contains(t, body, "detail-modal-test-detail-route")
		assert.Contains(t, body, `data-close-redirect="/"`)
	})
}
