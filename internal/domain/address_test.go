package domain_test

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestExtractCityFromAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "street-city-country",
			address:  "123 Main St, Lagos, Nigeria",
			expected: "Lagos",
		},
		{
			name:     "city-country",
			address:  "Accra, Ghana",
			expected: "Accra",
		},
		{
			name:     "city-only",
			address:  "Nairobi",
			expected: "Nairobi",
		},
		{
			name:     "extra-parts",
			address:  "Suite 101, 456 Park Ave, New York, USA",
			expected: "456 Park Ave",
		},
		{
			name:     "empty",
			address:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.ExtractCityFromAddress(tt.address)
			assert.Equal(t, tt.expected, result)
		})
	}
}
