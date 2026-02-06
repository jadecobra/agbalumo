package ui

import (
	"errors"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new instance of TemplateRenderer with parsed templates
func NewTemplateRenderer(patterns ...string) (*TemplateRenderer, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
		"add": func(i, j int) int { return i + j },
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
	})

	// Base to allow appending
	if _, err := tmpl.Parse("{{define \"base\"}}{{end}}"); err != nil {
		return nil, err
	}

	// Parse glob patterns
	for _, pattern := range patterns {
		if _, err := tmpl.ParseGlob(pattern); err != nil {
			return nil, err
		}
	}

	return &TemplateRenderer{
		templates: tmpl,
	}, nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
