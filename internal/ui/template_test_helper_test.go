package ui

import (
	"bytes"
	"html/template"
	"testing"
)

// renderPartial loads the real project templates and executes a named partial with data.
// This enables testing partials in isolation without a server.
func renderPartial(t *testing.T, name string, data interface{}) string {
	t.Helper()
	renderer, err := NewTemplateRenderer(
		"../../ui/templates/*.html",
		"../../ui/templates/partials/*.html",
		"../../ui/templates/components/*.html",
	)
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}
	// Find any template set (all include all partials)
	var tmpl *template.Template
	for _, t := range renderer.templates {
		tmpl = t
		break
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		t.Fatalf("Failed to render %q: %v", name, err)
	}
	return buf.String()
}
