package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

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
