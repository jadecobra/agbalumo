package handler_test

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
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
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
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

	// Re-parse index.html last to ensure its "content" block takes precedence
	// over other pages (like profile.html) that also define "content"
	indexPath := filepath.Join(projectRoot, "ui", "templates", "index.html")
	_, err = tmpl.ParseFiles(indexPath)
	if err != nil {
		t.Fatalf("Failed to re-parse index.html: %v", err)
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
	mockRepo.On("FindAll", testifyMock.Anything, "", "", "", "", false, 20, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetLocations", testifyMock.Anything).Return([]string{}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	body := rec.Body.String()

	// 1. Verify Listings Container Theme
	if !strings.Contains(body, `bg-earth-dark text-earth-sand`) {
		t.Error("Regression: Listings container missing dark theme classes")
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
	mockRepo.On("FindAll", testifyMock.Anything, "", "", "", "", false, 20, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetLocations", testifyMock.Anything).Return([]string{}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	if !strings.Contains(body, `data-action="toggle-filters"`) {
		t.Error("Regression: Filter button missing data-action toggle-filters")
	}

	if !containsNormalized(body, "All Categories") {
		t.Error("Regression: 'All Categories' option missing")
	}

	if !strings.Contains(rec.Body.String(), `src="/static/js/app.js?v=4"`) {
		t.Errorf("Regression: app.js script tag missing")
	}
}

func TestFilterPanelStructure(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", "", "", false, 20, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{
		"Business": 5,
		"Food":     3,
	}, nil)
	mockRepo.On("GetLocations", testifyMock.Anything).Return([]string{"Lagos", "Accra"}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)
	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	// Filter panel must NOT contain native <select> elements (replaced with inline dropdown list)
	panelIdx := strings.Index(body, `id="filter-dropdown-panel"`)
	if panelIdx == -1 {
		t.Fatal("Filter panel with id='filter-dropdown-panel' not found")
	}
	// Look backward to find the opening <div tag to capture class attributes too
	tagStart := strings.LastIndex(body[:panelIdx], "<div")
	if tagStart == -1 {
		tagStart = panelIdx
	}
	panelSection := body[tagStart:min(len(body), panelIdx+3000)]

	if strings.Contains(panelSection, "<select") {
		t.Error("Filter panel should NOT contain <select> elements — use custom dropdown layout instead")
	}

	// Must have category filter buttons
	if !strings.Contains(panelSection, `data-filter-type="category"`) {
		t.Error("Filter panel missing category filter buttons")
	}

	// Theme: panel should use earth-sand background
	if !strings.Contains(panelSection, "bg-earth-sand") {
		t.Error("Filter panel missing bg-earth-sand theme class")
	}

	// Theme: panel should have ochre accent border
	if !strings.Contains(panelSection, "border-earth-ochre") {
		t.Error("Filter panel missing border-earth-ochre accent")
	}
}

func TestFilterPanelPositioning(t *testing.T) {
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

	// Find the filter panel div
	panelIdx := strings.Index(content, `id="filter-dropdown-panel"`)
	if panelIdx == -1 {
		t.Fatal("Filter panel with id='filter-dropdown-panel' not found in index.html")
	}

	// Get the surrounding div attributes (look backward for the opening tag)
	tagStart := strings.LastIndex(content[:panelIdx], "<div")
	if tagStart == -1 {
		t.Fatal("Could not find opening <div for filter panel")
	}
	panelTag := content[tagStart : panelIdx+100]

	// Panel must have absolute positioning since it's a dropdown relative to search bar
	if !strings.Contains(panelTag, "absolute") {
		t.Error("Filter panel missing 'absolute' positional class")
	}
	if !strings.Contains(panelTag, "right-0") {
		t.Error("Filter panel missing 'right-0' positional class to align below filter button")
	}

	// Panel must have top-full for positioning below
	if !strings.Contains(panelTag, "top-full") {
		t.Error("Filter panel missing 'top-full' positioning class")
	}
}

