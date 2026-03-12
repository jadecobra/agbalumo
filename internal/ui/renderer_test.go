package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	Sub: {{ sub 5 2 }}
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

	out := buf.String()
	if !contains(out, "Add: 3") {
		t.Errorf("Expected Add: 3 in %q", out)
	}
	if !contains(out, "Sub: 3") {
		t.Errorf("Expected Sub: 3 in %q", out)
	}
	if !contains(out, "Mod: 1") {
		t.Errorf("Expected Mod: 1 in %q", out)
	}
	if !contains(out, "Seq: [1 2 3]") {
		t.Errorf("Expected Seq: [1 2 3] in %q", out)
	}
	if !contains(out, "Dict: v") {
		t.Errorf("Expected Dict: v in %q", out)
	}
}

func TestTemplateRenderer_Render_IsNew(t *testing.T) {
	tempDir := t.TempDir()
	tmplContent := `{{ if isNew .CreatedAt }}new{{ else }}old{{ end }}`
	tmplPath := filepath.Join(tempDir, "isnew.html")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write func template: %v", err)
	}

	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	e := echo.New()
	c := e.NewContext(nil, nil)

	t.Run("recent item", func(t *testing.T) {
		buf := new(bytes.Buffer)
		data := map[string]interface{}{"CreatedAt": time.Now()}
		if err := renderer.Render(buf, "isnew.html", data, c); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		if !bytes.Contains(buf.Bytes(), []byte("new")) {
			t.Errorf("Expected 'new' for recent item, got %s", buf.String())
		}
	})

	t.Run("old item", func(t *testing.T) {
		buf := new(bytes.Buffer)
		data := map[string]interface{}{"CreatedAt": time.Now().Add(-8 * 24 * time.Hour)}
		if err := renderer.Render(buf, "isnew.html", data, c); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		if !bytes.Contains(buf.Bytes(), []byte("old")) {
			t.Errorf("Expected 'old' for old item, got %s", buf.String())
		}
	})

	t.Run("zero time", func(t *testing.T) {
		buf := new(bytes.Buffer)
		data := map[string]interface{}{"CreatedAt": time.Time{}}
		if err := renderer.Render(buf, "isnew.html", data, c); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		if !bytes.Contains(buf.Bytes(), []byte("old")) {
			t.Errorf("Expected 'old' for zero time, got %s", buf.String())
		}
	})
}

func TestTemplateRenderer_Render_ToJson(t *testing.T) {
	tempDir := t.TempDir()
	tmplContent := `{{ $j := toJson . }}{{ $j }}`
	tmplPath := filepath.Join(tempDir, "tojson.html")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write func template: %v", err)
	}

	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	e := echo.New()
	c := e.NewContext(nil, nil)

	t.Run("valid object", func(t *testing.T) {
		buf := new(bytes.Buffer)
		data := map[string]interface{}{"name": "test", "count": 42}
		if err := renderer.Render(buf, "tojson.html", data, c); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		if !bytes.Contains(buf.Bytes(), []byte(`name`)) || !bytes.Contains(buf.Bytes(), []byte(`test`)) {
			t.Errorf("Expected JSON output with name/test, got %s", buf.String())
		}
	})
}

