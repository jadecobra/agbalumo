package testutil

import (
	"html/template"
	"io"
	"github.com/gorilla/sessions"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"strings"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
)

// AssertContainsPagination verifies that the response contains the pagination controls.
func AssertContainsPagination(t testing.TB, body string) {
	t.Helper()
	if !strings.Contains(body, `hx-swap-oob="true"`) || !strings.Contains(body, `id="pagination-controls"`) {
		t.Error("response missing pagination controls")
	}
}

// AssertErrorPage verifies that the response contains the error page content.
func AssertErrorPage(t testing.TB, body string) {
	t.Helper()
	if !strings.Contains(body, "Error Page") {
		t.Error("response missing error page content")
	}
}

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
		{{define "` + domain.TemplateIndex + `"}}{{.TotalCount}} {{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "modal_detail"}}{{.Listing.Title}}{{end}}
		{{define "listing_list"}}{{range .Listings}}{{.Title}}{{end}}{{template "pagination_controls" dict "OOB" true}}{{end}}
		{{define "pagination_controls"}}{{if .OOB}}hx-swap-oob="true" id="pagination-controls"{{end}}{{end}}
		{{define "listing_card"}}{{.Listing.Title}}{{end}}
		{{define "modal_edit_listing"}}{{.Listing.Title}}{{end}}
		{{define "modal_profile"}}{{.User.Name}}{{end}}
		{{define "profile.html"}}{{.User.Name}}{{end}}
		{{define "about.html"}}About agbalumo{{end}}
		{{define "` + domain.TemplateError + `"}}Error Page: {{.Message}}{{end}}
		{{define "admin_listings.html"}}{{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "admin_listing_table_row"}}<tr id="listing-row-{{.ID}}"><input type="checkbox" /></tr>{{end}}
		{{define "admin_dashboard.html"}}Admin Dashboard{{end}}
		{{define "modal_feedback.html"}}{{if .}}Feedback Modal: {{.}}{{else}}Feedback Modal{{end}}{{end}}
	`))
}

// SetupTestContext prepares a basic Echo context for testing.
func SetupTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &TestRenderer{Templates: NewMainTemplate()}
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
 
// SetupTestContextWithSession prepares an Echo context with a functional session store.
func SetupTestContextWithSession(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	c, rec := SetupTestContext(method, target, body)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)
	return c, rec
}

// GetAuthSession returns the session associated with the context.
func GetAuthSession(c echo.Context) (*sessions.Session, error) {
	s, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return nil, domain.ErrLoginRequired
	}
	return s, nil
}

// NewRealTemplate returns a template object parsed from actual filesystem files.
// It includes all templates, partials, and components.
func NewRealTemplate(t *testing.T) *template.Template {
	return NewRealTemplateForPage(t, domain.TemplateIndex)
}

// NewRealTemplateForPage returns a template object for a specific page.
func NewRealTemplateForPage(t *testing.T, pageName string) *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	projectRoot := findProjectRoot(t, wd)

	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")
	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

	funcMap := ui.BuildGlobalFuncMap()
	tmpl := template.New("base").Funcs(funcMap)

	paths := []string{
		filepath.Join(projectRoot, "ui", "templates", domain.TemplateBase),
		filepath.Join(projectRoot, "ui", "templates", domain.TemplateError),
	}

	if pageName != domain.TemplateBase && pageName != domain.TemplateError {
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

func findProjectRoot(t *testing.T, wd string) string {
	tempWd := wd
	for {
		if _, err := os.Stat(filepath.Join(tempWd, "ui", "templates")); err == nil {
			return tempWd
		}
		parent := filepath.Dir(tempWd)
		if parent == tempWd {
			t.Fatalf("Could not find project root (containing ui/templates) starting from %s", wd)
		}
		tempWd = parent
	}
}
