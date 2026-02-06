package handler_test

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
)

// RealTemplateRenderer parses actual files from ui/templates
type RealTemplateRenderer struct {
	templates *template.Template
}

func (t *RealTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRealTemplate(t *testing.T) *template.Template {
	// Locate project root assuming we run from internal/handler
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Traverse up to find go.mod or specific root marker
	// Simpler approach: assume standard level depth or use relative ../../
	// If running from internal/handler, root is ../../
	projectRoot := filepath.Join(wd, "..", "..")
	templatePattern := filepath.Join(projectRoot, "ui", "templates", "*.html")
	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")

	// Helper functions map matching what's in main.go
	funcMap := template.FuncMap{
		"mod": func(i, j int) int { return i % j },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, nil // simplified
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil // simplified
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	tmpl, err := template.New("base").Funcs(funcMap).ParseGlob(templatePattern)
	if err != nil {
		t.Fatalf("Failed to parse base templates: %v", err)
	}
	_, err = tmpl.ParseGlob(partialPattern)
	if err != nil {
		t.Fatalf("Failed to parse partial templates: %v", err)
	}

	return tmpl
}

func TestHomePageUIValues(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			return []domain.Listing{}, nil
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	body := rec.Body.String()

	// 1. Verify Listings Container is Relative
	// Expected: <div class="relative flex-1 px-4 ...">
	if !strings.Contains(body, `class="relative flex-1`) {
		t.Error("Regression: Listings container missing 'relative' class")
	}

	// 2. Verify Helper Overlay classes
	// Expected: opacity-0 pointer-events-none
	if !strings.Contains(body, `opacity-0`) || !strings.Contains(body, `pointer-events-none`) {
		t.Error("Regression: Loading indicator missing default transparency/pointer-events classes")
	}

	// 3. Verify HTMX Interaction classes
	// Expected: [&.htmx-request]:opacity-100
	if !strings.Contains(body, `[&.htmx-request]:opacity-100`) {
		t.Error("Regression: Loading indicator missing active state classes")
	}

	// 4. Verify Search Indicator
	// Expected: hx-indicator="#listings-loading"
	if !strings.Contains(body, `hx-indicator="#listings-loading"`) {
		t.Error("Regression: Search input missing hx-indicator attribute")
	}
}
