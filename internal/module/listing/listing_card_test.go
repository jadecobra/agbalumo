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

	// Create renderer to get real CountryRegions data
	renderer, err := ui.NewTemplateRenderer("../../../ui/templates/*.html")
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	tmpl := template.New("listing_card.html").Funcs(renderer.GetFuncMap())

	_, err = tmpl.ParseFiles(
		"../../../ui/templates/partials/listing_card.html",
		"../../../ui/templates/partials/ui_components.html",
		"../../../ui/templates/partials/save_button.html",
	)
	if err != nil {
		t.Fatalf("Failed to parse listing_card.html: %v", err)
	}

	t.Run("BasicRendering", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]interface{}{
			"Listing": domain.Listing{
				ID:    "123",
				Title: "Test Biz",
			},
			"User": nil,
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
		var buf bytes.Buffer
		l := domain.Listing{
			ID:                "456",
			Title:             "Authentic Suya",
			OwnerOrigin:       "Ghana",
			RegionalSpecialty: "Nigeria",
			Type:              domain.Food,
		}

		// 1. Not verified - should show Ghana flag
		data := map[string]interface{}{"Listing": l, "User": nil}
		if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
			t.Fatalf("Render failed: %v", err)
		}
		html := buf.String()
		if !strings.Contains(html, "🇬🇭") {
			t.Error("Should show poster flag (🇬🇭) when not verified")
		}

		// 2. Verified - should show Nigeria flag
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
		if strings.Contains(html, "🇬🇭") {
			t.Error("Should NOT show poster flag (🇬🇭) when verified origin exists")
		}

		// 3. Redundancy check - if RegionalSpecialty is same as Flag Name, hide badge
		// Use a more specific check to avoid matching title attributes
		if strings.Contains(html, "tracking-wide\">Origin: Nigeria</span>") {
			t.Error("Should hide 'Origin: Nigeria' badge if it's already shown via flag")
		}

		// 4. Ethnic specialty check - should show flag AND badge
		buf.Reset()
		l.RegionalSpecialty = "Yoruba"
		data["Listing"] = l
		_ = tmpl.ExecuteTemplate(&buf, "listing_card", data)
		html = buf.String()
		if !strings.Contains(html, "🇳🇬") {
			t.Error("Should show Nigeria flag (🇳🇬) for Yoruba specialty")
		}
		if !strings.Contains(html, "tracking-wide\">Origin: Yoruba</span>") {
			t.Error("Should show 'Origin: Yoruba' badge for ethnic specialty")
		}
	})
}
