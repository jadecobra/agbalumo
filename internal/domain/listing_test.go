package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHoursOfOperationField(t *testing.T) {
	t.Parallel()
	// This test ensures the field exists and can be set.
	l := Listing{
		ID:               "test-hours",
		OwnerOrigin:      "Togo",
		Type:             Business,
		Title:            "Hours Test",
		ContactEmail:     "hours@example.com",
		Address:          "Main St",
		HoursOfOperation: "Mon-Fri 9-5",
		CreatedAt:        time.Now(),
	}

	assert.Equal(t, "Mon-Fri 9-5", l.HoursOfOperation)
}

func TestListing_MetaDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		listing  *Listing
		expected string
	}{
		{
			name: "substantive description preserved",
			listing: &Listing{
				Title:       "Aria Suya Kitchen",
				City:        "Arlington",
				Description: "Finest smoked suya and jollof in Arlington.",
				Rating:      4.4,
				ReviewCount: 440,
			},
			expected: "Finest smoked suya and jollof in Arlington.",
		},
		{
			name: "missing description with city rating and reviews",
			listing: &Listing{
				Title:       "Asafo Market",
				City:        "Grand Prairie",
				Description: "",
				Rating:      4.4,
				ReviewCount: 440,
			},
			expected: "Grand Prairie · ★ 4.4 from 440 Google reviews · find it on agbalumo",
		},
		{
			name: "short description overridden with data",
			listing: &Listing{
				Title:       "Asafo Market",
				City:        "Grand Prairie",
				Description: "Restaurant",
				Rating:      4.4,
				ReviewCount: 440,
			},
			expected: "Grand Prairie · ★ 4.4 from 440 Google reviews · find it on agbalumo",
		},
		{
			name: "singular review count",
			listing: &Listing{
				Title:       "Mama Put",
				City:        "Dallas",
				Description: "",
				Rating:      5.0,
				ReviewCount: 1,
			},
			expected: "Dallas · ★ 5.0 from 1 Google review · find it on agbalumo",
		},
		{
			name: "unrated listing with city",
			listing: &Listing{
				Title:       "Asafo Market",
				City:        "Grand Prairie",
				Description: "",
				Rating:      0.0,
				ReviewCount: 0,
			},
			expected: "Grand Prairie · find it on agbalumo",
		},
		{
			name: "rating without review count",
			listing: &Listing{
				Title:       "Lola's Kitchen",
				City:        "Grand Prairie",
				Description: "",
				Rating:      4.5,
				ReviewCount: 0,
			},
			expected: "Grand Prairie · ★ 4.5 · find it on agbalumo",
		},
		{
			name: "no city falls back to title",
			listing: &Listing{
				Title:       "Adom Market",
				City:        "",
				Description: "",
				Rating:      4.4,
				ReviewCount: 440,
			},
			expected: "Adom Market · ★ 4.4 from 440 Google reviews · find it on agbalumo",
		},
		{
			name: "no city unrated falls back to title",
			listing: &Listing{
				Title:       "MaxiMomkitchen",
				City:        "",
				Description: "",
				Rating:      0.0,
				ReviewCount: 0,
			},
			expected: "MaxiMomkitchen · find it on agbalumo",
		},
		{
			name: "whitespace only description treated as missing",
			listing: &Listing{
				Title:       "Suya Spot",
				City:        "Plano",
				Description: "   \n\t  ",
				Rating:      0.0,
				ReviewCount: 0,
			},
			expected: "Plano · find it on agbalumo",
		},
		{
			name:     "nil listing fallback",
			listing:  nil,
			expected: "Find authentic African food in under 60 seconds on agbalumo.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.listing.MetaDescription())
		})
	}
}
