package common_test

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestRenderer struct {
	templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func TestHandleAbout(t *testing.T) {
	e := echo.New()

	// Setup simple template for testing
	funcs := ui.BuildGlobalFuncMap()
	tmpl := template.Must(template.New("about").Funcs(funcs).Parse(`{{define "about.html"}}About agbalumo{{end}}`))
	e.Renderer = &TestRenderer{templates: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	cfg := &config.Config{}

	h := common.NewPageHandler(repo, cfg)
	if err := h.HandleAbout(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "About agbalumo")
}
