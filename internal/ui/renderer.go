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

type TemplateRenderer struct {
	templates      map[string]*template.Template
	funcMap        template.FuncMap
	patterns       []string
	CountryRegions []Region
}

type Region struct {
	Region    string    `json:"region"`
	Countries []Country `json:"countries"`
}

type Country struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
}

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
	funcMap["countryFlag"] = func(name string) string {
		return getCountryFlag(renderer.CountryRegions, name)
	}
	funcMap["isCountry"] = func(name string) string {
		if isCountry(renderer.CountryRegions, name) {
			return "true"
		}
		return ""
	}

	templates, err := compileTemplates(layoutFiles, partialFiles, pageFiles, funcMap)
	if err != nil {
		return nil, err
	}

	renderer.patterns = patterns
	renderer.funcMap = funcMap
	renderer.templates = templates
	return renderer, nil
}

// GetFuncMap returns the template function map
func (t *TemplateRenderer) GetFuncMap() template.FuncMap {
	return t.funcMap
}

// GetAllTemplateNames returns all defined template names (pages, layouts, partials)
func (t *TemplateRenderer) GetAllTemplateNames() []string {
	uniqueNames := make(map[string]bool)
	for _, tmpl := range t.templates {
		for _, def := range tmpl.Templates() {
			uniqueNames[def.Name()] = true
		}
	}
	var names []string
	for name := range uniqueNames {
		names = append(names, name)
	}
	return names
}

// RenderDefinition searches all template sets for the named definition and renders it.
// Used primarily for testing partials and page-specific blocks.
func (t *TemplateRenderer) RenderDefinition(w io.Writer, name string, data interface{}) error {
	for _, tmpl := range t.templates {
		if d := tmpl.Lookup(name); d != nil {
			return tmpl.ExecuteTemplate(w, name, data)
		}
	}
	return fmt.Errorf("template definition %q not found in any template set", name)
}

func (t *TemplateRenderer) loadCountryData() error {
	paths := []string{
		"ui/data/countries.json",
		"../../../ui/data/countries.json",
		"../../ui/data/countries.json",
		"../ui/data/countries.json",
	}

	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(filepath.Clean(p))
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}
	return json.Unmarshal(data, &t.CountryRegions)
}

func categorizeTemplateFiles(files []string) (layouts, partials, pages []string) {
	for _, file := range files {
		ext := filepath.Ext(file)
		if ext != ".html" {
			continue
		}

		dir := filepath.Base(filepath.Dir(file))
		baseName := filepath.Base(file)
		if dir == "layouts" || baseName == "base.html" {
			layouts = append(layouts, file)
		} else if dir == "partials" || dir == "components" {
			partials = append(partials, file)
		} else {
			pages = append(pages, file)
		}
	}
	return layouts, partials, pages
}

func BuildGlobalFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"div": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"split":            strings.Split,
		"seq":              seq,
		"dict":             dict,
		"toJson":           toJson,
		"isNew":            isNew,
		"safeHTML":         safeHTML,
		"safeHTMLAttr":     safeHTMLAttr,
		"safeJS":           safeJS,
		"displayCity":      displayCity,
		"fallbackImageURL": fallbackImageURL,
		"hasDelivery":      hasDelivery,
		"float64": func(i int) float64 {
			return float64(i)
		},
		"subF": func(a, b float64) float64 {
			return a - b
		},
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
	tmpl := template.New(filepath.Base(page)).Funcs(funcMap).Option("missingkey=error")

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

func (t *TemplateRenderer) recompileTemplates() {
	var allFiles []string
	for _, pattern := range t.patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		allFiles = append(allFiles, files...)
	}

	layouts, partials, pages := categorizeTemplateFiles(allFiles)
	templates, err := compileTemplates(layouts, partials, pages, t.funcMap)
	if err != nil {
		slog.Warn("Template recompilation failed", "error", err)
		return
	}
	t.templates = templates
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t.maybeRecompile()
	t.injectContext(data, c)

	// Semantic Tagging for Agentic Discovery (Non-Production Only)
	if os.Getenv(domain.EnvKeyAppEnv) != domain.EnvProduction {
		_, _ = fmt.Fprintf(w, "<!-- BEGIN TEMPLATE: %s -->", name)
	}

	tmpl, ok := t.templates[name]
	if !ok {
		return t.renderFallback(w, name, data)
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		slog.Error("Template execution failed", "template", name, "error", err)
	}
	return err
}

func (t *TemplateRenderer) maybeRecompile() {
	if os.Getenv(domain.EnvKeyAppEnv) != domain.EnvProduction && len(t.patterns) > 0 {
		t.recompileTemplates()
	}
}

func (t *TemplateRenderer) injectContext(data interface{}, c echo.Context) {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		token := c.Get("csrf")
		viewContext["CSRF"] = token
		viewContext["CountryRegions"] = t.CountryRegions
	}
}

func (t *TemplateRenderer) renderFallback(w io.Writer, name string, data interface{}) error {
	var tmpl *template.Template
	for _, tSet := range t.templates {
		tmpl = tSet
		break
	}

	if tmpl == nil {
		return errors.New("template not found and no default template available: " + name)
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		slog.Error("Template execution failed (fallback)", "template", name, "error", err)
	}
	return err
}
