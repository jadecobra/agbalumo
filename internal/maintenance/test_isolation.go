package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyTestIsolation ensures that no *_test.go file uses exec.Command("git", ...) directly.
// They must use testutil.IsolatedGitCommand instead to prevent environment variable leakage
// (e.g., GIT_INDEX_FILE, GIT_DIR) from the parent process, which causes index corruption
// when the tests are run from inside a git hook like pre-commit.
func VerifyTestIsolation(rootDir string) error {
	var violations []string
	execGitRegex := regexp.MustCompile(`exec\.Command(?:Context)?\([^)]*"git"`)

	err := filepath.Walk(rootDir, createIsolationWalker(&violations, execGitRegex))
	if err != nil {
		return fmt.Errorf("failed to scan for test isolation: %w", err)
	}

	if len(violations) > 0 {
		return fmt.Errorf("found exec.Command(\"git\") in %d test files (use testutil.IsolatedGitCommand instead to prevent env var leakage):\n  - %s", len(violations), strings.Join(violations, "\n  - "))
	}

	return nil
}

func createIsolationWalker(violations *[]string, regex *regexp.Regexp) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skip, walkErr := shouldSkipIsolationDir(info); skip {
			return walkErr
		}

		if isIsolationTarget(path, info) && isViolatingTestIsolation(path, regex) {
			*violations = append(*violations, path)
		}
		return nil
	}
}

func shouldSkipIsolationDir(info os.FileInfo) (bool, error) {
	if !info.IsDir() {
		return false, nil
	}
	if info.Name() == "vendor" || info.Name() == ".git" {
		return true, filepath.SkipDir
	}
	return true, nil
}

func isIsolationTarget(path string, info os.FileInfo) bool {
	return filepath.Ext(path) == ".go" && strings.HasSuffix(info.Name(), "_test.go")
}

func isViolatingTestIsolation(path string, regex *regexp.Regexp) bool {
	content, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // reading trusted test code
	if err != nil {
		return false
	}
	return regex.Match(content)
}
