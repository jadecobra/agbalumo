package ui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"html/template"

	"github.com/labstack/echo/v4"
)

func TestTemplateRenderer_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	e := echo.New()
	c := e.NewContext(nil, nil)
	buf := new(bytes.Buffer)

	t.Run("Dict Helper Errors", func(t *testing.T) {
		os.WriteFile(filepath.Join(tempDir, "bad_dict.html"), []byte(`{{ dict "key" }}`), 0644)
		renderer, _ := NewTemplateRenderer(filepath.Join(tempDir, "*.html"))
		if err := renderer.Render(buf, "bad_dict.html", nil, c); err == nil {
			t.Error("Expected error for odd number of dict args")
		}
	})

	t.Run("Template Not Found", func(t *testing.T) {
		emptyRenderer := &TemplateRenderer{templates: make(map[string]*template.Template)}
		rec := httptest.NewRecorder()
		c2 := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
		if err := emptyRenderer.Render(rec, "nonexistent", nil, c2); err == nil {
			t.Error("Expected error for nonexistent template")
		}
	})
}
