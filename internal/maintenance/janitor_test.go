package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

var debrisFiles = []string{
	"critique_full.txt",
	"fail_ci.yml",
	"success_ci.yml",
	"fail_ci_job.log",
	"full_ci_job.log",
	"success_ci_job.log",
	"precommit_check.txt",
	"index_test.html",
	"server.log",
}

func TestJanitor(t *testing.T) {
	// Setup temp directory
	tempDir, err := os.MkdirTemp("", "janitor_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create dummy debris files in the root of tempDir
	for _, d := range debrisFiles {
		path := filepath.Join(tempDir, d)
		if err := os.WriteFile(path, []byte("test data"), 0600); err != nil {
			t.Fatalf("failed to create debris file %s: %v", d, err)
		}
	}

	// Run Janitor
	if err := RunJanitor(tempDir); err != nil {
		t.Fatalf("RunJanitor failed: %v", err)
	}

	verifyJanitorResults(t, tempDir)
}

func verifyJanitorResults(t *testing.T, tempDir string) {
	t.Helper()

	// Verify .tester/ exists
	testerDir := filepath.Join(tempDir, ".tester")
	if _, err := os.Stat(testerDir); os.IsNotExist(err) {
		t.Errorf(".tester directory was not created")
	}

	// Verify files moved to .tester/
	for _, d := range debrisFiles {
		oldPath := filepath.Join(tempDir, d)
		newPath := filepath.Join(testerDir, d)

		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			t.Errorf("original file %s still exists in root", d)
		}

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Errorf("file %s was not moved to .tester/", d)
		}
	}
}
