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

	yamlContent, err := agent.RunCommand("npx", "-y", "swagger-cli", "bundle", "../../docs/openapi.yaml", "-r", "-t", "yaml")
	if err != nil {
		t.Skipf("Skipping integration test: failed to bundle openapi.yaml (likely environmental): %v\nOutput: %s", err, string(yamlContent))
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

func TestExtractCLICodeCommands(t *testing.T) {
	// Create a temporary mock cmd dir
	tmpDir := t.TempDir()
	code1 := `package main
	import "github.com/spf13/cobra"
	var cmd = &cobra.Command{ Use: "serve" }`
	code2 := `package main
var cmd = &cobra.Command{
	Use:   "admin",
}
var subCmd = &cobra.Command{
	Use: "approve [id]",
}`
	_ = os.WriteFile(tmpDir+"/1.go", []byte(code1), 0644)
	_ = os.WriteFile(tmpDir+"/2.go", []byte(code2), 0644)
	_ = os.WriteFile(tmpDir+"/ignore.txt", []byte(`Use: "ignoreMe"`), 0644) // Should be ignored

	expected := []string{"admin", "approve", "serve"}
	cmds, err := agent.ExtractCLICodeCommands(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestExtractCLIMarkdownCommands(t *testing.T) {
	tmpDir := t.TempDir()
	md1 := `
### Serve
Some description

### Admin
Another description

##### approve
Approves stuff

### Subcommands
(Ignored)
`
	md2 := `
### seed
Description	
`
	_ = os.WriteFile(tmpDir+"/cli.md", []byte(md1), 0644)
	_ = os.MkdirAll(tmpDir+"/cli", 0755)
	_ = os.WriteFile(tmpDir+"/cli/seed.md", []byte(md2), 0644)

	expected := []string{"admin", "approve", "seed", "serve"}

	cmds, err := agent.ExtractCLIMarkdownCommands(tmpDir+"/cli.md", tmpDir+"/cli/seed.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(cmds, expected) {
		t.Errorf("expected %v, got %v", expected, cmds)
	}
}

func TestCheckCLIDrift(t *testing.T) {
	codeCmds := []string{"serve", "admin"}
	mdCmds := []string{"serve", "seed"}

	diffs := agent.CheckCLIDrift(codeCmds, mdCmds)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d: %v", len(diffs), diffs)
	}
	
	expected1 := "❌ Missing in CLI Docs: admin (found in Code)"
	expected2 := "❌ Missing in Code: seed (found in CLI Docs)"

	found1, found2 := false, false
	for _, d := range diffs {
		if d == expected1 { found1 = true }
		if d == expected2 { found2 = true }
	}

	if !found1 || !found2 {
		t.Errorf("did not find expected differences, got: %v", diffs)
	}
}
