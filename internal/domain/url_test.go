package domain_test

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"  ", ""},
		{"example.com", "https://example.com"},
		{"http://test.com", "http://test.com"},
		{"https://secure.com", "https://secure.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.NormalizeURL(tt.input))
		})
	}
}
