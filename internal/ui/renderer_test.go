package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewTemplateRenderer(t *testing.T) {
	// Create temporary template directory structure for testing
	tempDir := t.TempDir()

	// Create a dummy template file
	tmplPath := filepath.Join(tempDir, "test.html")
	if err := os.WriteFile(tmplPath, []byte(`{{define "test"}}Hello Test{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to write temp template: %v", err)
	}

	// Test Success
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

	// Test Fail - Bad Pattern
	_, err = NewTemplateRenderer(filepath.Join(tempDir, "nonexistent/*.html"))
	if err == nil {
		t.Error("Expected error for non-existent pattern, got nil")
	}
}

func TestTemplateRenderer_Render(t *testing.T) {
	tempDir := t.TempDir()
	tmplContent := `
	{{- /* Test FuncMap */ -}}
	Add: {{ add 1 2 }}
	Mod: {{ mod 10 3 }}
	Seq: {{ seq 1 3 }}
	Dict: {{ $d := dict "k" "v" }}{{ $d.k }}
	`
	tmplPath := filepath.Join(tempDir, "funcs.html")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write func template: %v", err)
	}

	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	buf := new(bytes.Buffer)
	e := echo.New()
	c := e.NewContext(nil, nil)

	if err := renderer.Render(buf, "funcs.html", nil, c); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify Output
	out := buf.String()
	// weak check contains
	if !contains(out, "Add: 3") {
		t.Errorf("Expected Add: 3 in %q", out)
	}
	if !contains(out, "Mod: 1") {
		t.Errorf("Expected Mod: 1 in %q", out)
	}
	// Seq returns slice, print might be [1 2 3]
	if !contains(out, "Seq: [1 2 3]") {
		t.Errorf("Expected Seq: [1 2 3] in %q", out)
	}
	if !contains(out, "Dict: v") {
		t.Errorf("Expected Dict: v in %q", out)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func BenchmarkRender(b *testing.B) {
	tmpl := template.New("bench")
	_, _ = tmpl.Parse(`<h1>{{.Title}}</h1><p>{{.Description}}</p>`)
	
	// Mock the map structure
	templates := map[string]*template.Template{
		"bench": tmpl,
	}
	renderer := &TemplateRenderer{templates: templates}

	e := echo.New()
	c := e.NewContext(nil, nil)
	data := map[string]interface{}{
		"Title":       "Benchmark Listing",
		"Description": "This is a test description for benchmarking template rendering performance.",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = renderer.Render(&buf, "bench", data, c)
	}
}

func TestTemplateRenderer_Render_CSRF(t *testing.T) {
	tmpl := template.New("test")
	_, err := tmpl.Parse(`{{.CSRF}}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	renderer := &TemplateRenderer{
		templates: map[string]*template.Template{
			"test": tmpl,
		},
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	c.Set("csrf", "test-token-123")

	data := map[string]interface{}{
		"Title": "Test Page",
	}

	err = renderer.Render(rec, "test", data, c)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if rec.Body.String() != "test-token-123" {
		t.Errorf("Expected CSRF token 'test-token-123', got '%s'", rec.Body.String())
	}
}
