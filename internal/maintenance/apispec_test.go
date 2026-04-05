package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareRoutes(t *testing.T) {
	source := []Route{{Method: "GET", Path: "/v1"}}
	target := []Route{{Method: "POST", Path: "/v1"}}

	diffs := CompareRoutes("Code", "OpenAPI", source, target)

	if len(diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(diffs))
	}
}

func TestExtractCLICodeCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli_code_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	goCode := `
package test
var cmd = &cobra.Command{ Use: "test-cmd" }
`
	_ = os.WriteFile(/*nolint:gosec*/ filepath.Join(tmpDir, "cli.go"), []byte(goCode), 0600)

	cmds, err := ExtractCLICodeCommands(tmpDir)
	if err != nil {
		t.Fatalf("ExtractCLICodeCommands failed: %v", err)
	}

	if len(cmds) != 1 || cmds[0] != "test-cmd" {
		t.Errorf("expected [test-cmd], got %v", cmds)
	}
}

func TestExtractCLIMarkdownCommands(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli_md_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	mdCode := `
# CLI Docs
### category
##### add
`
	mdPath := filepath.Join(tmpDir, "cli.md")
	_ = os.WriteFile(/*nolint:gosec*/ mdPath, []byte(mdCode), 0600)

	cmds, err := ExtractCLIMarkdownCommands(mdPath)
	if err != nil {
		t.Fatalf("ExtractCLIMarkdownCommands failed: %v", err)
	}

	expected := map[string]bool{"category": true, "add": true}
	for _, c := range cmds {
		delete(expected, c)
	}

	if len(expected) > 0 {
		t.Errorf("missing expected commands from MD: %v", expected)
	}
}
