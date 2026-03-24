package auth_test

import (
	"context"
	"html/template"
	"io"

	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
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

// MockGoogleProvider
type MockGoogleProvider struct {
	testifyMock.Mock
}

func (m *MockGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	args := m.Called(state, scheme, host)
	return args.String(0)
}

func (m *MockGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	args := m.Called(ctx, code, scheme, host)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*auth.GoogleUser, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.GoogleUser), args.Error(1)
}
