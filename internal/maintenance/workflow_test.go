package maintenance

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestGitRepo(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	runGit := func(args ...string) {
		cmd := exec.Command("git", args...) //nolint:gosec // maintenance utility uses git in tests
		cmd.Dir = tmpDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
		}
	}

	runGit("init")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "Test User")

	// Initial commit
	if err := os.WriteFile( /*nolint:gosec*/ filepath.Join(tmpDir, "README.md"), []byte("# Test Repo"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit("add", "README.md")
	runGit("commit", "-m", "chore: initial commit")

	return tmpDir
}

func TestInferCurrentPhase(t *testing.T) {
	tmpDir := setupTestGitRepo(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	runGit := func(args ...string) {
		cmd := exec.Command("git", args...) //nolint:gosec // maintenance utility uses git in tests
		cmd.Dir = tmpDir
		_ = cmd.Run()
	}

	testIdlePhase(t, tmpDir)
	testRedPhase(t, tmpDir, runGit)
	testGreenPhase(t, tmpDir, runGit)
	testRefactorPhase(t, tmpDir, runGit)
}

func testIdlePhase(t *testing.T, tmpDir string) {
	phase, err := InferCurrentPhase(tmpDir)
	if err != nil || phase != PhaseIdle {
		t.Errorf("IDLE: expected (IDLE, nil), got (%s, %v)", phase, err)
	}
}

func testRedPhase(t *testing.T, tmpDir string, runGit func(...string)) {
	testFile := filepath.Join(tmpDir, "feature_test.go")
	if err := os.WriteFile( /*nolint:gosec*/ testFile, []byte("package main\nfunc TestRed(t *testing.T) {}"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit("add", "feature_test.go")
	phase, err := InferCurrentPhase(tmpDir)
	if err != nil || phase != PhaseRed {
		t.Errorf("RED: expected (RED, nil), got (%s, %v)", phase, err)
	}
}

func testGreenPhase(t *testing.T, tmpDir string, runGit func(...string)) {
	runGit("commit", "-m", "test: add failing test")
	implFile := filepath.Join(tmpDir, "feature.go")
	if err := os.WriteFile( /*nolint:gosec*/ implFile, []byte("package main\nfunc Feature() {}"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit("add", "feature.go")
	phase, err := InferCurrentPhase(tmpDir)
	if err != nil || phase != PhaseGreen {
		t.Errorf("GREEN: expected (GREEN, nil), got (%s, %v)", phase, err)
	}
}

func testRefactorPhase(t *testing.T, tmpDir string, runGit func(...string)) {
	runGit("commit", "-m", "feat: implement feature")
	implFile := filepath.Join(tmpDir, "feature.go")
	newImpl := []byte("package main\nfunc Feature() { /* optimized */ }")
	if err := os.WriteFile( /*nolint:gosec*/ implFile, newImpl, 0600); err != nil {
		t.Fatal(err)
	}
	runGit("add", "feature.go")
	phase, err := InferCurrentPhase(tmpDir)
	if err != nil || phase != PhaseRefactor {
		t.Errorf("REFACTOR: expected (REFACTOR, nil), got (%s, %v)", phase, err)
	}
}
