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
	"time"

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
		"split": strings.Split,
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

func TestFilterUIValues(t *testing.T) {
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

	body := rec.Body.String()

	// 1. Verify Filter Container ID
	if !strings.Contains(body, `id="filter-container"`) {
		t.Error("Regression: Filter container missing ID 'filter-container'")
	}

	// 2. Verify 'All' button has active state by default
	// Active state class unique part: bg-stone-900 text-white
	// We check for the specific combination on the All button
	if !strings.Contains(body, `All`) || !strings.Contains(body, `bg-stone-900 text-white`) {
		t.Error("Regression: 'All' button does not appear to have active classes")
	}

	// 3. Verify 'Food' button has inactive state by default
	// Inactive state includes: bg-white dark:bg-surface-dark
	if !strings.Contains(body, `Food (`) {
		t.Error("Regression: 'Food' filter button missing count")
	}
	// This check is a bit loose but ensures we didn't accidentally make Food active
	foodIndex := strings.Index(body, `Food (`)
	if foodIndex != -1 {
		// Look backwards for the button class definition
		buttonStart := strings.LastIndex(body[:foodIndex], `<button`)
		buttonTag := body[buttonStart:foodIndex]
		if strings.Contains(buttonTag, `bg-stone-900`) {
			t.Error("Regression: 'Food' button seems to have active class by default")
		}
	}


	// 4. Verify Delegated Script Logic Presence
	// We check for the specific code logic to be robust
	if !strings.Contains(body, `event.target.closest('#filter-container button')`) {
		t.Errorf("Regression: Delegated filter logic script missing. Body snippet: %s", body[len(body)-500:])
	}
}

func TestJobListingUI(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}

	job := domain.Listing{
		ID:           "job-123",
		Title:        "Senior Go Engineer",
		Company:      "Tech Innovators",
		PayRange:     "$120k - $160k",
		Type:         domain.Job,
		Description:  "Join our team to build amazing Go applications.",
		Skills:       "Go,Docker,K8s",
		JobStartDate: time.Now().Add(24 * time.Hour),
		OwnerOrigin:  "Nigeria",
		City:         "Lagos",
		ContactEmail: "jobs@techinnovators.com",
	}

	// 1. Verify Listing Card Rendering
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// We'll render the listing card partial directly for precision
	data := map[string]interface{}{
		"Listing":   job,
		"GridClass": "",
	}
	if err := e.Renderer.Render(rec, "listing_card.html", data, c); err != nil {
		t.Fatalf("Failed to render listing_card.html: %v", err)
	}

	cardBody := rec.Body.String()
	if !strings.Contains(cardBody, job.Company) {
		t.Errorf("Job card missing Company name: %s", job.Company)
	}
	if !strings.Contains(cardBody, job.PayRange) {
		t.Errorf("Job card missing Pay Range: %s", job.PayRange)
	}
	// Verify no emoji in type badge (we'll fix this in templates)
	if strings.Contains(cardBody, "💼") {
		t.Error("Job card UI contains emoji in badge")
	}

	// 2. Verify Detail Modal Rendering
	rec = httptest.NewRecorder()
	if err := e.Renderer.Render(rec, "modal_detail.html", data, c); err != nil {
		t.Fatalf("Failed to render modal_detail.html: %v", err)
	}

	detailBody := rec.Body.String()
	if !strings.Contains(detailBody, job.Company) {
		t.Errorf("Job detail missing Company name: %s", job.Company)
	}
	if !strings.Contains(detailBody, job.PayRange) {
		t.Errorf("Job detail missing Pay Range: %s", job.PayRange)
	}
	if !strings.Contains(detailBody, "Go") || !strings.Contains(detailBody, "Docker") {
		t.Error("Job detail missing skills tags")
	}
}

