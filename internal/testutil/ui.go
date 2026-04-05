package testutil

import (
	"html/template"
	"io"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// TestRenderer is a simple Template Renderer for testing
type TestRenderer struct {
	Templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// NewMainTemplate returns a minimal template for use in integration tests
func NewMainTemplate() *template.Template {
	return template.Must(template.New("main").Funcs(ui.BuildGlobalFuncMap()).Parse(`
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
		{{define "modal_feedback.html"}}{{if .}}Feedback Modal: {{.}}{{else}}Feedback Modal{{end}}{{end}}
	`))
}
