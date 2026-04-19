package maintenance

import (
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
	tmpDir, cleanup := setupTestDir(t, "cli_code_test")
	defer cleanup()

	goCode := `
package test
var cmd = &cobra.Command{ Use: "test-cmd" }
`
	_ = writeTestFile(t, tmpDir, "cli.go", goCode)

	cmds, err := ExtractCLICodeCommands(tmpDir)
	if err != nil {
		t.Fatalf("ExtractCLICodeCommands failed: %v", err)
	}

	if len(cmds) != 1 || cmds[0] != "test-cmd" {
		t.Errorf("expected [test-cmd], got %v", cmds)
	}
}

func TestExtractCLIMarkdownCommands(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t, "cli_md_test")
	defer cleanup()

	mdCode := `
# CLI Docs
### category
##### add
`
	path := writeTestFile(t, tmpDir, "cli.md", mdCode)
	cmds, err := ExtractCLIMarkdownCommands(path)
	if err != nil {
		t.Fatalf("ExtractCLIMarkdownCommands failed: %v", err)
	}

	assertStringsMatch(t, "CLI commands", cmds, map[string]bool{"category": true, "add": true})
}
