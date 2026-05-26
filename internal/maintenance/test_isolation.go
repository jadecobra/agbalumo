package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyTestIsolation ensures that:
// 1. No *_test.go file uses exec.Command("git", ...) directly (must use testutil.IsolatedGitCommand).
// 2. No *_test.go file uses raw os.Setenv(...) (must use t.Setenv() or clean up properly).
// Replaces Strict Lesson: Test Parallelism & Env Safety
func VerifyTestIsolation(rootDir string) error {
	var violations []string
	execGitRegex := regexp.MustCompile(`exec\.Command(?:Context)?\([^)]*"git"`)
	osSetenvRegex := regexp.MustCompile(`\bos\.Setenv\(`)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return handleDirWalk(path, rootDir, info)
		}

		v := checkFileIsolation(path, info, execGitRegex, osSetenvRegex)
		if v != nil {
			violations = append(violations, *v...)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to scan for test isolation: %w", err)
	}

	if len(violations) > 0 {
		return fmt.Errorf("found test isolation violations:\n  - %s", strings.Join(violations, "\n  - "))
	}

	return nil
}

// handleDirWalk handles skipping of hidden and vendor directories.
func handleDirWalk(path string, rootDir string, info os.FileInfo) error {
	// Skip hidden directories (like .git, .tester, .gemini) and vendor
	// But do not skip the root directory itself!
	if path != rootDir && path != "." && path != "./" {
		if info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
	}
	return nil
}

// checkFileIsolation scans a test file for git command and os.Setenv violations.
func checkFileIsolation(path string, info os.FileInfo, gitRegex, envRegex *regexp.Regexp) *[]string {
	name := info.Name()
	if name == "test_isolation_test.go" || name == "vcs_test.go" {
		return nil
	}

	if filepath.Ext(path) != ".go" || !strings.HasSuffix(name, "_test.go") {
		return nil
	}

	content, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // reading trusted test code
	if err != nil {
		return nil
	}

	var fileViolations []string
	if gitRegex.Match(content) {
		fileViolations = append(fileViolations, fmt.Sprintf("%s: uses exec.Command(\"git\") (use testutil.IsolatedGitCommand instead to prevent env var leakage)", path))
	}
	if envRegex.Match(content) {
		fileViolations = append(fileViolations, fmt.Sprintf("%s: uses os.Setenv (use t.Setenv() instead to prevent global environment variable pollution)", path))
	}

	if len(fileViolations) > 0 {
		return &fileViolations
	}
	return nil
}
