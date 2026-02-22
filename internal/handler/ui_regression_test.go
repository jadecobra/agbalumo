package handler_test

import (
	"encoding/json"
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
		"add":   func(i, j int) int { return i + j },
		"sub":   func(i, j int) int { return i - j },
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
		"toJson": func(v interface{}) (template.JS, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(b), nil
		},
		"isNew": func(createdAt time.Time) bool {
			if createdAt.IsZero() {
				return false
			}
			return time.Since(createdAt) < 7*24*time.Hour
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
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false, 20, 0).Return([]domain.Listing{}, nil)
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
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false, 20, 0).Return([]domain.Listing{}, nil)
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

func TestHomepageFiltersUseFlexWrap(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "index.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	content := string(templateContent)

	idx := strings.Index(content, `id="filter-container"`)
	if idx == -1 {
		t.Fatal("filter-container not found in template")
	}

	filterDiv := content[idx : idx+200]
	if strings.Contains(filterDiv, "overflow-x-auto") {
		t.Error("Homepage filter container should not use overflow-x-auto - use flex-wrap instead")
	}

	if !strings.Contains(filterDiv, "flex-wrap") {
		t.Error("Homepage filter container should use flex-wrap class")
	}
}

func TestAdminFiltersShowCounts(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_listings.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_listings.html: %v", err)
	}

	content := string(templateContent)

	filterSectionIdx := strings.Index(content, "Category Filters")
	if filterSectionIdx == -1 {
		t.Fatal("Category Filters section not found in admin_listings.html")
	}

	filterSection := content[filterSectionIdx : filterSectionIdx+2000]

	if !strings.Contains(filterSection, "Counts") && !strings.Contains(filterSection, "{{ .Counts") {
		t.Error("Admin filters should display listing counts like homepage filters")
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

func TestJobListingCardGradient(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}

	job := domain.Listing{
		ID:          "job-gradient-test",
		Title:       "Test Job",
		Type:        domain.Job,
		Description: "Test description",
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

	if !strings.Contains(cardBody, "from-amber-400 to-orange-500") {
		t.Error("Job card missing expected gradient (from-amber-400 to-orange-500). Check for typo in type comparison.")
	}
}

func TestJustAddedBadgeOnlyShowsForRecentListings(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}

	tests := []struct {
		name        string
		createdAt   time.Time
		expectBadge bool
	}{
		{
			name:        "Recent listing (1 day old) shows badge",
			createdAt:   time.Now().Add(-24 * time.Hour),
			expectBadge: true,
		},
		{
			name:        "Old listing (8 days old) no badge",
			createdAt:   time.Now().Add(-8 * 24 * time.Hour),
			expectBadge: false,
		},
		{
			name:        "Zero CreatedAt (edge case) no badge",
			createdAt:   time.Time{},
			expectBadge: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listing := domain.Listing{
				ID:          "badge-test",
				Title:       "Test Listing",
				Type:        domain.Business,
				Description: "Test description",
				CreatedAt:   tt.createdAt,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			data := map[string]interface{}{
				"Listing":   listing,
				"GridClass": "",
				"User":      domain.User{},
			}
			if err := e.Renderer.Render(rec, "listing_card", data, c); err != nil {
				t.Fatalf("Failed to render listing_card.html: %v", err)
			}

			cardBody := rec.Body.String()
			hasBadge := strings.Contains(cardBody, "Just Added")

			if tt.expectBadge && !hasBadge {
				t.Errorf("Expected 'Just Added' badge for %s, but it was not found", tt.name)
			}
			if !tt.expectBadge && hasBadge {
				t.Errorf("Did not expect 'Just Added' badge for %s, but it was found", tt.name)
			}
		})
	}
}

func TestFilterButtonClassesMatchJS(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	jsPath := filepath.Join(projectRoot, "ui", "static", "js", "app.js")
	jsContent, err := os.ReadFile(jsPath)
	if err != nil {
		t.Fatalf("Failed to read app.js: %v", err)
	}

	templatePath := filepath.Join(projectRoot, "ui", "templates", "index.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	jsStr := string(jsContent)
	templateStr := string(templateContent)

	if !strings.Contains(jsStr, "h-11") {
		t.Error("JS filter button classes should use h-11 to match HTML template")
	}

	if strings.Contains(jsStr, "h-8") && strings.Contains(templateStr, "h-11") {
		t.Error("JS uses h-8 but HTML template uses h-11 - they should match")
	}
}

func TestNoPlaceholderEmailInBaseTemplate(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "base.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read base.html: %v", err)
	}

	if strings.Contains(string(templateContent), "[EMAIL_ADDRESS]") {
		t.Error("base.html contains placeholder [EMAIL_ADDRESS] - replace with actual email or make it configurable")
	}
}

func TestFooterJoinCommunityLinkText(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "base.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read base.html: %v", err)
	}

	content := string(templateContent)

	idx := strings.Index(content, `Join the`)
	if idx == -1 {
		return
	}

	snippet := content[max(0, idx-50):min(len(content), idx+50)]

	if strings.Contains(snippet, "/auth/google/login") {
		if !strings.Contains(snippet, "Sign In") && !strings.Contains(snippet, "Sign Up") {
			t.Error("'Join the Community' link should clearly indicate sign-in action (e.g., 'Sign In to Join Community')")
		}
	}
}

func TestCreateListingFormHasLoadingIndicator(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_create_listing.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_create_listing.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "htmx-request") && !strings.Contains(content, "hx-indicator") {
		t.Error("Create listing form should have loading indicator (htmx-request class or hx-indicator)")
	}
}

func TestAdminDashboardColorConsistency(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_dashboard.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_dashboard.html: %v", err)
	}

	content := string(templateContent)

	if strings.Contains(content, "gray-") {
		t.Error("admin_dashboard.html should use stone-* colors instead of gray-* for consistency with other templates")
	}
}

func TestBaseTemplateHasSkipToContentLink(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "base.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read base.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "skip") && !strings.Contains(content, "sr-only") {
		t.Error("base.html should have a skip-to-content link for accessibility")
	}
}

func TestBaseTemplateHasAriaLiveRegion(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "base.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read base.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "aria-live") {
		t.Error("base.html should have an aria-live region for dynamic content announcements")
	}
}

func TestGoogleMapsApiLoadedLazily(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_create_listing.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_create_listing.html: %v", err)
	}

	content := string(templateContent)

	if strings.Contains(content, "maps.googleapis.com") && !strings.Contains(content, "loadGoogleMaps") {
		t.Error("Google Maps API should be loaded lazily via JS function, not directly in template")
	}
}

func TestListingListHasEmptyState(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "listing_list.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read listing_list.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "No listings") && !strings.Contains(content, "no results") && !strings.Contains(content, "empty") {
		t.Error("listing_list.html should have an empty state message when no listings are found")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
