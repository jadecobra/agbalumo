package maintenance

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/v1/user/{id}", "/api/v1/user/:id"},
		{"/api/v1/user/:id", "/api/v1/user/:id"},
		{"/health", "/health"},
		{"", "/"},
		{"/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizePath(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRouteExtraction_Formats(t *testing.T) {
	runTest := func(name string, content []byte, extract func([]byte) ([]Route, error)) {
		t.Run(name, func(t *testing.T) {
			routes, err := extract(content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(routes) != 2 {
				t.Errorf("expected 2 routes, got %d", len(routes))
			}
		})
	}

	openapi := []byte(`
paths:
  /users:
    get:
      summary: List users
  '/users/{id}':
    post:
      summary: Create user
`)
	runTest("OpenAPI", openapi, ExtractOpenAPIRoutes)

	markdown := []byte(`
| Method | Path |
| --- | --- |
| GET | /api/v1/health |
| POST | /api/v1/login |
`)
	runTest("Markdown", markdown, ExtractMarkdownRoutes)
}

func TestCalculateContextCost(t *testing.T) {
	_, _ = CalculateContextCost(".")
}

func TestExtractRendererFunctions_Empty(t *testing.T) {
	_, _ = ExtractRendererFunctions("non-existent.go")
}
