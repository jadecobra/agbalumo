package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T, prefix string) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tmpDir, func() { _ = os.RemoveAll(tmpDir) }
}

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	return path
}

func assertStringsMatch(t *testing.T, msg string, actual []string, expected map[string]bool) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Errorf("%s: expected %d items, got %d: %v", msg, len(expected), len(actual), actual)
		return
	}

	for _, a := range actual {
		if !expected[a] {
			t.Errorf("%s: unexpected item found: %s", msg, a)
		}
	}
}
