package testutil

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// TestRenderer is a simple Template Renderer for testing with an in-memory template.
type TestRenderer struct {
	Templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// RealTemplateRenderer parses actual files from ui/templates
type RealTemplateRenderer struct {
	Templates *template.Template
}

func (t *RealTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// NewMainTemplate returns a minimal template for use in unit tests.
func NewMainTemplate() *template.Template {
	return template.Must(template.New("main").Funcs(ui.BuildGlobalFuncMap()).Parse(`
		{{define "index.html"}}{{.TotalCount}} {{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "modal_detail"}}{{.Listing.Title}}{{end}}
		{{define "listing_list"}}{{range .Listings}}{{.Title}}{{end}}{{template "pagination_controls" dict "OOB" true}}{{end}}
		{{define "pagination_controls"}}{{if .OOB}}hx-swap-oob="true" id="pagination-controls"{{end}}{{end}}
		{{define "listing_card"}}{{.Listing.Title}}{{end}}
		{{define "modal_edit_listing"}}{{.Listing.Title}}{{end}}
		{{define "modal_profile"}}{{.User.Name}}{{end}}
		{{define "profile.html"}}{{.User.Name}}{{end}}
		{{define "about.html"}}About agbalumo{{end}}
		{{define "error.html"}}Error Page: {{.Message}}{{end}}
		{{define "admin_listings.html"}}{{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "admin_listing_table_row"}}<tr id="listing-row-{{.ID}}"><input type="checkbox" /></tr>{{end}}
		{{define "admin_dashboard.html"}}Admin Dashboard{{end}}
		{{define "modal_feedback.html"}}{{if .}}Feedback Modal: {{.}}{{else}}Feedback Modal{{end}}{{end}}
	`))
}

// SetupTestContext prepares a common Echo context and recorder.
func SetupTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &TestRenderer{Templates: NewMainTemplate()}
	req := httptest.NewRequest(method, target, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// NewRealTemplate returns a template object parsed from actual filesystem files.
// It includes all templates, partials, and components.
func NewRealTemplate(t *testing.T) *template.Template {
	return NewRealTemplateForPage(t, "index.html")
}

// NewRealTemplateForPage returns a template object for a specific page.
func NewRealTemplateForPage(t *testing.T, pageName string) *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Adjust path to find the ui directory correctly from various test locations
	var projectRoot string
	tempWd := wd
	for {
		if _, err := os.Stat(filepath.Join(tempWd, "ui", "templates")); err == nil {
			projectRoot = tempWd
			break
		}
		parent := filepath.Dir(tempWd)
		if parent == tempWd {
			t.Fatalf("Could not find project root (containing ui/templates) starting from %s", wd)
		}
		tempWd = parent
	}

	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")
	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

	funcMap := ui.BuildGlobalFuncMap()
	tmpl := template.New("base").Funcs(funcMap)

	paths := []string{
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "error.html"),
	}

	if pageName != "base.html" && pageName != "error.html" {
		paths = append(paths, filepath.Join(projectRoot, "ui", "templates", pageName))
	}

	tmpl, err = tmpl.ParseFiles(paths...)
	if err != nil {
		t.Fatalf("Failed to parse templates for %s: %v", pageName, err)
	}

	_, err = tmpl.ParseGlob(partialPattern)
	if err != nil {
		t.Fatalf("Failed to parse partial templates: %v", err)
	}
	_, err = tmpl.ParseGlob(componentPattern)
	if err != nil {
		t.Fatalf("Failed to parse component templates: %v", err)
	}

	return tmpl
}
