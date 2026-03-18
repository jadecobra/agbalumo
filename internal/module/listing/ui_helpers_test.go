package listing_test

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// RealTemplateRenderer parses actual files from ui/templates
type RealTemplateRenderer struct {
	templates *template.Template
}

func (t *RealTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

	projectRoot := filepath.Join(wd, "..", "..", "..")
	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")
	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

	funcMap := ui.BuildGlobalFuncMap()

	// We parse base.html, error.html and the specific page
	tmpl := template.New("base").Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles(
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "error.html"),
		filepath.Join(projectRoot, "ui", "templates", pageName),
	)
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
