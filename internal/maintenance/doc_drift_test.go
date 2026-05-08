package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDocDrift(t *testing.T) {
	tmpDir := t.TempDir()
	setupMockFilesystem(t, tmpDir)

	// Create a new md file in docs/ to test dynamic scanning
	newDoc := filepath.Join(tmpDir, "docs/new-feature.md")
	_ = os.WriteFile(newDoc, []byte("Reference to `internal/missing.go`"), 0600)

	t.Run("detects drift across all docs", func(t *testing.T) {
		violations, err := CheckDocDrift(tmpDir)
		if err != nil {
			t.Fatalf("CheckDocDrift failed: %v", err)
		}

		// architecture-overview.md has 2, new-feature.md has 1 = 3 total
		expectedViolations := 3
		if len(violations) != expectedViolations {
			t.Errorf("expected %d violations, got %d", expectedViolations, len(violations))
		}
	})
}

func TestCheckCommandDrift(t *testing.T) {
	tmpDir := t.TempDir()
	docPath := filepath.Join(tmpDir, "docs/commands.md")
	_ = os.MkdirAll(filepath.Join(tmpDir, "docs"), 0700)

	content := `
# Commands
Valid: ` + "`go run ./cmd/verify ci`" + `
Stale Task: ` + "`task pre-commit`" + `
Invalid Verify: ` + "`go run ./cmd/verify nonexistent`" + `
`
	_ = os.WriteFile(docPath, []byte(content), 0600)

	t.Run("detects stale commands", func(t *testing.T) {
		violations, err := CheckCommandDrift(tmpDir)
		if err != nil {
			t.Fatalf("CheckCommandDrift failed: %v", err)
		}

		// Should find: task pre-commit, go run ./cmd/verify nonexistent
		if len(violations) != 2 {
			t.Errorf("expected 2 violations, got %d", len(violations))
		}
	})
}

func TestCheckConfigPathDrift(t *testing.T) {
	tmpDir := t.TempDir()
	docPath := filepath.Join(tmpDir, "docs/config.md")
	_ = os.MkdirAll(filepath.Join(tmpDir, "docs"), 0700)
	_ = os.MkdirAll(filepath.Join(tmpDir, ".agents"), 0700)
	_ = os.WriteFile(filepath.Join(tmpDir, ".agents/coverage.json"), []byte("{}"), 0600)

	content := `
# Config
Valid: ` + "`.agents/coverage.json`" + `
Stale: ` + "`.agents/missing.yaml`" + `
`
	_ = os.WriteFile(docPath, []byte(content), 0600)

	t.Run("detects stale config paths", func(t *testing.T) {
		violations, err := CheckConfigPathDrift(tmpDir)
		if err != nil {
			t.Fatalf("CheckConfigPathDrift failed: %v", err)
		}

		if len(violations) != 1 {
			t.Errorf("expected 1 violation, got %d", len(violations))
		}
		if violations[0].ReferencedPath != ".agents/missing.yaml" {
			t.Errorf("expected violation for .agents/missing.yaml, got %s", violations[0].ReferencedPath)
		}
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
			content := ""
			if p == "docs/architecture-overview.md" {
				content = "Reference to `internal/stale/file.go` and **internal/missing/dir/**"
			}
			_ = os.WriteFile(fullPath, []byte(content), 0600)
		}
	}
}

func TestCheckCommandDrift_DynamicFromManifest(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, ".agents"), 0700)
	_ = os.MkdirAll(filepath.Join(tmpDir, "docs"), 0700)

	manifestContent := `
commands:
  - name: my-new-cmd
tools:
  - name: my-tool
`
	_ = os.WriteFile(filepath.Join(tmpDir, ".agents/verify-manifest.yaml"), []byte(manifestContent), 0600)

	docContent := `
# Docs
Valid: ` + "`go run ./cmd/verify my-new-cmd`" + `
Valid: ` + "`go run ./cmd/verify my-tool`" + `
Stale: ` + "`go run ./cmd/verify nonexistent`" + `
`
	_ = os.WriteFile(filepath.Join(tmpDir, "docs/test.md"), []byte(docContent), 0600)

	violations, err := CheckCommandDrift(tmpDir)
	if err != nil {
		t.Fatalf("CheckCommandDrift failed: %v", err)
	}

	// With current hardcoded list, "my-new-cmd" and "my-tool" are missing, so it should find 3 violations.
	// Once implemented, it should only find 1 (nonexistent).
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
		for _, v := range violations {
			t.Logf("Found unexpected violation: %s", v.ReferencedPath)
		}
	}
}
