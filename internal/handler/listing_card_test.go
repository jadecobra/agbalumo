package handler_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// TestListingCardRendering verifies the logic within listing_card.html
func TestListingCardRendering(t *testing.T) {
	// Parse the template
	tmpl := template.New("listing_card.html").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, _ := values[i].(string)
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"isNew": func(createdAt time.Time) bool {
			if createdAt.IsZero() {
				return false
			}
			return time.Since(createdAt) < 7*24*time.Hour
		},
	})

	// Load the actual template content
	// In a real execution, we rely on the file system, but for this specific "reproduction" test
	// we will read the file in the test or assume we can parse it.
	// However, `view_file` gave us the path: ui/templates/partials/listing_card.html
	// We will use ParseFiles to load it directly.
	_, err := tmpl.ParseFiles("../../ui/templates/partials/listing_card.html")
	if err != nil {
		t.Fatalf("Failed to parse listing_card.html: %v", err)
	}

	tests := []struct {
		name           string
		listing        domain.Listing
		expectedHref   string
		expectedAction string // "WhatsApp", "Call", "Visit Website", etc.
	}{
		{
			name: "WhatsApp and Phone available -> Prefer WhatsApp",
			listing: domain.Listing{
				ID:              "123",
				Title:           "Test Biz",
				ContactWhatsApp: "1234567890",
				ContactPhone:    "0987654321",
			},
			expectedHref:   "https://wa.me/1234567890",
			expectedAction: "WhatsApp",
		},
		{
			name: "Only Phone available -> Prefer Phone",
			listing: domain.Listing{
				ID:           "456",
				Title:        "Phone Biz",
				ContactPhone: "0987654321",
			},
			expectedHref:   "tel:0987654321",
			expectedAction: "Call",
		},
		{
			name: "Only Website available -> Prefer Website",
			listing: domain.Listing{
				ID:         "789",
				Title:      "Web Biz",
				WebsiteURL: "https://example.com",
			},
			expectedHref:   "https://example.com",
			expectedAction: "Visit Website", // or just "Contact" depending on implementation
		},
		{
			name: "Only Email available -> Prefer Email",
			listing: domain.Listing{
				ID:           "101",
				Title:        "Email Biz",
				ContactEmail: "test@example.com",
			},
			expectedHref:   "mailto:test@example.com",
			expectedAction: "Email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			data := map[string]interface{}{
				"Listing": tt.listing,
				"User":    nil,
			}
			if err := tmpl.ExecuteTemplate(&buf, "listing_card", data); err != nil {
				t.Fatalf("Failed to render template: %v", err)
			}

			html := buf.String()

			// 1. Check Container does NOT have HTMX (it should be on the overlay)
			if strings.Contains(html, `<div class="flex flex-col rounded-2xl bg-white dark:bg-surface-dark shadow-soft overflow-hidden group animate-in fade-in zoom-in-95 duration-300 cursor-pointer" hx-get="/listings/`) {
				t.Error("Container still has hx-get, it should be removed")
			}

			// 2. Check Overlay Link Existence
			// Note: The newlines/formatting in template might affect exact string matching,
			// but we can check for key unique substrings or clean up whitespace if needed.
			// Let's check for the key parts near each other if exact match fails, or simplistic contains.
			if !strings.Contains(html, `absolute inset-0 z-10 cursor-pointer`) {
				t.Error("Overlay link div missing class 'absolute inset-0 z-10 cursor-pointer'")
			}
			if !strings.Contains(html, `hx-get="/listings/`+tt.listing.ID+`"`) {
				t.Error("Overlay link div missing hx-get attribute")
			}

			// 3. Check Contact Button Layering
			// The button wrapper should have z-20 and pointer-events-auto
			if !strings.Contains(html, `class="relative z-20 pointer-events-auto"`) {
				t.Error("Contact button wrapper missing 'relative z-20 pointer-events-auto' class for layering")
			}

			// 4. Check Contact Button logic (href) matches expectations
			if !strings.Contains(html, `href="`+tt.expectedHref+`"`) {
				t.Errorf("Expected contact button to have href='%s', got html:\n%s", tt.expectedHref, html)
			}
		})
	}
}
