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

func TestTemplateRenderer_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Test Dict Helper Errors
	// Create a template that uses dict incorrectly
	badDictTmpl := `{{ dict "key" }}` // Odd number of args
	badDictPath := filepath.Join(tempDir, "bad_dict.html")
	if err := os.WriteFile(badDictPath, []byte(badDictTmpl), 0644); err != nil {
		t.Fatal(err)
	}

	badKeyTmpl := `{{ dict 1 "value" }}` // Non-string key
	badKeyPath := filepath.Join(tempDir, "bad_key.html")
	if err := os.WriteFile(badKeyPath, []byte(badKeyTmpl), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a valid renderer first
	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	c := e.NewContext(nil, nil)
	buf := new(bytes.Buffer)

	// executing bad_dict should fail
	if err := renderer.Render(buf, "bad_dict.html", nil, c); err == nil {
		t.Error("Expected error for odd number of dict args, got nil")
	}

	// executing bad_key should fail
	if err := renderer.Render(buf, "bad_key.html", nil, c); err == nil {
		t.Error("Expected error for non-string dict key, got nil")
	}

	// 2. Test Partial Fallback (Render a partial that isn't a main page)
	// We need a base page and a partial
	baseContent := `{{define "base"}}Base: {{template "content" .}}{{end}}`
	partialContent := `{{define "mypartial"}}Partial Content{{end}}`

	// Write files
	os.WriteFile(filepath.Join(tempDir, "base.html"), []byte(baseContent), 0644)
	os.WriteFile(filepath.Join(tempDir, "partials/frag.html"), []byte(partialContent), 0644) // Need subfolder for "partials" regex?
	// The code checks strings.Contains(file, "partials").
	// So we need a filename or path with "partials".

	partDir := filepath.Join(tempDir, "partials")
	os.Mkdir(partDir, 0755)
	os.WriteFile(filepath.Join(partDir, "frag.html"), []byte(partialContent), 0644)

	// We also need a "page" to exist for the fallback loop to find a template set
	os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(`{{define "content"}}Index{{end}}`), 0644)

	// Re-init renderer with new files
	// Note: NewTemplateRenderer takes patterns.
	renderer, err = NewTemplateRenderer(filepath.Join(tempDir, "*.html"), filepath.Join(partDir, "*.html"))
	if err != nil {
		t.Fatal(err)
	}

	// Try rendering "mypartial" directly. It's not a file in the map (index.html, base.html are keys probably?)
	// base.html is in layoutFiles. index.html is in pageFiles.
	// partials/frag.html is in partialFiles.
	// templates map keys are base filenames of pageFiles. So "index.html".
	// "mypartial" is NOT a key. This triggers the fallback.

	if err := renderer.Render(buf, "mypartial", nil, c); err != nil {
		t.Errorf("Expected partial render fallback to succeed, got: %v", err)
	} else if !bytes.Contains(buf.Bytes(), []byte("Partial Content")) {
		t.Errorf("Expected 'Partial Content', got %s", buf.String())
	}

	// 3. Test NewTemplateRenderer with no files
	// Create empty dir
	emptyDir := t.TempDir()
	_, err = NewTemplateRenderer(filepath.Join(emptyDir, "*.html"))
	if err == nil {
		t.Error("Expected error when no files found, got nil")
	}

	// 4. Test Parse Error (Bad Syntax)
	badSyntaxDir := t.TempDir()
	os.WriteFile(filepath.Join(badSyntaxDir, "bad.html"), []byte(`{{ .Open `), 0644)
	_, err = NewTemplateRenderer(filepath.Join(badSyntaxDir, "*.html"))
	if err == nil {
		t.Error("Expected error for bad syntax, got nil")
	}
}
