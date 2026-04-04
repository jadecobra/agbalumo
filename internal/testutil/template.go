package testutil

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
	Templates *template.Template
}

func (t *RealTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
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

	// Adjust path to find the ui directory correctly from module tests or util
	var projectRoot string
	if filepath.Base(wd) == "agbalumo" {
		projectRoot = "."
	} else if filepath.Base(filepath.Dir(wd)) == "module" {
		projectRoot = filepath.Join(wd, "..", "..", "..")
	} else if filepath.Base(wd) == "testutil" {
		projectRoot = filepath.Join(wd, "..", "..")
	} else {
		projectRoot = filepath.Join(wd, "..", "..")
	}

	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")
	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

	funcMap := ui.BuildGlobalFuncMap()

	// We parse base.html, error.html and the specific page
	tmpl := template.New("base").Funcs(funcMap)

	// Add potential relative paths
	paths := []string{
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "error.html"),
	}

	// Add the specific page if it's not base or error
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