func TestFilterPanelLocationsRendered(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", "", "", false, 20, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetLocations", testifyMock.Anything).Return([]string{"Lagos", "Accra", "London"}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)
	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	body := rec.Body.String()

	// All mock locations must appear in the rendered HTML
	for _, loc := range []string{"Lagos", "Accra", "London"} {
		if !strings.Contains(body, loc) {
			t.Errorf("Location '%s' not found in rendered filter panel", loc)
		}
	}

	// Must have "All Locations" default option
	if !containsNormalized(body, "All Locations") {
		t.Error("Filter panel missing 'All Locations' default option")
	}

	// Must have "All Categories" default option
	if !containsNormalized(body, "All Categories") {
		t.Error("Filter panel missing 'All Categories' default option")
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

			if tt.expectBadge {
				if !strings.Contains(cardBody, `NEW`) {
					t.Errorf("Expected 'NEW' badge for %s, but it was not found", tt.name)
				}
			} else {
				if strings.Contains(cardBody, `NEW`) {
					t.Errorf("Did not expect 'NEW' badge for %s, but it was found", tt.name)
				}
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
			t.Error("'Join the Community' link should clearly indicate sign-in action (e.g., 'Sign In to Make a Request or Post a Listing')")
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

func TestGoogleFontsLoadInterAndPlayfair(t *testing.T) {
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

	if !strings.Contains(content, "family=Inter") {
		t.Error("base.html should load the Inter font from Google Fonts (defined in Stitch designs)")
	}

	if !strings.Contains(content, "family=Playfair+Display") {
		t.Error("base.html should load the Playfair Display font from Google Fonts (defined in Stitch designs)")
	}
}

func TestProfileTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

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
			b, err := json.Marshal(v)
			return template.JS(b), err
		},
		"isNew": func(createdAt time.Time) bool {
			return false
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl := template.New("base").Funcs(funcMap)
	tmpl, _ = tmpl.ParseFiles(
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "profile.html"),
	)
	tmpl, _ = tmpl.ParseGlob(filepath.Join(projectRoot, "ui", "templates", "partials", "*.html"))

	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: tmpl}

	user := domain.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := map[string]interface{}{
		"User": user,
	}
	if err := e.Renderer.Render(rec, "profile.html", data, c); err != nil {
		t.Fatalf("Failed to render profile.html: %v", err)
	}

	body := rec.Body.String()

	if !strings.Contains(body, "bg-earth-dark") {
		t.Error("profile.html should use bg-earth-dark for the main container background")
	}
	if !strings.Contains(body, "text-earth-cream") {
		t.Error("profile.html should use text-earth-cream for high contrast text")
	}
	if !strings.Contains(body, "font-serif") {
		t.Error("profile.html should use font-serif (Playfair Display) for headings")
	}
}

func TestAboutTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

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
			b, err := json.Marshal(v)
			return template.JS(b), err
		},
		"isNew": func(createdAt time.Time) bool {
			return false
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl := template.New("base").Funcs(funcMap)
	tmpl, _ = tmpl.ParseFiles(
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "about.html"),
	)
	tmpl, _ = tmpl.ParseGlob(filepath.Join(projectRoot, "ui", "templates", "partials", "*.html"))

	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: tmpl}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := e.Renderer.Render(rec, "about.html", nil, c); err != nil {
		t.Fatalf("Failed to render about.html: %v", err)
	}

	body := rec.Body.String()

	// --- Dark-theme class checks (existing) ---
	if !strings.Contains(body, "bg-earth-dark") {
		t.Error("about.html should use bg-earth-dark for the main container background")
	}
	if !strings.Contains(body, "text-white") {
		t.Error("about.html should use text-white for high contrast text to match homepage")
	}
	if !strings.Contains(body, "font-serif font-black italic") {
		t.Error("about.html should use font-serif font-black italic for headings to match homepage")
	}
	if !strings.Contains(body, "bg-earth-sand/10 border border-white/20") {
		t.Error("about.html should use sharp nested containers (bg-earth-sand/10 border border-white/20 p-1)")
	}

	// --- New content section checks (Stitch "Our Story" design) ---

	// Hero section has been removed as per user request

	// Mission section
	if !containsNormalized(body, "find what you want") {
		t.Error("about.html should contain mission heading 'find what you want' section")
	}

	// Agbalumo metaphor section
	if !containsNormalized(body, "find what you want") {
		t.Error("about.html should contain 'find what you want' section")
	}
	if !containsNormalized(body, "Leaf Veins (Our Connections)") {
		t.Error("about.html should contain 'Leaf Veins (Our Connections)' section")
	}
	if !containsNormalized(body, "Juice Drops (The Vibrancy)") {
		t.Error("about.html should contain 'Juice Drops (The Vibrancy)' section")
	}

}

func TestErrorTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "error.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read error.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "bg-earth-dark") {
		t.Error("error.html should use bg-earth-dark for the main background")
	}
	if !strings.Contains(content, "text-earth-cream") {
		t.Error("error.html should use text-earth-cream for high contrast text")
	}
	if !strings.Contains(content, "font-serif") {
		t.Error("error.html should use font-serif (Playfair Display) for the heading")
	}
	if !strings.Contains(content, "bg-earth-accent") {
		t.Error("error.html should use bg-earth-accent for the primary CTA button")
	}
}

