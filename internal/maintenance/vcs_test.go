package maintenance

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestVerifyGitClean(t *testing.T) {
	// Setup a temporary git repo
	tmpDir, err := os.MkdirTemp("", "git-clean-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	runGit := func(args ...string) {
		//nolint:gosec // maintenance utility runs trusted commands in tests
		cmd := exec.Command("git", args...)
		cmd.Dir = tmpDir
		if runErr := cmd.Run(); runErr != nil {
			t.Fatalf("git %v failed: %v", args, runErr)
		}
	}

	runGit("init")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "Test User")

	// Case 1: Empty repo (clean)
	if err := VerifyGitClean(tmpDir); err != nil {
		t.Errorf("expected clean for empty repo, got error: %v", err)
	}

	// Case 2: Untracked file (dirty)
	//nolint:gosec // test file creation
	if err := os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("test"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := VerifyGitClean(tmpDir); err == nil {
		t.Error("expected dirty for untracked file, got nil")
	}

	// Case 3: Staged file (clean - we allow staged changes in our definition)
	runGit("add", "untracked.txt")
	if err := VerifyGitClean(tmpDir); err != nil {
		t.Errorf("expected clean for staged file, got error: %v", err)
	}

	// Case 4: Unstaged modification (dirty)
	//nolint:gosec // test file creation
	if err := os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("modified"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := VerifyGitClean(tmpDir); err == nil {
		t.Error("expected dirty for unstaged modification, got nil")
	}

	// Case 5: Staged modification (clean)
	runGit("add", "untracked.txt")
	if err := VerifyGitClean(tmpDir); err != nil {
		t.Errorf("expected clean for staged modification, got error: %v", err)
	}
}
