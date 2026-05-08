package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckAgentsCoverage(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t, "agents-coverage-")
	defer cleanup()

	// Setup structure:
	// pkg_a/ (AGENTS.md + .go) -> OK
	// pkg_b/ (.go, no AGENTS.md) -> MISSING
	// not_go/ (.txt only) -> IGNORE
	// vendor/pkg/ (.go, no AGENTS.md) -> EXCLUDED

	dirs := []string{"pkg_a", "pkg_b", "not_go", "vendor/pkg"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0750); err != nil {
			t.Fatalf("failed to create dir %s: %v", d, err)
		}
	}

	writeTestFile(t, tmpDir, filepath.Join("pkg_a", "AGENTS.md"), "agents file")
	writeTestFile(t, tmpDir, filepath.Join("pkg_a", "main.go"), "package main")
	writeTestFile(t, tmpDir, filepath.Join("pkg_b", "main.go"), "package main")
	writeTestFile(t, tmpDir, filepath.Join("not_go", "info.txt"), "some text")
	writeTestFile(t, tmpDir, filepath.Join("vendor/pkg", "main.go"), "package main")

	missing, err := CheckAgentsCoverage(tmpDir)
	if err != nil {
		t.Fatalf("CheckAgentsCoverage failed: %v", err)
	}

	if len(missing) != 1 || missing[0] != "pkg_b" {
		t.Errorf("expected [pkg_b], got %v", missing)
	}
}
