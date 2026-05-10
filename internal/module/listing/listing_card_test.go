package listing_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/ui"
)

// TestListingCardRendering verifies the logic within listing_card.html
func TestListingCardRendering(t *testing.T) {
	t.Parallel()
	tmpl := setupListingCardTmpl(t)

	t.Run("BasicRendering", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]interface{}{
			"Listing": domain.Listing{
				ID:    "123",
				Title: "Test Biz",
			},
			"User":      nil,
			"IDPrefix":  "",
			"GridClass": "",
		}
		if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
			t.Fatalf("Failed to render template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `hx-get="/listings/123"`) {
			t.Error("Overlay link div missing hx-get attribute")
		}
	})

	t.Run("VerifiedOriginOverride", func(t *testing.T) {
		verifyOriginOverride(t, tmpl)
	})
}

func setupListingCardTmpl(t *testing.T) *template.Template {
	t.Helper()
	renderer, err := ui.NewTemplateRenderer("../../../ui/templates/*.html")
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	tmpl := template.New("listing_card.html").Funcs(renderer.GetFuncMap())
	_, err = tmpl.ParseFiles(
		"../../../ui/templates/partials/listing_card.html",
		"../../../ui/templates/partials/ui_components.html",
		"../../../ui/templates/partials/save_button.html",
		"../../../ui/templates/partials/rating_stars.html",
	)
	if err != nil {
		t.Fatalf("Failed to parse listing_card.html: %v", err)
	}
	return tmpl
}

func verifyOriginOverride(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:                "456",
		Title:             "Authentic Suya",
		OwnerOrigin:       "Ghana",
		RegionalSpecialty: "Nigeria",
		Type:              domain.Food,
	}

	data := map[string]interface{}{"Listing": l, "User": nil, "IDPrefix": "", "GridClass": ""}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "🇬🇭") {
		t.Error("Should show poster flag (🇬🇭) when not verified")
	}

	// Verified - should show Nigeria flag
	buf.Reset()
	now := time.Now()
	l.EnrichmentAttemptedAt = &now
	data["Listing"] = l
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	html = buf.String()
	if !strings.Contains(html, "🇳🇬") {
		t.Error("Should show verified origin flag (🇳🇬) when enriched")
	}

	// Ethnic specialty check
	buf.Reset()
	l.RegionalSpecialty = "Yoruba"
	data["Listing"] = l
	_ = tmpl.ExecuteTemplate(&buf, "listing_card", data)
	html = buf.String()
	if !strings.Contains(html, "🇳🇬") {
		t.Error("Should show Nigeria flag (🇳🇬) for Yoruba specialty")
	}
}
