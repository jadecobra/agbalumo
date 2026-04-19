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
	remaining := make(map[string]bool)
	for k, v := range expected {
		remaining[k] = v
	}

	for _, a := range actual {
		delete(remaining, a)
	}

	if len(remaining) > 0 {
		t.Errorf("%s: missing expected items: %v", msg, remaining)
	}
}
