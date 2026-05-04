package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePrunedTree_Unlimited(t *testing.T) {
	tmpDir := t.TempDir()
	setupDummyStructure(t, tmpDir)

	got := GeneratePrunedTree(tmpDir, 10)
	want := []string{"file1.txt", "file2.txt", "file3.txt"}
	notWant := []string{"node_modules", ".git"}

	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("missing %q\nGot:\n%s", w, got)
		}
	}
	for _, nw := range notWant {
		if strings.Contains(got, nw) {
			t.Errorf("should not contain %q\nGot:\n%s", nw, got)
		}
	}
}

func TestGeneratePrunedTree_Depth(t *testing.T) {
	tmpDir := t.TempDir()
	setupDummyStructure(t, tmpDir)

	got := GeneratePrunedTree(tmpDir, 1)
	if !strings.Contains(got, "file1.txt") || !strings.Contains(got, "dir1/") {
		t.Errorf("missing expected files at depth 1\nGot:\n%s", got)
	}
	if strings.Contains(got, "file2.txt") {
		t.Errorf("should not contain file2.txt at depth 1\nGot:\n%s", got)
	}
}

func setupDummyStructure(t *testing.T, tmpDir string) {
	t.Helper()
	createFile(t, filepath.Join(tmpDir, "file1.txt"))
	createFile(t, filepath.Join(tmpDir, "dir1", "file2.txt"))
	createFile(t, filepath.Join(tmpDir, "dir1", "dir2", "file3.txt"))
	createFile(t, filepath.Join(tmpDir, "node_modules", "skipped.txt"))
	createFile(t, filepath.Join(tmpDir, ".git", "config"))
}

func createFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
}
