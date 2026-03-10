package handler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTypographyConstraints(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	tailwindPath := filepath.Join(projectRoot, "tailwind.config.js")
	tailwindContent, err := os.ReadFile(tailwindPath)
	if err != nil {
		t.Fatalf("Failed to read tailwind.config.js: %v", err)
	}
	content := string(tailwindContent)
	if !strings.Contains(content, `"Inter"`) || !strings.Contains(content, `"Playfair Display"`) {
		t.Error("tailwind.config.js does not define Inter and Playfair Display fonts")
	}
}
