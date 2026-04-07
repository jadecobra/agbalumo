package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	defer cleanup()
	app.Cfg.AdminCode = "secret"

	// Seed data
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify deletion
	_, err := app.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err) // Should not be found
}

func TestHandleAdminDeleteView(t *testing.T) {
	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "listing1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete?id=listing1", nil)
	c.Echo().Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_delete_confirm.html")}

	if assert.NoError(t, h.HandleAdminDeleteView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleAdminDeleteView_NoIDs_Redirects(t *testing.T) {
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteView_FindByIDError_Returns404(t *testing.T) {
	_, h, cleanup := setupAdminTest(t)
	defer cleanup()
	// No data seeded, so "bad-id" will not be found

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete?id=bad-id", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleAdminDeleteAction_NoIDs_Redirects(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "secret")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	app, h, cleanup := setupAdminTest(t)
	defer cleanup()
	app.Cfg.AdminCode = "secret"

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteAction_WrongCode_RendersConfirmWithError(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "wrong")
	formData.Add("id", "l1")
	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	defer cleanup()
	app.Cfg.AdminCode = "correct"

	// Seed so it doesn't fail on something else
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleAdminDeleteAction_PartialSuccess(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	formData.Add("id", "l2") // Does not exist
	app, h, c, rec, cleanup := setupAdminBulkTest(t, http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	defer cleanup()
	app.Cfg.AdminCode = "secret"

	// Seed only l1
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

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
