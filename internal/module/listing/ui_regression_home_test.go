package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHomePageUIValues(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	saveTestListing(t, app.DB, "1", "Business A", func(l *domain.Listing) { l.CreatedAt = time.Now() })
	saveTestListing(t, app.DB, "2", "Job B", func(l *domain.Listing) { l.Type = domain.Job; l.CreatedAt = time.Now().Add(time.Second) })

	h := listmod.NewListingHandler(app)
	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	body := rec.Body.String()
	assert.Contains(t, body, "Business A")
	assert.Contains(t, body, "Job B")
	assert.Contains(t, body, "2 listings and growing")
}

func TestTemplateTailwindCleanup(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..", "..")

	var templates []string
	_ = filepath.Walk(filepath.Join(projectRoot, "ui", "templates"), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".html") {
			templates = append(templates, path)
		}
		return nil
	})

	for _, tmpl := range templates {
		contentBytes, err := os.ReadFile(filepath.Clean(tmpl))
		if err != nil {
			t.Fatalf("Failed to read template %s: %v", tmpl, err)
		}
		checkTemplateStyles(t, tmpl, string(contentBytes))
	}
}

func checkTemplateStyles(t *testing.T, tmpl string, content string) {
	checks := []struct {
		pattern string
		msg     string
	}{
		{"class=\"[^\"]*gray-", "contains raw 'gray' Tailwind classes. Use 'stone' or 'earth-...' tokens instead."},
		{"class=\"[^\"]*primary", "contains legacy 'primary' class. Use 'earth-accent' instead."},
		{"class=\"[^\"]*orange-", "contains raw 'orange' class. Use 'earth-accent' instead."},
	}

	for _, c := range checks {
		matched, _ := regexp.MatchString(c.pattern, content)
		if matched {
			t.Errorf("Template %s %s", filepath.Base(tmpl), c.msg)
		}
	}
}

func TestSearchBarTheme(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "components", "navigation.html")
	templateContent, err := os.ReadFile(filepath.Clean(templatePath))
	if err != nil {
		t.Fatalf("Failed to read navigation.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-transparent shadow-sm border border-white/20`) {
		t.Error("Search Bar wrapper missing transparent sharp-edged styling")
	}

	if !strings.Contains(content, `text-earth-cream bg-transparent`) {
		t.Error("Search Bar input missing transparent styling")
	}
}
