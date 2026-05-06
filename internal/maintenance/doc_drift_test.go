package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDocDrift(t *testing.T) {
	tmpDir := t.TempDir()
	setupMockFilesystem(t, tmpDir)
	setupMockDocs(t, tmpDir)

	t.Run("detects drift", func(t *testing.T) {
		violations, err := CheckDocDrift(tmpDir)
		if err != nil {
			t.Fatalf("CheckDocDrift failed: %v", err)
		}

		expectedViolations := 2
		if len(violations) != expectedViolations {
			t.Errorf("expected %d violations, got %d", expectedViolations, len(violations))
		}

		checkSpecificViolations(t, violations)
	})
}

func setupMockFilesystem(t *testing.T, tmpDir string) {
	paths := []string{
		"internal/domain/listing.go",
		"internal/handler/",
		"docs/architecture-overview.md",
		"docs/ATLAS.md",
	}

	for _, p := range paths {
		fullPath := filepath.Join(tmpDir, p)
		if filepath.Ext(p) == "" {
			_ = os.MkdirAll(fullPath, 0700)
		} else {
			_ = os.MkdirAll(filepath.Dir(fullPath), 0700)
			_ = os.WriteFile(fullPath, []byte(""), 0600)
		}
	}
}

func setupMockDocs(t *testing.T, tmpDir string) {
	overviewContent := `
# Architecture Overview
The core logic resides in ` + "`internal/domain/listing.go`" + `.
Handlers are in **internal/handler/**.
A stale reference to ` + "`internal/stale/file.go`" + `.
Another stale one: **internal/missing/dir/**.
`
	_ = os.WriteFile(filepath.Join(tmpDir, "docs/architecture-overview.md"), []byte(overviewContent), 0600)

	atlasContent := `
# ATLAS
Reference to ` + "`internal/domain/listing.go`" + `.
`
	_ = os.WriteFile(filepath.Join(tmpDir, "docs/ATLAS.md"), []byte(atlasContent), 0600)
}

func checkSpecificViolations(t *testing.T, violations []DriftViolation) {
	foundStaleFile := false
	foundStaleDir := false
	for _, v := range violations {
		if v.ReferencedPath == "internal/stale/file.go" {
			foundStaleFile = true
		}
		if v.ReferencedPath == "internal/missing/dir/" {
			foundStaleDir = true
		}
	}

	if !foundStaleFile {
		t.Error("failed to detect stale file: internal/stale/file.go")
	}
	if !foundStaleDir {
		t.Error("failed to detect stale dir: internal/missing/dir/")
	}
}
