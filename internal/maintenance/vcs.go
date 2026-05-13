package maintenance

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckIgnoredFiles checks for ignored files that are accidentally staged for commit.
func CheckIgnoredFiles(rootDir string) error {
	fmt.Println("🔍 Checking for ignored files staged for commit...")

	// Get staged files (Added, Copied, Modified, Renamed)
	cmdDiff := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR")
	cmdDiff.Dir = rootDir
	staged, err := cmdDiff.Output()
	if err != nil {
		return fmt.Errorf("failed to get staged files: %w", err)
	}

	if len(staged) == 0 {
		return nil
	}

	// Filter with git check-ignore
	cmdCheck := exec.Command("git", "check-ignore", "--no-index", "--stdin")
	cmdCheck.Dir = rootDir
	cmdCheck.Stdin = strings.NewReader(string(staged))
	ignored, _ := cmdCheck.Output() // git check-ignore returns exit code 1 if no files match

	if len(ignored) > 0 {
		fmt.Println("❌ Error: The following ignored files are staged for commit:")
		lines := strings.Split(strings.TrimSpace(string(ignored)), "\n")
		for _, line := range lines {
			fmt.Printf("    %s\n", line)
		}
		return fmt.Errorf("ignored files staged for commit")
	}

	fmt.Println("✅ No ignored files found in stage.")
	return nil
}

// CheckGitleaks runs gitleaks on staged files to detect potential secrets.
func CheckGitleaks(rootDir string) error {
	fmt.Println("🛡️  Running gitleaks secrets scan on staged files...")

	// Check if gitleaks is installed
	if _, err := exec.LookPath("gitleaks"); err != nil {
		fmt.Println("⚠️  Warning: 'gitleaks' is not installed. Skipping secrets scan.")
		return nil
	}

	cmd := exec.Command("gitleaks", "protect", "--staged", "--verbose", "--redact")
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Secrets detected by gitleaks:\n%s\n", string(out))
		return fmt.Errorf("gitleaks secrets scan failed")
	}

	fmt.Println("✅ Gitleaks: No secrets detected in staged files.")
	return nil
}

// VerifyGitClean ensures that there are no uncommitted changes in the repository.
func VerifyGitClean(rootDir string) error {
	fmt.Println("🔍 Verifying repository cleanliness (git status)...")

	// --porcelain=v1 ensures the output is easy to parse and consistent across versions.
	// We want to detect both modified (staged/unstaged) and untracked files.
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = rootDir
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run git status: %w", err)
	}

	if len(out) > 0 {
		dirty := filterDirtyLines(string(out))
		if len(dirty) > 0 {
			fmt.Println("❌ Error: Repository has unstaged or untracked changes:")
			for _, d := range dirty {
				fmt.Printf("  %s\n", d)
			}
			return fmt.Errorf("working tree is dirty (run 'git add' or 'git restore')")
		}
	}

	fmt.Println("✅ Working tree is clean (all changes staged or none exist).")
	return nil
}

func filterDirtyLines(statusOutput string) []string {
	lines := strings.Split(strings.TrimSpace(statusOutput), "\n")
	var dirty []string
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		workTree := line[1]
		index := line[0]

		// workTree != ' ' means unstaged change
		// index == '?' means untracked file
		if workTree != ' ' || index == '?' {
			dirty = append(dirty, line)
		}
	}
	return dirty
}

// GetLastCommitMessage returns the message of the most recent commit.
func GetLastCommitMessage(rootDir string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = rootDir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit message: %w", err)
	}
	return string(out), nil
}
