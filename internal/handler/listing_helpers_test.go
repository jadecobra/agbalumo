package handler_test

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

// Simple Template Renderer for testing
type TestRenderer struct {
	templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewMainTemplate() *template.Template {
	funcMap := template.FuncMap{
		"mod":   func(i, j int) int { return i % j },
		"add":   func(i, j int) int { return i + j },
		"sub":   func(i, j int) int { return i - j },
		"split": strings.Split,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, nil
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"toJson": func(v interface{}) (template.JS, error) {
			b, jErr := json.Marshal(v)
			if jErr != nil {
				return "", jErr
			}
			return template.JS(b), nil
		},
		"isNew": func(createdAt time.Time) bool {
			if createdAt.IsZero() {
				return false
			}
			return time.Since(createdAt) < 7*24*time.Hour
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	return template.Must(template.New("listing").Funcs(funcMap).Parse(`
		{{define "index.html"}}{{.TotalCount}} {{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "modal_detail"}}{{.Listing.Title}}{{end}}
		{{define "listing_list"}}{{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "listing_card"}}{{.Listing.Title}}{{end}}
		{{define "modal_edit_listing"}}{{.Listing.Title}}{{end}}
		{{define "modal_profile"}}{{.User.Name}}{{end}}
		{{define "profile.html"}}{{.User.Name}}{{end}}
		{{define "about.html"}}About agbalumo{{end}}
		{{define "error.html"}}Error Page: {{.Message}}{{end}}
		{{define "admin_listings.html"}}{{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "admin_dashboard.html"}}Admin Dashboard{{end}}
		{{define "modal_feedback.html"}}Feedback Modal{{end}}
	`))
}

func setupTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(method, target, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func setupRequest(method, target string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, target, body)
}

func setupResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// MockImageService for testing
type MockImageService struct {
	testifyMock.Mock
}

func (m *MockImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	args := m.Called(ctx, file, listingID)
	return args.String(0), args.Error(1)
}

func (m *MockImageService) DeleteImage(ctx context.Context, imageURL string) error {
	args := m.Called(ctx, imageURL)
	return args.Error(0)
}
