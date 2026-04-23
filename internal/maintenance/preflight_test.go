package maintenance

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreflight(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "preflight-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	setupPreflightTestRepo(t, tempDir)
	createTestFiles(t, tempDir)
	createAgentsAndStandards(t, tempDir)

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	errRun := RunPreflight(tempDir)
	
	_ = w.Close()
	os.Stdout = oldStdout

	if errRun != nil {
		t.Fatalf("RunPreflight failed: %v", errRun)
	}

	validatePreflightOutput(t, r)
}

func setupPreflightTestRepo(t *testing.T, dir string) {
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func createTestFiles(t *testing.T, tempDir string) {
	// Domain file
	if errDir := os.MkdirAll(filepath.Join(tempDir, "internal/domain"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	if errFile := os.WriteFile(filepath.Join(tempDir, "internal/domain/listing.go"), []byte("package domain"), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	runGit(t, tempDir, "add", ".")
	runGit(t, tempDir, "commit", "-m", "initial commit")

	// Modified/Staged files
	if errFile := os.WriteFile(filepath.Join(tempDir, "internal/domain/listing.go"), []byte("package domain\n// modified"), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	if errDir := os.MkdirAll(filepath.Join(tempDir, "ui"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	if errFile := os.WriteFile(filepath.Join(tempDir, "ui/index.html"), []byte("<html></html>"), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	runGit(t, tempDir, "add", "ui/index.html")

	// Test file
	if errDir := os.MkdirAll(filepath.Join(tempDir, "internal/service"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	if errFile := os.WriteFile(filepath.Join(tempDir, "internal/service/foo_test.go"), []byte("package service"), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	runGit(t, tempDir, "add", "internal/service/foo_test.go")
}

func createAgentsAndStandards(t *testing.T, tempDir string) {
	if errDir := os.MkdirAll(filepath.Join(tempDir, ".agents/workflows"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	content := "### UI & Frontend\n* UI lesson 1\n### Testing\n* Test lesson 1\n"
	if errFile := os.WriteFile(filepath.Join(tempDir, ".agents/workflows/coding-standards.md"), []byte(content), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	if errDir := os.MkdirAll(filepath.Join(tempDir, "internal/handler"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	if errFile := os.WriteFile(filepath.Join(tempDir, "internal/handler/AGENTS.md"), []byte("Handler constraints"), 0600); errFile != nil {
		t.Fatal(errFile)
	}
	if errFile := os.WriteFile(filepath.Join(tempDir, ".agents/invariants.json"), []byte(`{"port": 8443}`), 0600); errFile != nil {
		t.Fatal(errFile)
	}
}

func validatePreflightOutput(t *testing.T, r io.Reader) {
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	checks := []string{"domain", "ui", "testing", "port: 8443", "* UI lesson 1", "* Test lesson 1"}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected %q in output, but not found\nOutput: %s", check, output)
		}
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	// #nosec G204 -- test helper running git commands in temp directory
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
	}
}
