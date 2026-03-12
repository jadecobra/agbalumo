package handler_test

import (
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

	projectRoot := filepath.Join(wd, "..", "..")
	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")
	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

	funcMap := template.FuncMap{
		"mod":   func(i, j int) int { return i % j },
		"add":   func(i, j int) int { return i + j },
		"sub":   func(i, j int) int { return i - j },
		"split": strings.Split,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, nil
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"toJson": func(v interface{}) (template.JS, error) {
			b, jErr := json.Marshal(v)
			if jErr != nil {
				return "", jErr
			}
			return template.JS(b), nil
		},
		"isNew": func(createdAt time.Time) bool {
			if createdAt.IsZero() {
				return false
			}
			return time.Since(createdAt) < 7*24*time.Hour
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

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
