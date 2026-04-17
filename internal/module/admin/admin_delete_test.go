package admin_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleAdminDeleteAction_Success(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)
	env.App.Cfg.AdminCode = "secret"

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	testutil.SaveTestListing(t, env.App.DB, "l1", "To Delete")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify deletion
	_, err := env.App.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err) // Should not be found
}

func TestHandleAdminDeleteView(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)
	testutil.SaveTestListing(t, env.App.DB, "listing1", "To Delete")

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings/delete?id=listing1", nil)
	c.Echo().Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_delete_confirm.html")}

	if assert.NoError(t, h.HandleAdminDeleteView(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleAdminDeleteView_NoIDs_Redirects(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings/delete", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteView_FindByIDError_Returns404(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)
	// No data seeded, so "bad-id" will not be found

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings/delete?id=bad-id", nil)

	_ = h.HandleAdminDeleteView(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleAdminDeleteAction_NoIDs_Redirects(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/delete", nil)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin/listings", rec.Header().Get("Location"))
}

func TestHandleAdminDeleteAction_WrongCode_RendersConfirmWithError(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.Cfg.AdminCode = "correct"
	h := admin.NewAdminHandler(env.App)

	formData := url.Values{}
	formData.Set("admin_code", "wrong")
	formData.Add("id", "l1")
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	testutil.SaveTestListing(t, env.App.DB, "l1", "To Delete")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleAdminDeleteAction_PartialSuccess(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.App.Cfg.AdminCode = "secret"
	h := admin.NewAdminHandler(env.App)

	formData := url.Values{}
	formData.Set("admin_code", "secret")
	formData.Add("id", "l1")
	formData.Add("id", "l2")
	c, rec := testutil.SetupAdminContext(http.MethodPost, "/admin/listings/delete", strings.NewReader(formData.Encode()))
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	testutil.SaveTestListing(t, env.App.DB, "l1", "To Delete")

	_ = h.HandleAdminDeleteAction(c)
	assert.Equal(t, http.StatusFound, rec.Code)

	// Verify l1 deleted
	_, err := env.App.DB.FindByID(context.Background(), "l1")
	assert.Error(t, err)

	// Verify flash message
	sess := middleware.GetSession(c)
	flashes := sess.Flashes(domain.FlashMessageKey)
	assert.Len(t, flashes, 1)
	assert.Contains(t, flashes[0], "Successfully deleted 1 listings")
}
