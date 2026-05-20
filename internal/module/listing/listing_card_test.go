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
		verifyBasicRendering(t, tmpl)
	})

	t.Run("MobileCompactCardCTAs", func(t *testing.T) {
		verifyMobileCompactCardCTAs(t, tmpl)
	})

	t.Run("VerifiedOriginOverride", func(t *testing.T) {
		verifyOriginOverride(t, tmpl)
	})
}

func verifyBasicRendering(t *testing.T, tmpl *template.Template) {
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

	// Layout and Title Position Assertions
	parts := strings.Split(html, "flex-grow min-h-0")
	if len(parts) < 2 {
		t.Fatal("Could not find bottom half signature 'flex-grow min-h-0' in rendered card HTML")
	}
	topHalf := parts[0]
	bottomHalf := parts[1]

	// 1. Check card height constraint
	cardClassStart := strings.Index(topHalf, `class="listing-card`)
	if cardClassStart == -1 {
		t.Fatal("Could not find listing-card class definition in HTML")
	}
	remainingTop := topHalf[cardClassStart+7:]
	cardClassEnd := strings.Index(remainingTop, `"`)
	if cardClassEnd == -1 {
		t.Fatal("Malformed class list quotes in HTML")
	}
	classContent := remainingTop[:cardClassEnd]
	if !strings.Contains(classContent, "h-full") {
		t.Error("Listing card container is missing 'h-full' height constraint class")
	}

	// 2. Check title position: title must not be overlayed in the Top Half
	if strings.Contains(topHalf, "font-serif") {
		t.Error("Title element h3 is still overlayed inside the Image Section (Top Half)")
	}
	if !strings.Contains(bottomHalf, "font-serif") || !strings.Contains(bottomHalf, "Test Biz") {
		t.Error("Title element h3 is missing from the Content Section (Bottom Half)")
	}
}

func verifyMobileCompactCardCTAs(t *testing.T, tmpl *template.Template) {
	t.Run("Food_AllAvailable", func(t *testing.T) {
		verifyFoodAllAvailable(t, tmpl)
	})

	t.Run("Food_NoMenu_HasWebsite", func(t *testing.T) {
		verifyFoodNoMenuHasWebsite(t, tmpl)
	})

	t.Run("Food_NoMenu_NoWebsite_HasPhone", func(t *testing.T) {
		verifyFoodNoMenuNoWebsiteHasPhone(t, tmpl)
	})

	t.Run("Food_NoButtons", func(t *testing.T) {
		verifyFoodNoButtons(t, tmpl)
	})

	t.Run("NonFood_HasWebsite_HasPhone", func(t *testing.T) {
		verifyNonFoodHasWebsiteHasPhone(t, tmpl)
	})

	t.Run("NonFood_NoWebsite_HasPhone", func(t *testing.T) {
		verifyNonFoodNoWebsiteHasPhone(t, tmpl)
	})

	t.Run("NonFood_NoButtons", func(t *testing.T) {
		verifyNonFoodNoButtons(t, tmpl)
	})
}

func verifyFoodAllAvailable(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:           "789",
		Title:        "Suya Grill",
		MenuURL:      "https://suyagrill.com/menu",
		WebsiteURL:   "https://suyagrill.com",
		ContactPhone: "1234567890",
		Type:         domain.Food,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, `data-testid="ag-mobile-view-menu"`) {
		t.Error("Missing 'View Menu' action button")
	}
	if strings.Contains(html, `data-testid="ag-mobile-website"`) {
		t.Error("Website button should not be rendered when View Menu is available")
	}
	if strings.Contains(html, `data-testid="ag-mobile-phone"`) {
		t.Error("Phone button should not be rendered when View Menu is available")
	}
	if strings.Contains(html, `data-testid="ag-mobile-view-details"`) || strings.Contains(html, "View Spot") {
		t.Error("View Spot button must not be rendered")
	}
}

func verifyFoodNoMenuHasWebsite(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:           "789",
		Title:        "Suya Grill",
		WebsiteURL:   "https://suyagrill.com",
		ContactPhone: "1234567890",
		Type:         domain.Food,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, `data-testid="ag-mobile-view-menu"`) {
		t.Error("View Menu button should not be rendered when no MenuURL is available")
	}
	if !strings.Contains(html, `data-testid="ag-mobile-website"`) {
		t.Error("Missing 'Website' action button")
	}
	if strings.Contains(html, `data-testid="ag-mobile-phone"`) {
		t.Error("Phone button should not be rendered when Website is available")
	}
	if strings.Contains(html, `data-testid="ag-mobile-view-details"`) || strings.Contains(html, "View Spot") {
		t.Error("View Spot button must not be rendered")
	}
}

func verifyFoodNoMenuNoWebsiteHasPhone(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:           "789",
		Title:        "Suya Grill",
		ContactPhone: "1234567890",
		Type:         domain.Food,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, `data-testid="ag-mobile-view-menu"`) {
		t.Error("View Menu button should not be rendered")
	}
	if strings.Contains(html, `data-testid="ag-mobile-website"`) {
		t.Error("Website button should not be rendered")
	}
	if !strings.Contains(html, `data-testid="ag-mobile-phone"`) {
		t.Error("Missing 'Phone' action button")
	}
	if strings.Contains(html, `data-testid="ag-mobile-view-details"`) || strings.Contains(html, "View Spot") {
		t.Error("View Spot button must not be rendered")
	}
}

func verifyFoodNoButtons(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:    "789",
		Title: "Suya Grill",
		Type:  domain.Food,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, `data-testid="ag-mobile-view-details"`) {
		t.Error("Expected fallback View Details button to be rendered when no links are present")
	}
}

func verifyNonFoodHasWebsiteHasPhone(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:           "789",
		Title:        "Not Food Biz",
		WebsiteURL:   "https://notfood.com",
		ContactPhone: "1234567890",
		Type:         domain.Job,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, `data-testid="ag-mobile-website"`) {
		t.Error("Missing 'Website' action button")
	}
	if strings.Contains(html, `data-testid="ag-mobile-phone"`) {
		t.Error("Phone button should not be rendered when Website is available")
	}
	if strings.Contains(html, `data-testid="ag-mobile-view-details"`) || strings.Contains(html, "View Spot") {
		t.Error("View Spot button must not be rendered")
	}
}

func verifyNonFoodNoWebsiteHasPhone(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:           "789",
		Title:        "Not Food Biz",
		ContactPhone: "1234567890",
		Type:         domain.Job,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, `data-testid="ag-mobile-website"`) {
		t.Error("Website button should not be rendered")
	}
	if !strings.Contains(html, `data-testid="ag-mobile-phone"`) {
		t.Error("Missing 'Phone' action button")
	}
	if strings.Contains(html, `data-testid="ag-mobile-view-details"`) || strings.Contains(html, "View Spot") {
		t.Error("View Spot button must not be rendered")
	}
}

func verifyNonFoodNoButtons(t *testing.T, tmpl *template.Template) {
	var buf bytes.Buffer
	l := domain.Listing{
		ID:    "789",
		Title: "Not Food Biz",
		Type:  domain.Job,
	}
	data := map[string]interface{}{
		"Listing":   l,
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	}
	if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, `data-testid="ag-mobile-view-details"`) {
		t.Error("Expected fallback View Details button to be rendered when no links are present")
	}
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