func TestTemplateRenderer_Render_DisplayCity(t *testing.T) {
	tempDir := t.TempDir()
	tmplContent := `{{ displayCity .City .Address }}`
	tmplPath := filepath.Join(tempDir, "displaycity.html")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write func template: %v", err)
	}

	renderer, err := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	e := echo.New()
	c := e.NewContext(nil, nil)

	tests := []struct {
		name     string
		city     string
		address  string
		expected string
	}{
		{"with city", "Dallas", "123 Main St", "Dallas"},
		{"empty city standard address", "", "10051 Whitehurst Dr, Dallas, TX 75243", "Dallas"},
		{"empty city city only", "", "Houston", ""},
		{"empty city and address", "", "", ""},
		{"empty city complex address", "", "10828C Beechnut St, Houston, TX 77072", "Houston"},
		{"empty city street only", "", "123 Test St", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			data := map[string]interface{}{"City": tt.city, "Address": tt.address}
			if err := renderer.Render(buf, "displaycity.html", data, c); err != nil {
				t.Fatalf("Render failed: %v", err)
			}
			if strings.TrimSpace(buf.String()) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
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
	if err2 := renderer.Render(buf, "bad_dict.html", nil, c); err2 == nil {
		t.Error("Expected error for odd number of dict args, got nil")
	}

	// executing bad_key should fail
	if err3 := renderer.Render(buf, "bad_key.html", nil, c); err3 == nil {
		t.Error("Expected error for non-string dict key, got nil")
	}

	// 2. Test Partial Fallback (Render a partial that isn't a main page)
	// We need a base page and a partial
	baseContent := `{{define "base"}}Base: {{template "content" .}}{{end}}`
	partialContent := `{{define "mypartial"}}Partial Content{{end}}`

	// Write files
	_ = os.WriteFile(filepath.Join(tempDir, "base.html"), []byte(baseContent), 0644)

	partDir := filepath.Join(tempDir, "partials")
	_ = os.Mkdir(partDir, 0755)
	_ = os.WriteFile(filepath.Join(partDir, "frag.html"), []byte(partialContent), 0644)

	// We also need a "page" to exist for the fallback loop to find a template set
	_ = os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(`{{define "content"}}Index{{end}}`), 0644)

	// Re-init renderer with new files
	renderer, err = NewTemplateRenderer(filepath.Join(tempDir, "*.html"), filepath.Join(partDir, "*.html"))
	if err != nil {
		t.Fatal(err)
	}

	if err4 := renderer.Render(buf, "mypartial", nil, c); err4 != nil {
		t.Errorf("Expected partial render fallback to succeed, got: %v", err4)
	} else if !bytes.Contains(buf.Bytes(), []byte("Partial Content")) {
		t.Errorf("Expected 'Partial Content', got %s", buf.String())
	}

	// 5. Test ParseFiles Errors (Layouts/Partials/Pages)
	// We'll manually call compileTemplates or trigger it via NewTemplateRenderer with bad files
	
	// Create a dir with a bad layout
	badLayoutDir := t.TempDir()
	_ = os.Mkdir(filepath.Join(badLayoutDir, "partials"), 0755)
	_ = os.WriteFile(filepath.Join(badLayoutDir, "base.html"), []byte(`{{ .Oops `), 0644)
	_ = os.WriteFile(filepath.Join(badLayoutDir, "index.html"), []byte(`Index`), 0644)
	
	_, err = NewTemplateRenderer(filepath.Join(badLayoutDir, "*.html"))
	if err == nil {
		t.Error("Expected error for bad layout syntax, got nil")
	}

	// Create a dir with a bad partial
	badPartialDir := t.TempDir()
	badPartSubDir := filepath.Join(badPartialDir, "partials")
	_ = os.Mkdir(badPartSubDir, 0755)
	_ = os.WriteFile(filepath.Join(badPartialDir, "index.html"), []byte(`Index`), 0644)
	_ = os.WriteFile(filepath.Join(badPartSubDir, "bad.html"), []byte(`{{ .Oops `), 0644)

	_, err = NewTemplateRenderer(filepath.Join(badPartialDir, "*.html"), filepath.Join(badPartSubDir, "*.html"))
	if err == nil {
		t.Error("Expected error for bad partial syntax, got nil")
	}

	// 6. Test Render Error (Template Not Found)
	// We already have a renderer from previous steps
	rec := httptest.NewRecorder()
	c2 := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	
	// renderer with NO templates
	emptyRenderer := &TemplateRenderer{templates: make(map[string]*template.Template)}
	if err := emptyRenderer.Render(rec, "nonexistent", nil, c2); err == nil {
		t.Error("Expected error for nonexistent template in empty renderer, got nil")
	}
}
