package main

import (
	"testing"
)

func TestCheckFlyConfig(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		missing        bool
		expectedResult bool
	}{
		{
			name:           "Valid Config",
			content:        "app = 'agbalumo'\nprimary_region = 'dfw'",
			missing:        false,
			expectedResult: true,
		},
		{
			name:           "Secret Leak",
			content:        "[env]\nMY_SECRET = 'supersecret'",
			missing:        false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkFlyConfig(tt.content, tt.missing)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestContainsSensitive(t *testing.T) {
	tests := []struct {
		content  string
		key      string
		expected bool
	}{
		{"SECRET_KEY=123", "SECRET", true},
		{"APP_NAME=test", "SECRET", false},
		{"# A comment about SECRET", "SECRET", false},
	}

	for _, tt := range tests {
		if got := containsSensitive(tt.content, tt.key); got != tt.expected {
			t.Errorf("containsSensitive(%q, %q) = %v", tt.content, tt.key, got)
		}
	}
}
