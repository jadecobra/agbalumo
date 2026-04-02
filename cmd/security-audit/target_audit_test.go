package main

import (
	"testing"
)

func TestIsValidTarget(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://localhost:8443", true},
		{"http://localhost:8080", true},
		{"https://google.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := isValidTarget(tt.url); got != tt.expected {
				t.Errorf("isValidTarget(%q) = %v, want %v", tt.url, got, tt.expected)
			}
		})
	}
}
