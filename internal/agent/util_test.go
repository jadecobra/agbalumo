package agent

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "root path",
			input:    "/",
			expected: "/",
		},
		{
			name:     "simple path",
			input:    "/users",
			expected: "/users",
		},
		{
			name:     "trailing slash",
			input:    "/users/",
			expected: "/users",
		},
		{
			name:     "duplicate slashes",
			input:    "//users//details",
			expected: "/users/details",
		},
		{
			name:     "parameter replacement :id",
			input:    "/users/:id",
			expected: "/users/{id}",
		},
		{
			name:     "multiple parameters",
			input:    "/orgs/:org_id/users/:user_id/",
			expected: "/orgs/{org_id}/users/{user_id}",
		},
		{
			name:     "empty path returns root",
			input:    "",
			expected: "/",
		},
		{
			name:     "mixed case parameters",
			input:    "/api/:UserId",
			expected: "/api/{UserId}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NormalizePath(tt.input)
			if actual != tt.expected {
				t.Errorf("NormalizePath(%q) = %q; want %q", tt.input, actual, tt.expected)
			}
		})
	}
}
