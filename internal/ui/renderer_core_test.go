package ui

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestTemplateRenderer_Render_Core(t *testing.T) {
	e := echo.New()

	t.Run("CSRF", func(t *testing.T) {
		tmpl := template.New("test").Funcs(BuildGlobalFuncMap())
		_, _ = tmpl.Parse(`{{.CSRF}}`)
		renderer := &TemplateRenderer{templates: map[string]*template.Template{"test": tmpl}}
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
		c.Set("csrf", "token-123")

		if err := renderer.Render(rec, "test", map[string]interface{}{}, c); err != nil {
			t.Fatal(err)
		}
		if rec.Body.String() != "token-123" {
			t.Errorf("Expected token-123, got %q", rec.Body.String())
		}
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
