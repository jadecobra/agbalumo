package admin_test

import (
	"html/template"
	"io"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

// AdminMockRenderer is a simple renderer for admin tests
type AdminMockRenderer struct{}

func (m *AdminMockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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
