package admin_test

import (
	"html/template"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
)

// AdminMockRenderer is a simple renderer for admin tests
type AdminMockRenderer struct{}

func (m *AdminMockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if name == "admin_listing_table_row" {
		listing, ok := data.(domain.Listing)
		if ok {
			_, err := w.Write([]byte(`<tr id="listing-row-` + listing.ID + `">Mock HTML Row</tr>`))
			return err
		}
	}
	// Simple mock that does nothing but satisfying the interface
	return nil
}

// setupAdminTestContext sets up a basic Echo context for admin tests
func setupAdminTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &AdminMockRenderer{}
	req := httptest.NewRequest(method, target, body)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func setupAdminAuth(t *testing.T, c echo.Context) {
	c.Set("User", domain.User{ID: "admin1", Role: domain.UserRoleAdmin})
}

func setupAdminTest(t *testing.T) (*env.AppEnv, *admin.AdminHandler, func()) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	h := admin.NewAdminHandler(app)
	return app, h, cleanup
}


func setupAdminBulkTest(t *testing.T, method, target string, body io.Reader) (*env.AppEnv, *admin.AdminHandler, echo.Context, *httptest.ResponseRecorder, func()) {
	c, rec := setupAdminTestContext(method, target, body)
	setupAdminAuth(t, c)
	app, h, cleanup := setupAdminTest(t)
	app.CSVService = service.NewCSVService()

	// Set session
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	return app, h, c, rec, cleanup
}

// setupAdminAuthWithID sets a user with a specific ID as admin
func setupAdminAuthWithID(c echo.Context, userID string) {
	c.Set("User", domain.User{ID: userID, Role: domain.UserRoleAdmin})
}

// NewAdminTemplate returns a mock template for admin tests
func NewAdminTemplate() *template.Template {
	return template.Must(template.New("admin").Parse(`
		{{define "admin_login.html"}}Login{{end}}
		{{define "admin_dashboard.html"}}Dashboard{{end}}
		{{define "admin_listings.html"}}Listings{{end}}
		{{define "admin_users.html"}}Users{{end}}
		{{define "admin_delete_confirm.html"}}Delete Confirm{{end}}
		{{define "error.html"}}Error: {{.Message}}{{end}}
	`))
}
