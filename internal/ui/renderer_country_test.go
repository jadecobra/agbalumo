package ui

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestTemplateRenderer_CountryData(t *testing.T) {
	e := echo.New()

	t.Run("InjectedCountryData", func(t *testing.T) {
		// Setup renderer with mock data and custom funcMap
		renderer := &TemplateRenderer{
			CountryRegions: []Region{
				{
					Region: "Region1",
					Countries: []Country{
						{Name: "Country1", Flag: "F1"},
					},
				},
			},
		}

		funcMap := BuildGlobalFuncMap()
		funcMap["Countries"] = func() []Region {
			return renderer.CountryRegions
		}

		tmpl := template.New("test").Funcs(funcMap)
		_, _ = tmpl.Parse(`{{range Countries}}{{.Region}}:{{range .Countries}}{{.Name}},{{end}}{{end}}`)
		renderer.templates = map[string]*template.Template{"test": tmpl}

		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

		err := renderer.Render(rec, "test", map[string]interface{}{}, c)
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}

		expected := "Region1:Country1,"
		if rec.Body.String() != expected {
			t.Errorf("Expected %q, got %q. Country data might not be injected.", expected, rec.Body.String())
		}
	})

	t.Run("CheckCategorizeFiles", func(t *testing.T) {
		layouts, partials, pages := categorizeTemplateFiles([]string{
			"ui/templates/base.html",
			"ui/templates/components/foo.html",
			"ui/templates/index.html",
		})
		if len(layouts) != 1 || len(partials) != 1 || len(pages) != 1 {
			t.Errorf("Categorization failed: %v, %v, %v", layouts, partials, pages)
		}
	})
}
