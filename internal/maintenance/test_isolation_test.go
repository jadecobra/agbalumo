package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyTestIsolation(t *testing.T) {
	tempDir := t.TempDir()

	// Valid file using testutil.IsolatedGitCommand
	validContent := `package foo_test
import "github.com/jadecobra/agbalumo/internal/testutil"
func TestValid(t *testing.T) {
	testutil.IsolatedGitCommand("dir", "status")
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "valid_test.go"), []byte(validContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Invalid file using exec.Command("git", ...)
	invalidContent := `package foo_test
import "os/exec"
func TestInvalid(t *testing.T) {
	exec.Command("git", "status")
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "invalid_test.go"), []byte(invalidContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Invalid file using raw os.Setenv
	invalidEnvContent := `package foo_test
import "os"
func TestInvalidEnv(t *testing.T) {
	os.Setenv("FOO", "BAR")
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "invalid_env_test.go"), []byte(invalidEnvContent), 0600); err != nil {
		t.Fatal(err)
	}

	err := VerifyTestIsolation(tempDir)
	if err == nil {
		t.Fatalf("expected error for invalid test isolation, got nil")
	}

	if !strings.Contains(err.Error(), "os.Setenv") {
		t.Fatalf("expected error message to contain 'os.Setenv', got: %v", err.Error())
	}
}
