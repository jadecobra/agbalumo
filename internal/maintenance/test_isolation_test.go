package maintenance

import (
	"os"
	"path/filepath"
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

	err := VerifyTestIsolation(tempDir)
	if err == nil {
		t.Fatalf("expected error for invalid test isolation, got nil")
	}

	if err.Error() == "" {
		t.Fatalf("expected error message")
	}
}
