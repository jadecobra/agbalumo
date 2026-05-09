package common_test

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
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
	t.Parallel()
	e := echo.New()

	// Setup simple template for testing
	funcs := ui.BuildGlobalFuncMap()
	tmpl := template.Must(template.New("about").Funcs(funcs).Parse(`{{define "about.html"}}About agbalumo{{end}}`))
	e.Renderer = &TestRenderer{templates: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := common.NewPageHandler(app)
	if err := h.HandleAbout(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "About agbalumo")
}

func TestHandleSandbox(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		env            string
		expectedStatus int
	}{
		{
			name:           "Sandbox allowed in development",
			env:            domain.EnvDevelopment,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Sandbox forbidden in production",
			env:            domain.EnvProduction,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()

			// Setup simple template for testing
			funcs := ui.BuildGlobalFuncMap()
			tmpl := template.Must(template.New("sandbox").Funcs(funcs).Parse(`{{define "sandbox.html"}}Sandbox Content{{end}}`))
			e.Renderer = &TestRenderer{templates: tmpl}

			app, cleanup := testutil.SetupTestAppEnv(t)
			defer cleanup()
			app.Cfg.Env = tt.env

			h := common.NewPageHandler(app)

			runSandboxTest(t, h, tt.expectedStatus, e)
		})
	}
}

func runSandboxTest(t *testing.T, h *common.PageHandler, expectedStatus int, e *echo.Echo) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/sandbox", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.HandleSandbox(c)
	if expectedStatus == http.StatusOK {
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sandbox Content")
		return
	}

	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(t, expectedStatus, he.Code)
		} else {
			t.Errorf("expected echo.HTTPError, got %v", err)
		}
	} else {
		assert.Equal(t, expectedStatus, rec.Code)
	}
}
