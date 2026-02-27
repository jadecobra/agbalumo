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
			return false
		},
	})

	_, err := tmpl.ParseFiles("../../ui/templates/partials/listing_card.html")
	if err != nil {
		t.Fatalf("Failed to parse listing_card.html: %v", err)
	}

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

	if !strings.Contains(html, `absolute inset-0 z-10 cursor-pointer`) {
		t.Error("Overlay link div missing class 'absolute inset-0 z-10 cursor-pointer'")
	}
	if !strings.Contains(html, `hx-get="/listings/123"`) {
		t.Error("Overlay link div missing hx-get attribute")
	}
}
