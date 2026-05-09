package testutil

import (
	"os"
	"os/exec"
	"strings"
)

// IsolatedGitCommand creates an exec.Cmd for git that is isolated from the parent's
// git-related environment variables (like GIT_INDEX_FILE, GIT_DIR, etc.)
// This prevents tests from accidentally corrupting or reading the real repository's state
// when run inside a git hook (like pre-commit).
func IsolatedGitCommand(dir string, args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...) //nolint:gosec // intended for test utilities
	cmd.Dir = dir
	
	// Filter out GIT_ environment variables
	var cleanEnv []string
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "GIT_") {
			cleanEnv = append(cleanEnv, env)
		}
	}
	cmd.Env = cleanEnv
	return cmd
}
