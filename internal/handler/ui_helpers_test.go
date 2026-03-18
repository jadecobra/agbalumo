package handler_test

import (
	"context"
	"html/template"
	"io"
	"mime/multipart"

	"github.com/jadecobra/agbalumo/internal/ui"
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
	return template.Must(template.New("listing").Funcs(ui.BuildGlobalFuncMap()).Parse(`
		{{define "index.html"}}{{.TotalCount}} {{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "modal_detail"}}{{.Listing.Title}}{{end}}
		{{define "listing_list"}}{{range .Listings}}{{.Title}}{{end}}{{template "pagination_controls" dict "OOB" true}}{{end}}
		{{define "pagination_controls"}}{{if .OOB}}hx-swap-oob="true" id="pagination-controls"{{end}}{{end}}
		{{define "listing_card"}}{{.Listing.Title}}{{end}}
		{{define "modal_edit_listing"}}{{.Listing.Title}}{{end}}
		{{define "modal_profile"}}{{.User.Name}}{{end}}
		{{define "profile.html"}}{{.User.Name}}{{end}}
		{{define "about.html"}}About agbalumo{{end}}
		{{define "error.html"}}Error Page: {{.Message}}{{end}}
		{{define "admin_listings.html"}}{{range .Listings}}{{.Title}}{{end}}{{end}}
		{{define "admin_listing_table_row"}}<tr id="listing-row-{{.ID}}"><input type="checkbox" /></tr>{{end}}
		{{define "admin_dashboard.html"}}Admin Dashboard{{end}}
		{{define "modal_feedback.html"}}Feedback Modal{{end}}
	`))
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
