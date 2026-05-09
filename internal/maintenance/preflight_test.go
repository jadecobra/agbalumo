package maintenance

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/testutil"
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
	createVerifyManifest(t, tempDir)

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

func createVerifyManifest(t *testing.T, dir string) {
	content := `
commands:
  - name: test-cmd
    trigger: test_authoring
    description: A test command

skills:
  - name: go-tdd
    trigger: test_authoring
    path: .agents/skills/go-tdd/SKILL.md
`
	if errDir := os.MkdirAll(filepath.Join(dir, ".agents"), 0750); errDir != nil {
		t.Fatal(errDir)
	}
	if errFile := os.WriteFile(filepath.Join(dir, ".agents/verify-manifest.yaml"), []byte(content), 0600); errFile != nil {
		t.Fatal(errFile)
	}
}

func validatePreflightOutput(t *testing.T, r io.Reader) {
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	checks := []string{
		"domain", "ui", "testing", "port: 8443", "* UI lesson 1", "* Test lesson 1",
		"Relevant Skills:", "go-tdd → .agents/skills/go-tdd/SKILL.md",
		"Relevant Verify Commands:", "test-cmd (A test command)",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected %q in output, but not found\nOutput: %s", check, output)
		}
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	cmd := testutil.IsolatedGitCommand(dir, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
	}
}

func TestFilterLessonsByTrigger(t *testing.T) {
	input := `* **Lesson 1** [TRIGGER: handler_change]: Description 1.
* **Lesson 2** [TRIGGER: security_change]: Description 2.
* **Lesson 3** [TRIGGER: template_change]: Description 3.
* **Lesson 4**: Untagged lesson.`

	tests := []struct {
		name            string
		matchedTriggers map[string]bool
		expected        []string
	}{
		{
			name:            "match handler_change",
			matchedTriggers: map[string]bool{"handler_change": true},
			expected:        []string{"Lesson 1", "Lesson 4"},
		},
		{
			name:            "match template and security",
			matchedTriggers: map[string]bool{"template_change": true, "security_change": true},
			expected:        []string{"Lesson 2", "Lesson 3", "Lesson 4"},
		},
		{
			name:            "no matches",
			matchedTriggers: map[string]bool{"something_else": true},
			expected:        []string{"Lesson 4"},
		},
		{
			name:            "empty matches",
			matchedTriggers: map[string]bool{},
			expected:        []string{"Lesson 4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterLessonsByTrigger(input, tt.matchedTriggers)
			verifyFilteredOutput(t, got, tt.expected)
		})
	}
}

func verifyFilteredOutput(t *testing.T, got string, expected []string) {
	t.Helper()
	for _, exp := range expected {
		if !strings.Contains(got, exp) {
			t.Errorf("expected %q in output, but not found\nOutput: %s", exp, got)
		}
	}
	// Verify count
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != len(expected) && !(len(lines) == 1 && lines[0] == "" && len(expected) == 0) {
		t.Errorf("expected %d lines, got %d\nOutput: %s", len(expected), len(lines), got)
	}
}
