package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTemplateRenderer_Render_Core(t *testing.T) {
	t.Parallel()
	e := echo.New()

	t.Run("CSRF", func(t *testing.T) {
		t.Parallel()
		tmpl := template.New("test").Funcs(BuildGlobalFuncMap())
		_, _ = tmpl.Parse(`{{.CSRF}}`)
		renderer := &TemplateRenderer{templates: map[string]*template.Template{"test": tmpl}}
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
		c.Set("csrf", "token-123")

		if err := renderer.Render(rec, "test", map[string]interface{}{}, c); err != nil {
			t.Fatal(err)
		}
		actual := rec.Body.String()
		if actual != "token-123" && actual != "<!-- BEGIN TEMPLATE: test -->token-123" {
			t.Errorf("Expected token-123 or tagged version, got %q", actual)
		}
	})

	t.Run("ComponentAttributes", func(t *testing.T) {
		t.Parallel()
		tmpl := template.New("test-comp").Funcs(BuildGlobalFuncMap())
		// Test rendering button_sharp with custom attributes
		_, _ = tmpl.Parse(`{{template "button_sharp" dict "Label" "Test" "Attr" "data-testid=\"btn-123\""}}
{{define "button_sharp"}}<button {{if .Attr}}{{.Attr | safeHTMLAttr}}{{end}}>{{.Label}}</button>{{end}}`)

		renderer := &TemplateRenderer{templates: map[string]*template.Template{"test-comp": tmpl}}
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

		if err := renderer.Render(rec, "test-comp", map[string]interface{}{}, c); err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, rec.Body.String(), `data-testid="btn-123"`)
	})
}

func BenchmarkRender(b *testing.B) {
	tmpl := template.New("bench").Funcs(BuildGlobalFuncMap())
	_, _ = tmpl.Parse(`<h1>{{.Title}}</h1><p>{{.Description}}</p>`)
	renderer := &TemplateRenderer{templates: map[string]*template.Template{"bench": tmpl}}
	e := echo.New()
	c := e.NewContext(nil, nil)
	data := map[string]interface{}{"Title": "Bench", "Description": "Desc"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = renderer.Render(&buf, "bench", data, c)
	}
}

func TestTemplateRenderer_HotReload(t *testing.T) {
	// 1. Create a temp directory for templates
	tempDir := t.TempDir()
	tmplPath := filepath.Join(tempDir, "test.html")

	// 2. Write initial template content
	err := os.WriteFile(tmplPath, []byte("hello"), 0600)
	assert.NoError(t, err)

	// 3. Create renderer with the temp dir pattern
	pattern := filepath.Join(tempDir, "*.html")
	renderer, err := NewTemplateRenderer(pattern)
	assert.NoError(t, err)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	// 4. Set env to test (not production) to trigger recompile
	t.Setenv(domain.EnvKeyAppEnv, domain.EnvDevelopment)

	// 5. Render -> assert "hello"
	err = renderer.Render(rec, "test.html", map[string]interface{}{}, c)
	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "hello")

	// 6. Overwrite template file content
	err = os.WriteFile(tmplPath, []byte("goodbye"), 0600)
	assert.NoError(t, err)

	// 7. Render again -> assert "goodbye" (proves recompilation happened)
	rec2 := httptest.NewRecorder()
	err = renderer.Render(rec2, "test.html", map[string]interface{}{}, c)
	assert.NoError(t, err)
	assert.Contains(t, rec2.Body.String(), "goodbye")
}
