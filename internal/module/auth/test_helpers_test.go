package auth_test

import (
	"html/template"
	"io"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type TestRenderer struct {
	templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewMainTemplate() *template.Template {
	return template.Must(template.New("listing").Funcs(ui.BuildGlobalFuncMap()).Parse(`
		{{define "error.html"}}Error Page: {{.Message}}{{end}}
	`))
}
