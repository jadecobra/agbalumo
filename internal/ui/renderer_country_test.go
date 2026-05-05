package ui

import (
	"testing"
)

func TestTemplateRenderer_CountryData(t *testing.T) {
	t.Parallel()

	t.Run("GetCountryFlagLogic", func(t *testing.T) {
		t.Parallel()
		regions := []Region{
			{
				Region: "West Africa",
				Countries: []Country{
					{Name: "Nigeria", Flag: "🇳🇬"},
					{Name: "Ghana", Flag: "🇬🇭"},
				},
			},
		}

		tests := []struct {
			name     string
			expected string
		}{
			{"Nigeria", "🇳🇬"},
			{"ghana", "🇬🇭"},
			{"Unknown", ""},
		}

		for _, tt := range tests {
			actual := getCountryFlag(regions, tt.name)
			if actual != tt.expected {
				t.Errorf("getCountryFlag(%q) = %q, expected %q", tt.name, actual, tt.expected)
			}
		}
	})

	t.Run("CheckCategorizeFiles", func(t *testing.T) {
		t.Parallel()
		layouts, partials, pages := categorizeTemplateFiles([]string{
			"ui/templates/base.html",
			"ui/templates/partials/foo.html",
			"ui/templates/index.html",
		})
		if len(layouts) != 1 || len(partials) != 1 || len(pages) != 1 {
			t.Errorf("Categorization failed: %v, %v, %v", layouts, partials, pages)
		}
	})
}
