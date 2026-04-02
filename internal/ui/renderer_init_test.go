package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewTemplateRenderer(t *testing.T) {
	tempDir := t.TempDir()

	tmplPath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(tmplPath, []byte(`{{define "test"}}Hello Test{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to write temp template: %v", err)
	}

	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	if renderer == nil {
		t.Fatal("Renderer is nil")
	}

	if len(renderer.templates) == 0 {
		t.Error("Renderer templates map is empty")
	}

	_, err = NewTemplateRenderer(filepath.Join(tempDir, "nonexistent/*.html"))
	if err == nil {
		t.Error("Expected error for non-existent pattern, got nil")
	}
}
