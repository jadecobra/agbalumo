package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHomePageUIValues(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplate(t)}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	repo := handler.SetupTestRepository(t)
	err := repo.Save(ctx, domain.Listing{ID: "1", Title: "Business A", Type: domain.Business, IsActive: true, CreatedAt: time.Now()})
	require.NoError(t, err)
	err = repo.Save(ctx, domain.Listing{ID: "2", Title: "Job B", Type: domain.Job, IsActive: true, CreatedAt: time.Now().Add(time.Second)})
	require.NoError(t, err)

	// Verify repo has them
	all, _, err := repo.FindAll(ctx, "", "", "", "", false, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(all))

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	body := rec.Body.String()
	assert.Contains(t, body, "Business A")
	assert.Contains(t, body, "Job B")
	assert.Contains(t, body, "2 listings and growing")
}

func TestTemplateTailwindCleanup(t *testing.T) {
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
		contentBytes, err := os.ReadFile(tmpl)
		if err != nil {
			t.Fatalf("Failed to read template %s: %v", tmpl, err)
		}
		content := string(contentBytes)

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

func TestSearchBarTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "components", "navigation.html")
	templateContent, err := os.ReadFile(templatePath)
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
