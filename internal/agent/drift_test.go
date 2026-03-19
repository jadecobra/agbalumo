package agent_test

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestExtractOpenAPIRoutes(t *testing.T) {
	yamlContent := []byte(`
paths:
  /users:
    get:
      summary: List users
    post:
      summary: Create user
  /users/{id}:
    get:
      summary: Get user
`)

	expected := []agent.Route{
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/users"},
		{Method: "GET", Path: "/users/{id}"},
	}

	routes, err := agent.ExtractOpenAPIRoutes(yamlContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(routes, expected) {
		t.Errorf("expected %v, got %v", expected, routes)
	}
}

func TestExtractMarkdownRoutes(t *testing.T) {
	mdContent := []byte(`
| Method | Path | Description |
|--------|------|-------------|
| GET | ` + "`" + `/auth/dev` + "`" + ` | Dev login (development only) |
| POST | /listings/:id/claim | Claim listing |
`)

	// Our bash script normalized :id to {id}
	expected := []agent.Route{
		{Method: "GET", Path: "/auth/dev"},
		{Method: "POST", Path: "/listings/{id}/claim"},
	}

	routes, err := agent.ExtractMarkdownRoutes(mdContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(routes, expected) {
		t.Errorf("expected %v, got %v", expected, routes)
	}
}

func TestCompareRoutes(t *testing.T) {
	source := []agent.Route{
		{Method: "GET", Path: "/a"},
		{Method: "POST", Path: "/b"},
	}
	target := []agent.Route{
		{Method: "GET", Path: "/a"},
	}

	diffs := agent.CompareRoutes("Source", "Target", source, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diffs), diffs)
	}
	expectedDiff := "❌ Missing in Target: POST /b (found in Source)"
	if diffs[0] != expectedDiff {
		t.Errorf("expected %q, got %q", expectedDiff, diffs[0])
	}
}

func TestIntegrationAPIDrift(t *testing.T) {
	// Skip if we are running unit tests quickly, or just run it since it's local
	codeRoutes, err := agent.ExtractRoutes("../../cmd", "../../internal/handler", "../../internal/module")
	if err != nil {
		t.Fatalf("failed to extract code routes: %v", err)
	}

	yamlContent, err := os.ReadFile("../../docs/openapi.yaml")
	if err != nil {
		t.Fatalf("failed to read openapi.yaml: %v", err)
	}
	openapiRoutes, err := agent.ExtractOpenAPIRoutes(yamlContent)
	if err != nil {
		t.Fatalf("failed to extract openapi routes: %v", err)
	}

	mdContent, err := os.ReadFile("../../docs/api.md")
	if err != nil {
		t.Fatalf("failed to read api.md: %v", err)
	}
	mdRoutes, err := agent.ExtractMarkdownRoutes(mdContent)
	if err != nil {
		t.Fatalf("failed to extract md routes: %v", err)
	}

	diffs := agent.CheckAPIDrift(codeRoutes, openapiRoutes, mdRoutes)
	if len(diffs) > 0 {
		t.Errorf("expected 0 drift, got %d:\n%v", len(diffs), strings.Join(diffs, "\n"))
	}
}
