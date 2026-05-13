package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompileAgentBundle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{"strips html comments", "Line 1\n<!-- comment -->\nLine 2", "Line 1\nLine 2"},
		{"reduces multiple newlines", "Line 1\n\n\nLine 2", "Line 1\nLine 2"},
		{"strips multiple comments", "Line 1\n<!-- c1 -->\n\n<!-- c2 -->\nLine 2", "Line 1\nLine 2"},
		{"complex", "# Title\n<!-- n1 -->\n\nSome.\n\n<!-- n2 -->\n\nEnd.", "# Title\nSome.\nEnd."},
	}

	for _, tt := range tests {
		runMinifyTest(t, tt.name, tt.content, tt.expected)
	}
}

func runMinifyTest(t *testing.T, name, content, expected string) {
	t.Helper()
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.md")
	destFile := filepath.Join(tmpDir, "bundle.min.md")

	if err := os.WriteFile(sourceFile, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	//nolint:gosec // test utility reads from temp directory
	if err := CompileAgentBundle([]string{sourceFile}, destFile); err != nil {
		t.Fatalf("CompileAgentBundle failed: %v", err)
	}

	got, err := os.ReadFile(destFile) //nolint:gosec // test utility reads from temp directory
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}

	if string(got) != expected {
		t.Errorf("%s: got %q, expected %q", name, string(got), expected)
	}
}
