package ui

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestTemplateRenderer_Funcs(t *testing.T) {
	tempDir := t.TempDir()
	e := echo.New()
	c := e.NewContext(nil, nil)

	t.Run("Basic Arithmetic", func(t *testing.T) {
		tmplContent := `Add: {{ add 1 2 }} Sub: {{ sub 5 2 }} Mod: {{ mod 10 3 }} Seq: {{ seq 1 3 }} Dict: {{ $d := dict "k" "v" }}{{ $d.k }}`
		renderer := setupRenderer(t, tempDir, "funcs.html", tmplContent)
		buf := new(bytes.Buffer)
		if err := renderer.Render(buf, "funcs.html", nil, c); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "Add: 3") || !strings.Contains(out, "Sub: 3") || !strings.Contains(out, "Mod: 1") || !strings.Contains(out, "Seq: [1 2 3]") || !strings.Contains(out, "Dict: v") {
			t.Errorf("Unexpected output: %q", out)
		}
	})

	verifyOtherTemplateFuncs(t, tempDir, c)
}

func verifyOtherTemplateFuncs(t *testing.T, tempDir string, c echo.Context) {
	t.Run("isNew", func(t *testing.T) {
		tmplContent := `{{ if isNew .CreatedAt }}new{{ else }}old{{ end }}`
		renderer := setupRenderer(t, tempDir, "isnew.html", tmplContent)

		buf := new(bytes.Buffer)
		_ = renderer.Render(buf, "isnew.html", map[string]interface{}{"CreatedAt": time.Now()}, c)
		if !bytes.Contains(buf.Bytes(), []byte("new")) {
			t.Errorf("Expected 'new', got %s", buf.String())
		}
	})

	t.Run("toJson", func(t *testing.T) {
		tmplContent := `<script>{{ toJson . }}</script>`
		renderer := setupRenderer(t, tempDir, "tojson.html", tmplContent)
		buf := new(bytes.Buffer)
		_ = renderer.Render(buf, "tojson.html", map[string]interface{}{"name": "test"}, c)
		if !strings.Contains(buf.String(), `"name":"test"`) {
			t.Errorf("Expected JSON, got %s", buf.String())
		}
	})

	t.Run("displayCity", func(t *testing.T) {
		tmplContent := `City: {{ displayCity .City .Address }}`
		renderer := setupRenderer(t, tempDir, "city.html", tmplContent)

		tests := []struct {
			City    string
			Address string
			Want    string
		}{
			{"Lagos", "", "Lagos"},
			{"", "123 St, Accra, Ghana", "Accra"},
			{"", "123 St", ""},
			{"", "", ""},
		}

		for _, tt := range tests {
			buf := new(bytes.Buffer)
			_ = renderer.Render(buf, "city.html", map[string]interface{}{"City": tt.City, "Address": tt.Address}, c)
			if !strings.Contains(buf.String(), "City: "+tt.Want) {
				t.Errorf("displayCity(%q, %q) = %q, want %q", tt.City, tt.Address, buf.String(), tt.Want)
			}
		}
	})
}

func setupRenderer(t *testing.T, dir, name, content string) *TemplateRenderer {
	t.Helper()
	path := filepath.Join(dir, name)
	_ = os.WriteFile(filepath.Clean(path), []byte(content), 0600)
	r, _ := NewTemplateRenderer(filepath.Join(dir, "*.html"))
	return r
}
