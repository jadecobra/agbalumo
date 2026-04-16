package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	t.Parallel()
	app, h, c, rec, cleanup := setupAdminDeleteTest(t, "secret", "l1")
	defer cleanup()

	seedListing(t, app, "l1")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify deletion
	_, err := app.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err) // Should not be found
}

func TestHandleAdminDeleteView(t *testing.T) {
	t.Parallel()
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	seedListing(t, app, "listing1")

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete?id=listing1", nil)
	c.Echo().Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_delete_confirm.html")}

	if assert.NoError(t, h.HandleAdminDeleteView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleAdminDeleteView_NoIDs_Redirects(t *testing.T) {
	t.Parallel()
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteView_FindByIDError_Returns404(t *testing.T) {
	t.Parallel()
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()
	// No data seeded, so "bad-id" will not be found

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete?id=bad-id", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleAdminDeleteAction_NoIDs_Redirects(t *testing.T) {
	t.Parallel()
	_, h, c, rec, cleanup := setupAdminDeleteTest(t, "secret")
	defer cleanup()

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteAction_WrongCode_RendersConfirmWithError(t *testing.T) {
	t.Parallel()
	app, h, c, rec, cleanup := setupAdminDeleteTest(t, "wrong", "l1")
	defer cleanup()
	app.Cfg.AdminCode = "correct"

	seedListing(t, app, "l1")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleAdminDeleteAction_PartialSuccess(t *testing.T) {
	t.Parallel()
	app, h, c, rec, cleanup := setupAdminDeleteTest(t, "secret", "l1", "l2")
	defer cleanup()

	seedListing(t, app, "l1")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify l1 deleted
	_, err := app.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err)

	// Verify flash message
	sess := middleware.GetSession(c)
	flashes := sess.Flashes("message")
	assert.Len(t, flashes, 1)
	assert.Contains(t, flashes[0], "Successfully deleted 1 listings")
}

func setupAdminDeleteTest(t *testing.T, adminCode string, ids ...string) (*env.AppEnv, *admin.AdminHandler, echo.Context, *httptest.ResponseRecorder, func()) {
	formData := url.Values{}
	if adminCode != "" {
		formData.Set("admin_code", adminCode)
	}
	for _, id := range ids {
		formData.Add("id", id)
	}
	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	app.Cfg.AdminCode = "secret"
	return app, h, c, rec, cleanup
}

func seedListing(t *testing.T, app *env.AppEnv, id string) {
	err := app.DB.Save(context.Background(), domain.Listing{
		ID:          id,
		Title:       "To Delete",
		OwnerOrigin: "Nigeria",
		Type:        "business",
	})
	assert.NoError(t, err)
}
