package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// Country represents a country with a name and flag emoji.
type Country struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// Region represents a geographical region containing countries.
type Region struct {
	Region    string    `json:"region"`
	Countries []Country `json:"countries"`
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates      map[string]*template.Template
	CountryRegions []Region
}

// NewTemplateRenderer creates a new instance of TemplateRenderer with parsed templates
func NewTemplateRenderer(patterns ...string) (*TemplateRenderer, error) {
	var allFiles []string
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, files...)
	}

	if len(allFiles) == 0 {
		return nil, errors.New("no template files found")
	}

	layoutFiles, partialFiles, pageFiles := categorizeTemplateFiles(allFiles)

	renderer := &TemplateRenderer{}
	if err := renderer.loadCountryData(); err != nil {
		slog.Warn("Failed to load country data", "error", err)
	}

	funcMap := BuildGlobalFuncMap()
	funcMap["Countries"] = func() []Region {
		return renderer.CountryRegions
	}

	templates, err := compileTemplates(layoutFiles, partialFiles, pageFiles, funcMap)
	if err != nil {
		return nil, err
	}

	renderer.templates = templates
	return renderer, nil
}

func (t *TemplateRenderer) loadCountryData() error {
	data, err := os.ReadFile("ui/data/countries.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.CountryRegions)
}

func categorizeTemplateFiles(files []string) (layouts, partials, pages []string) {
	for _, file := range files {
		baseName := filepath.Base(file)
		if baseName == "base.html" {
			layouts = append(layouts, file)
		} else if strings.Contains(file, "partials") || strings.Contains(file, "components") {
			partials = append(partials, file)
		} else {
			pages = append(pages, file)
		}
	}
	return
}

func BuildGlobalFuncMap() template.FuncMap {
	return template.FuncMap{
		"split":       strings.Split,
		"mod":         func(i, j int) int { return i % j },
		"add":         func(i, j int) int { return i + j },
		"sub":         func(i, j int) int { return i - j },
		"seq":         seq,
		"dict":        dict,
		"toJson":      toJson,
		"isNew":       isNew,
		"safeHTML":    safeHTML,
		"displayCity": displayCity,
		"Countries": func() []Region {
			return nil
		},
	}
}

func compileTemplates(layouts, partials, pages []string, funcMap template.FuncMap) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	for _, pageFile := range pages {
		tmpl, err := parseTemplateFiles(pageFile, layouts, partials, funcMap)
		if err != nil {
			return nil, err
		}
		templates[filepath.Base(pageFile)] = tmpl
	}
	return templates, nil
}

func parseTemplateFiles(page string, layouts, partials []string, funcMap template.FuncMap) (*template.Template, error) {
	tmpl := template.New(filepath.Base(page)).Funcs(funcMap)

	if len(layouts) > 0 {
		if _, err := tmpl.ParseFiles(layouts...); err != nil {
			return nil, err
		}
	}

	if len(partials) > 0 {
		if _, err := tmpl.ParseFiles(partials...); err != nil {
			return nil, err
		}
	}

	return tmpl.ParseFiles(page)
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Inject CSRF token if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		token := c.Get("csrf")
		viewContext["CSRF"] = token
		viewContext["CountryRegions"] = t.CountryRegions
	}

	// Semantic Tagging for Agentic Discovery (Non-Production Only)
	if os.Getenv(domain.EnvKeyAppEnv) != domain.EnvProduction {
		_, _ = fmt.Fprintf(w, "<!-- BEGIN TEMPLATE: %s -->", name)
	}

	tmpl, ok := t.templates[name]
	if !ok {
		// Fallback: Check if it's a partial by trying to execute it on a default template set
		// We can use any existing template set because they all include all partials.
		// Let's try to find "index.html" or just use the first available one.
		for _, t := range t.templates {
			tmpl = t
			break
		}
		if tmpl == nil {
			return errors.New("template not found and no default template available: " + name)
		}
		// Try to execute the named partial on this template set
		// Note: ExecuteTemplate returns error if name is not found
		return tmpl.ExecuteTemplate(w, name, data)
	}
	return tmpl.ExecuteTemplate(w, name, data)
}
