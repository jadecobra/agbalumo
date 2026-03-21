package ui

import (
	"errors"
	"html/template"
	"io"
	"testing"
)

func TestAllTemplates_SyntaxAndEscaping(t *testing.T) {
	// Initialize TemplateRenderer using the actual template paths relative to this test folder (internal/ui)
	renderer, err := NewTemplateRenderer(
		"../../ui/templates/*.html",
		"../../ui/templates/partials/*.html",
		"../../ui/templates/components/*.html",
		"../../ui/templates/listings/*.html",
		"../../ui/templates/about.html",
	)
	if err != nil {
		t.Fatalf("Failed to parse project templates. (Check file paths?): %v", err)
	}

	for name, tmpl := range renderer.templates {
		t.Run(name, func(t *testing.T) {
			// Trigger the context-aware HTML escaper
			err := tmpl.Execute(io.Discard, nil)
			if err != nil {
				// We ONLY care about HTML escaping errors (syntax errors)
				var htmlErr *template.Error
				if errors.As(err, &htmlErr) {
					t.Fatalf("🚨 HTML escaping error in template %q:\n%v\nThis indicates missing quotes, unclosed tags, or invalid HTML structure.", name, htmlErr)
				}
				// Ignore other errors (e.g. template.ExecError) caused by nil data
				t.Logf("Expected execution error for %q (ignored): %v", name, err)
			}
		})
	}
}
