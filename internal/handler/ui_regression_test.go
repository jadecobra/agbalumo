package handler_test

import (
	"context"
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
	"github.com/stretchr/testify/assert"
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

	componentPattern := filepath.Join(projectRoot, "ui", "templates", "components", "*.html")

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

	tmpl, err := template.New("base").Funcs(funcMap).ParseGlob(templatePattern)
	if err != nil {
		t.Fatalf("Failed to parse base templates: %v", err)
	}
	_, err = tmpl.ParseGlob(partialPattern)
	if err != nil {
		t.Fatalf("Failed to parse partial templates: %v", err)
	}
	_, err = tmpl.ParseGlob(componentPattern)
	if err != nil {
		t.Fatalf("Failed to parse component templates: %v", err)
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
	ctx := context.Background()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", ctx, "", "", "", "", false, 20, 0).Return([]domain.Listing{
		{ID: "1", Title: "Business A", Type: domain.Business, IsActive: true, CreatedAt: time.Now()},
		{ID: "2", Title: "Job B", Type: domain.Job, IsActive: true, CreatedAt: time.Now()},
	}, nil).Maybe()
	mockRepo.On("GetCounts", ctx).Return(map[domain.Category]int{domain.Business: 1, domain.Job: 1}, nil).Maybe()
	mockRepo.On("GetLocations", ctx).Return([]string{"Lagos"}, nil).Maybe()
	mockRepo.On("GetFeaturedListings", ctx).Return([]domain.Listing{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	body := rec.Body.String()
	assert.Contains(t, body, "Business A")
	assert.Contains(t, body, "Job B")
}

func TestTemplateTailwindCleanup(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	var templates []string
	_ = filepath.Walk(filepath.Join(projectRoot, "ui", "templates"), func(path string, info os.FileInfo, err error) error {
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

	if !strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`) {
		t.Error("Edit Listing modal inputs missing new sharp border wrapper styling")
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

	// Bug fix: close button must use data-modal-action="close" (standardized pattern)
	// This ensures mobile touch events reliably reach the button element
	if !strings.Contains(content, `data-modal-action="close"`) {
		t.Error("Profile modal close button must use data-modal-action=\"close\" for reliable mobile touch handling")
	}

	// Backdrop click-to-close is handled globally by modals.js for all dialog[id] elements
	// Verify the dialog has an id so modals.js can attach the listener
	if !strings.Contains(content, `id="profile-modal"`) {
		t.Error("Profile modal <dialog> must have id='profile-modal' for modals.js backdrop-click handling")
	}

	// Bug fix: dialog must NOT be forced full-height (h-full) — it should show the backdrop around it
	if strings.Contains(content, `h-full`) {
		t.Error("Profile modal <dialog> must not use h-full — it should be constrained so the backdrop is visible and clickable")
	}

	// Bug fix: item count badge must be legible on dark background
	if !strings.Contains(content, `bg-white/10`) || !strings.Contains(content, `text-earth-cream`) {
		t.Error("Profile modal item count badge must use bg-white/10 and text-earth-cream for legibility on dark background")
	}

	// Bug fix: Sign Out link must be reachable on mobile (not buried inside an overflowing flex row)
	// Verify the sign-out link exists and is not inside the same flex row as the avatar/name block
	signOutIdx := strings.Index(content, `/auth/logout`)
	if signOutIdx == -1 {
		t.Error("Profile modal Sign Out link (/auth/logout) not found")
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

	componentsDir := filepath.Join(projectRoot, "ui", "templates", "components")
	files, _ := os.ReadDir(componentsDir)
	for _, f := range files {
		compContent, _ := os.ReadFile(filepath.Join(componentsDir, f.Name()))
		content += string(compContent)
	}

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

	if !strings.Contains(content, `bg-earth-dark flex-1`) {
		t.Error("Admin listings missing expected base dark theme wrapper classes (bg-earth-dark flex-1)")
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

// TestModalCloseButtons verifies that every modal's close button
// has the styled CLOSE text (not just an icon), ensuring it is tappable on mobile.
func TestModalCloseButtons(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	type modalCheck struct {
		file     string
		wantText string
		wantAttr string
	}

	checks := []modalCheck{
		{
			file:     "modal_create_listing.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
		{
			file:     "modal_create_request.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
		{
			file:     "modal_profile.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
		{
			file:     "modal_feedback.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
		{
			file:     "modal_detail.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
		{
			file:     "modal_edit_listing.html",
			wantText: "CLOSE",
			wantAttr: `data-modal-action="close"`,
		},
	}

	for _, check := range checks {
		t.Run(check.file, func(t *testing.T) {
			path := filepath.Join(partialsDir, check.file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", check.file, err)
			}
			body := string(content)

			if !strings.Contains(body, check.wantText) {
				t.Errorf("%s: expected CLOSE button with text %q, not found", check.file, check.wantText)
			}
			if !strings.Contains(body, check.wantAttr) {
				t.Errorf("%s: expected close button with attribute %q, not found", check.file, check.wantAttr)
			}
		})
	}
}

// TestModalCloseButtonStyle verifies that modal CLOSE buttons use the ochre style
// matching the ASK / POST buttons on the homepage.
func TestModalCloseButtonStyle(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	modalFiles := []string{
		"modal_create_listing.html",
		"modal_create_request.html",
	}

	for _, file := range modalFiles {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(partialsDir, file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", file, err)
			}
			body := string(content)

			// The CLOSE button should use the ochre colour to match ASK/POST buttons
			if !strings.Contains(body, "bg-earth-ochre") {
				t.Errorf("%s: CLOSE button missing bg-earth-ochre class (should match ASK/POST button style)", file)
			}
		})
	}
}

// TestAdminDashboardModalCloseButtons verifies the admin dashboard div-based
// modals all have accessible CLOSE buttons (already at the bottom via button_sharp).
func TestAdminDashboardModalCloseButtons(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	componentsDir := filepath.Join(projectRoot, "ui", "templates", "components")
	files, err := os.ReadDir(componentsDir)
	if err != nil {
		t.Fatalf("Failed to read components dir: %v", err)
	}

	closeCount := 0
	for _, f := range files {
		content, _ := os.ReadFile(filepath.Join(componentsDir, f.Name()))
		closeCount += strings.Count(string(content), `"Label" "Close"`)
	}

	if closeCount < 4 {
		t.Errorf("admin components should have at least 4 CLOSE bottom buttons (one per modal), found %d", closeCount)
	}
}

// TestModalNoOrphanIconOnlyCloseButton ensures there are no bare icon-only close
// buttons left (tiny <span>close</span> inside a plain <button> with no text).
func TestModalNoOrphanIconOnlyCloseButton(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	// Regex to find buttons that have ONLY an icon span and no text label visible
	// Pattern: button ... > (whitespace) <span class="material-symbols-outlined">close</span> (whitespace) </button>
	iconOnlyPattern := regexp.MustCompile(`(?s)<button[^>]*data-modal-action="close"[^>]*>\s*<span class="material-symbols-outlined[^"]*">close</span>\s*</button>`)

	modalFiles := []string{
		"modal_create_listing.html",
		"modal_create_request.html",
		"modal_profile.html",
		"modal_feedback.html",
		"modal_detail.html",
		"modal_edit_listing.html",
	}

	for _, file := range modalFiles {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(partialsDir, file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", file, err)
			}
			if iconOnlyPattern.Match(content) {
				t.Errorf("%s: found icon-only close button with data-modal-action=\"close\" — replace with CLOSE text label for mobile accessibility", file)
			}
		})
	}
}

func TestAdminListingsUIElements(t *testing.T) {
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

	// 1. Container should match admin_dashboard.html
	if !strings.Contains(content, "pt-32") || !strings.Contains(content, "max-w-6xl") || !strings.Contains(content, "bg-earth-dark") {
		t.Error("Regression: admin_listings.html missing standard admin container classes (pt-32, max-w-6xl, bg-earth-dark)")
	}

	// 2. Headings should use admin_dashboard.html style
	if !strings.Contains(content, "text-[10px]") || !strings.Contains(content, "uppercase") || !strings.Contains(content, "tracking-[0.3em]") {
		t.Error("Regression: admin_listings.html typography missing premium admin styling (text-[10px] uppercase tracking-[0.3em])")
	}

	// 3. Table headers should use bg-white/5 and text-white/50
	if !strings.Contains(content, "bg-white/5") || !strings.Contains(content, "text-white/50") {
		t.Error("Regression: admin_listings.html table headers missing premium dark styling (bg-white/5, text-white/50)")
	}
}
