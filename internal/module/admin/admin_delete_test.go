package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.AdminCode = "secret"

	// Seed data
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	h := admin.NewAdminHandler(app)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify deletion
	_, err := app.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err) // Should not be found
}

func TestHandleAdminDeleteView(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "listing1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})
	h := admin.NewAdminHandler(app)

	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete?id=listing1", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplateForPage(t, "admin_delete_confirm.html")}
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.HandleAdminDeleteView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleAdminDeleteView_NoIDs_Redirects(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := admin.NewAdminHandler(app)

	c, rec := setupAdminTestContext(http.MethodGet, "/admin/listings/delete", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteView_FindByIDError_Returns404(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	// No data seeded, so "bad-id" will not be found
	h := admin.NewAdminHandler(app)

	req := httptest.NewRequest(http.MethodGet, "/admin/listings/delete?id=bad-id", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Renderer = &AdminMockRenderer{}
	c := e.NewContext(req, rec)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleAdminDeleteAction_NoIDs_Redirects(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "secret")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.AdminCode = "secret"
	h := admin.NewAdminHandler(app)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteAction_WrongCode_RendersConfirmWithError(t *testing.T) {
	formData := url.Values{}
	formData.Set("admin_code", "wrong")
	formData.Add("id", "l1")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.AdminCode = "correct"

	// Seed so it doesn't fail on something else
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	h := admin.NewAdminHandler(app)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleAdminDeleteAction_PartialSuccess(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.AdminCode = "secret"

	// Seed only l1
	_ = app.DB.Save(context.Background(), domain.Listing{ID: "l1", Title: "To Delete", OwnerOrigin: "Nigeria", Type: "business"})

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	formData.Add("id", "l2") // Does not exist
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))

	h := admin.NewAdminHandler(app)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

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
