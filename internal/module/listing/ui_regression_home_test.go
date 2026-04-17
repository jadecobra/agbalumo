package listing_test

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHomePageUIValues(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplate(t)}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	testutil.SaveTestListing(t, env.App.DB, "1", "African Food A", func(l *domain.Listing) { l.Type = domain.Food; l.CreatedAt = time.Now() })
	testutil.SaveTestListing(t, env.App.DB, "2", "African Food B", func(l *domain.Listing) { l.Type = domain.Food; l.CreatedAt = time.Now().Add(time.Second) })

	c, rec := testutil.SetupModuleContext(http.MethodGet, "/", nil)
	c.Echo().Renderer = e.Renderer

	h := listmod.NewListingHandler(env.App)
	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	body := rec.Body.String()
	assert.Contains(t, body, "African Food A")
	assert.Contains(t, body, "African Food B")
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