func TestNoRawColors(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	var templates []string
	filepath.Walk(filepath.Join(projectRoot, "ui", "templates"), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".html") {
			templates = append(templates, path)
		}
		return nil
	})

	for _, tmpl := range templates {
		contentBytes, err := os.ReadFile(tmpl)
		if err != nil {
			t.Fatalf("Failed to read template %s: %v", tmpl, err)
		}
		content := string(contentBytes)

		// Find any class attributes and check for raw gray/yellow/blue etc (except custom classes if needed)
		if strings.Contains(content, "bg-gray-") || strings.Contains(content, "text-gray-") || strings.Contains(content, "border-gray-") {
			t.Errorf("Template %s contains raw 'gray' Tailwind classes. Use 'stone' or 'earth-...' tokens instead.", filepath.Base(tmpl))
		}
		if strings.Contains(content, "bg-primary") || strings.Contains(content, "text-primary") {
			t.Errorf("Template %s contains legacy 'primary' class. Use 'earth-accent' instead.", filepath.Base(tmpl))
		}
		if strings.Contains(content, "bg-orange-") || strings.Contains(content, "text-orange-") {
			t.Errorf("Template %s contains raw 'orange' class. Use 'earth-accent' instead.", filepath.Base(tmpl))
		}
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

// normalizeWS collapses runs of whitespace to a single space, for HTML text matching.
var wsRE = regexp.MustCompile(`\s+`)

func normalizeWS(s string) string {
	return strings.TrimSpace(wsRE.ReplaceAllString(s, " "))
}

// containsNormalized checks if haystack contains needle after collapsing whitespace.
func containsNormalized(haystack, needle string) bool {
	return strings.Contains(normalizeWS(haystack), normalizeWS(needle))
}

func TestCreateListingModalTheme(t *testing.T) {
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

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl p-6`) {
		t.Error("Create Listing modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`) {
		t.Error("Create Listing modal inputs missing new sharp border wrapper styling")
	}

	if strings.Contains(content, `multiple`) && strings.Contains(content, `type="file"`) {
		t.Error("Regression: Create Listing modal image input should NOT have 'multiple' attribute (single file upload only)")
	}
}

func TestEditListingModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_edit_listing.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_edit_listing.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Edit Listing modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `bg-transparent border-0 border-b border-white/20`) {
		t.Error("Edit Listing modal inputs missing transparent bottom border styling")
	}
}

func TestCreateRequestModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_create_request.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_create_request.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl p-6`) {
		t.Error("Create Request modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`) {
		t.Error("Create Request modal inputs missing sharp border wrapper styling")
	}

	if !strings.Contains(content, `bg-earth-ochre hover:bg-earth-ochre-light`) {
		t.Error("Create Request modal button missing expected earth-ochre styling")
	}
}

func TestDetailModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_detail.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_detail.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Detail modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `<h2 class="text-2xl font-bold font-serif leading-tight">`) {
		t.Error("Detail modal title missing font-serif (Playfair Display) class")
	}
}

func TestProfileModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_profile.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_profile.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Profile modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `<h2 class="text-xl font-bold font-serif text-earth-cream`) {
		t.Error("Profile modal title missing font-serif class")
	}
}

func TestFeedbackModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_feedback.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_feedback.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Feedback modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `border border-white/20 bg-white/5`) {
		t.Error("Feedback modal textarea missing translucent styling")
	}
}

func TestSearchBarTheme(t *testing.T) {
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

	if !strings.Contains(content, `bg-transparent shadow-sm border border-white/20`) {
		t.Error("Search Bar wrapper missing transparent sharp-edged styling")
	}

	if !strings.Contains(content, `text-earth-cream bg-transparent`) {
		t.Error("Search Bar input missing transparent styling")
	}
}

func TestTypographyConstraints(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	// Verify tailwind.config.js enforces Inter and Playfair Display
	tailwindPath := filepath.Join(projectRoot, "tailwind.config.js")
	tailwindContent, err := os.ReadFile(tailwindPath)
	if err != nil {
		t.Fatalf("Failed to read tailwind.config.js: %v", err)
	}
	content := string(tailwindContent)
	if !strings.Contains(content, `"Inter"`) || !strings.Contains(content, `"Playfair Display"`) {
		t.Error("tailwind.config.js does not define Inter and Playfair Display fonts")
	}
}

func TestAdminDashboardTheme(t *testing.T) {
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

	if !strings.Contains(content, `bg-earth-sand`) {
		t.Error("Admin dashboard metrics card missing semantic sand styling bg-earth-sand")
	}

	if !strings.Contains(content, `bg-earth-dark`) {
		t.Error("Admin dashboard page missing dark theme background bg-earth-dark")
	}

	if !strings.Contains(content, `text-earth-dark`) {
		t.Error("Admin dashboard text missing semantic dark text color text-earth-dark")
	}
}

func TestAdminListingsTheme(t *testing.T) {
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

	if !strings.Contains(content, `bg-earth-dark min-h-screen`) {
		t.Error("Admin listings missing expected base dark theme wrapper classes")
	}

	if !strings.Contains(content, `divide-white/10`) {
		t.Error("Admin listings table missing translucent divide styling")
	}
}

func TestAdminUsersTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_users.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_users.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark min-h-screen`) {
		t.Error("Admin users missing expected base dark theme wrapper classes")
	}

	if !strings.Contains(content, `divide-white/10`) {
		t.Error("Admin users table missing translucent divide styling")
	}
}

func TestAdminLoginTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_login.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_login.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark font-sans`) {
		t.Error("Admin login body missing expected dark theme classes")
	}

	if !strings.Contains(content, `border-b border-white/20`) {
		t.Error("Admin login input missing border-bottom styling")
	}
}
