package ui

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates map[string]*template.Template
}

// NewTemplateRenderer creates a new instance of TemplateRenderer with parsed templates
func NewTemplateRenderer(patterns ...string) (*TemplateRenderer, error) {
	// 1. Identify all files
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

	// 2. Classify files
	var layoutFiles []string
	var partialFiles []string
	var pageFiles []string

	for _, file := range allFiles {
		baseName := filepath.Base(file)
		if baseName == "base.html" {
			layoutFiles = append(layoutFiles, file)
		} else if strings.Contains(file, "partials") {
			partialFiles = append(partialFiles, file)
		} else {
			pageFiles = append(pageFiles, file)
		}
	}

	// 3. Logic to create FuncMap (shared)
	funcMap := template.FuncMap{
		"split": strings.Split,
		"mod":   func(i, j int) int { return i % j },
		"add":   func(i, j int) int { return i + j },
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"toJson": func(v interface{}) (template.JS, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(b), nil
		},
	}

	// 4. Compile Templates per Page
	templates := make(map[string]*template.Template)

	for _, pageFile := range pageFiles {
		fileName := filepath.Base(pageFile)

		// Create a new template set for this page
		tmpl := template.New(fileName).Funcs(funcMap)

		// Parse Layouts (base.html) - verify if exists
		if len(layoutFiles) > 0 {
			if _, err := tmpl.ParseFiles(layoutFiles...); err != nil {
				return nil, err
			}
		}

		// Parse Partials
		if len(partialFiles) > 0 {
			if _, err := tmpl.ParseFiles(partialFiles...); err != nil {
				return nil, err
			}
		}

		// Parse the Page itself
		if _, err := tmpl.ParseFiles(pageFile); err != nil {
			return nil, err
		}

		templates[fileName] = tmpl
	}

	return &TemplateRenderer{
		templates: templates,
	}, nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Inject CSRF token if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		token := c.Get("csrf")
		viewContext["CSRF"] = token
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
