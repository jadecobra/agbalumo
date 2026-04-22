package ui

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

func TestModalMenuRendering(t *testing.T) {
	// Initialize renderer with real templates
	// We need to point to the correct relative path for templates
	renderer, err := NewTemplateRenderer("../../ui/templates/*.html", "../../ui/templates/partials/*.html", "../../ui/templates/components/*.html")
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	listing := domain.Listing{
		ID:          "test-123",
		Title:       "Test Restaurant",
		Type:        domain.Food,
		Description: "A great place for food.",
		City:        "Dallas",
		Address:     "123 Main St",
		MenuURL:     "https://example.com/menu.pdf",
		IsActive:    true,
	}

	data := echo.Map{
		"Listing": listing,
	}

	t.Run("prominent_menu_button_exists", func(t *testing.T) {
		var buf bytes.Buffer
		if err := renderer.Render(&buf, "modal_detail", data, c); err != nil {
			t.Fatalf("Failed to render modal_detail: %v", err)
		}

		output := buf.String()
		
		// This specific check will FAIL currently because the code uses 'flex items-center' in the contact section
		// The user wants a "button button-secondary" or similar prominent CTA.
		if !strings.Contains(output, "View Menu") {
			t.Errorf("Expected prominent 'View Menu' text not found")
		}

		if !strings.Contains(output, "cta-block") {
			t.Errorf("CTA block not found in output")
		}
	})
}
