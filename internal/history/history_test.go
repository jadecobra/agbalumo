package history

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStore(t *testing.T) {
	// Setup: Use a temporary directory for tests
	tmpDir := t.TempDir()

	// Override default storage directory for testing
	oldDir := DefaultStorageDir
	DefaultStorageDir = tmpDir
	defer func() { DefaultStorageDir = oldDir }()

	decision := SquadDecision{
		FeatureName:      "test_feature",
		SystemsArchitect: "John Architect",
		ProductOwner:     "Jane Owner",
		SDET:             "Jim Tester",
		BackendEngineer:  "Jill Coder",
		DecisionSummary:  "This is a test decision summary.",
	}

	path, err := Store(decision)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	if path == "" {
		t.Fatal("Store returned empty path")
	}

	// Verify file existence (find the file based on the feature name in tmpDir)
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	var foundFile string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_test_feature.md") {
			foundFile = f.Name()
			break
		}
	}

	if foundFile == "" {
		t.Fatal("Decision file not found")
	}

	// Verify content
	content, err := os.ReadFile(tmpDir + "/" + foundFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedFrontmatter := []string{
		"FeatureName: test_feature",
		"SystemsArchitect: John Architect",
		"ProductOwner: Jane Owner",
		"SDET: Jim Tester",
		"BackendEngineer: Jill Coder",
	}

	for _, s := range expectedFrontmatter {
		if !bytes.Contains(content, []byte(s)) {
			t.Errorf("Expected frontmatter %q not found in %s", s, foundFile)
		}
	}

	if !bytes.Contains(content, []byte("[SUMMARY]")) {
		t.Error("Expected [SUMMARY] header not found")
	}

	if !bytes.Contains(content, []byte(decision.DecisionSummary)) {
		t.Error("Expected decision summary text not found")
	}
}

func TestStore_Errors(t *testing.T) {
	decision := SquadDecision{FeatureName: "err"}

	t.Run("MkdirError", func(t *testing.T) {
		// Use a file as the directory path to trigger MkdirAll failure
		tmpDir := t.TempDir()
		fileAsDir := tmpDir + "/is_a_file"
		_ = os.WriteFile(fileAsDir, []byte("data"), 0644)

		oldDir := DefaultStorageDir
		DefaultStorageDir = fileAsDir + "/subdir"
		defer func() { DefaultStorageDir = oldDir }()

		_, err := Store(decision)
		if err == nil {
			t.Error("expected Store to fail on mkdir error")
		}
	})
	t.Run("WriteError", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Create a directory where the file should be, making WriteFile fail
		dir := tmpDir + "/write_err_dir"
		_ = os.MkdirAll(dir, 0755)

		oldDir := DefaultStorageDir
		DefaultStorageDir = dir
		defer func() { DefaultStorageDir = oldDir }()

		// To trigger SafeWriteFile error, we can make the file path a directory
		// But filename has a timestamp.
		// If we wait and try to create a directory with the SAME name as the file?
		// Hard to predict timestamp.
		// Better: make the directory read-only.
		_ = os.Chmod(dir, 0555) // Read and execute only, no write

		decision := SquadDecision{FeatureName: "werr"}
		_, err := Store(decision)
		if err == nil {
			t.Error("expected Store to fail on write error in read-only directory")
		}
	})
}
