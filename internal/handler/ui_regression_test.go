package handler_test

import (
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
	testifyMock "github.com/stretchr/testify/mock"
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

	projectRoot := filepath.Join(wd, "..", "..")
	templatePattern := filepath.Join(projectRoot, "ui", "templates", "*.html")
	partialPattern := filepath.Join(projectRoot, "ui", "templates", "partials", "*.html")

	funcMap := template.FuncMap{
		"mod":   func(i, j int) int { return i % j },
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
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	// Expect calls for Home Page
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	body := rec.Body.String()

	// 1. Verify Listings Container is Relative
	if !strings.Contains(body, `class="relative flex-1`) {
		t.Error("Regression: Listings container missing 'relative' class")
	}

	// 2. Verify Helper Overlay classes
	if !strings.Contains(body, `opacity-0`) || !strings.Contains(body, `pointer-events-none`) {
		t.Error("Regression: Loading indicator missing default transparency/pointer-events classes")
	}

	// 3. Verify HTMX Interaction classes
	if !strings.Contains(body, `[&.htmx-request]:opacity-100`) {
		t.Error("Regression: Loading indicator missing active state classes")
	}

	// 4. Verify Search Indicator
	if !strings.Contains(body, `hx-indicator="#listings-loading"`) {
		t.Error("Regression: Search input missing hx-indicator attribute")
	}
}

func TestFilterUIValues(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	if !strings.Contains(body, `id="filter-container"`) {
		t.Error("Regression: Filter container missing ID 'filter-container'")
	}

	if !strings.Contains(body, `All`) || !strings.Contains(body, `bg-stone-900 text-white`) {
		t.Error("Regression: 'All' button does not appear to have active classes")
	}

	if !strings.Contains(body, `Food (`) {
		t.Error("Regression: 'Food' filter button missing count")
	}
	foodIndex := strings.Index(body, `Food (`)
	if foodIndex != -1 {
		buttonStart := strings.LastIndex(body[:foodIndex], `<button`)
		buttonTag := body[buttonStart:foodIndex]
		if strings.Contains(buttonTag, `bg-stone-900`) {
			t.Error("Regression: 'Food' button seems to have active class by default")
		}
	}

	if !strings.Contains(rec.Body.String(), `src="/static/js/app.js?v=2"`) {
		t.Errorf("Regression: app.js script tag missing")
	}
}

func TestJobListingUI(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := map[string]interface{}{
		"Listing":   job,
		"GridClass": "",
		"User":      domain.User{},
	}
	if err := e.Renderer.Render(rec, "listing_card", data, c); err != nil {
		t.Fatalf("Failed to render listing_card.html: %v", err)
	}

	cardBody := rec.Body.String()
	if !strings.Contains(cardBody, job.Company) {
		t.Errorf("Job card missing Company name: %s", job.Company)
	}
	if !strings.Contains(cardBody, job.PayRange) {
		t.Errorf("Job card missing Pay Range: %s", job.PayRange)
	}
	if strings.Contains(cardBody, "💼") {
		t.Error("Job card UI contains emoji in badge")
	}

	rec = httptest.NewRecorder()
	if err := e.Renderer.Render(rec, "modal_detail", data, c); err != nil {
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
