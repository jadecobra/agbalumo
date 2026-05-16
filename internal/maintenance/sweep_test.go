package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSweep(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestEnvironment(t, tmpDir)

	t.Run("All Pass (Partial Mock)", func(t *testing.T) {
		runAllPassTest(t, tmpDir)
	})

	t.Run("One Fail (Doc Drift)", func(t *testing.T) {
		runDocDriftFailTest(t, tmpDir)
	})

	t.Run("Mixed (Deprecated Warn)", func(t *testing.T) {
		runDeprecatedWarnTest(t, tmpDir)
	})
}

func setupTestEnvironment(t *testing.T, tmpDir string) {
	dirs := []string{
		"docs",
		".agents/skills",
		".agents/workflows",
		"cmd/verify",
		"internal/maintenance",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0750); err != nil {
			t.Fatal(err)
		}
	}

	manifestContent := "commands:\n  - name: test\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".agents/verify-manifest.yaml"), []byte(manifestContent), 0600); err != nil {
		t.Fatal(err)
	}

	resolverContent := "# Resolver\n| Trigger | Skill |\n|---------|-------|\n| test | `.agents/skills/test.md` |\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".agents/skills/RESOLVER.md"), []byte(resolverContent), 0600); err != nil {
		t.Fatal(err)
	}

	skillContent := "---\nname: test\ntriggers: [test]\nmutating: false\n---\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".agents/skills/test.md"), []byte(skillContent), 0600); err != nil {
		t.Fatal(err)
	}
}

func runAllPassTest(t *testing.T, tmpDir string) {
	results, err := RunSweep(tmpDir)
	if err != nil {
		t.Fatalf("RunSweep failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected results, got none")
	}

	var foundBuild bool
	for _, r := range results {
		if r.Gate == "build" {
			foundBuild = true
			if r.Status != "FAIL" {
				t.Errorf("Expected build to fail in empty temp dir, got %s", r.Status)
			}
		}
	}
	if !foundBuild {
		t.Error("build gate not found in results")
	}
}

func runDocDriftFailTest(t *testing.T, tmpDir string) {
	docContent := "Reference to `internal/non-existent-file.go`"
	if err := os.WriteFile(filepath.Join(tmpDir, "docs/broken.md"), []byte(docContent), 0600); err != nil {
		t.Fatal(err)
	}

	results, err := RunSweep(tmpDir)
	if err != nil {
		t.Fatalf("RunSweep failed: %v", err)
	}

	var foundDrift bool
	for _, r := range results {
		if r.Gate == "doc-drift" {
			foundDrift = true
			if r.Status != "FAIL" {
				t.Errorf("Expected doc-drift to fail, got %s", r.Status)
			}
		}
	}
	if !foundDrift {
		t.Error("doc-drift gate not found in results")
	}
}

func runDeprecatedWarnTest(t *testing.T, tmpDir string) {
	goContent := "package main\nvar _ = map[string]interface{}{}\n"
	modDir := filepath.Join(tmpDir, "internal/module")
	_ = os.MkdirAll(modDir, 0750)
	if err := os.WriteFile(filepath.Join(modDir, "test.go"), []byte(goContent), 0600); err != nil {
		t.Fatal(err)
	}

	results, err := RunSweep(tmpDir)
	if err != nil {
		t.Fatalf("RunSweep failed: %v", err)
	}

	var foundDepr bool
	for _, r := range results {
		if r.Gate == "deprecated" {
			foundDepr = true
			if r.Status != "WARN" {
				t.Errorf("Expected deprecated to warn, got %s", r.Status)
			}
		}
	}
	if !foundDepr {
		t.Error("deprecated gate not found in results")
	}
}
