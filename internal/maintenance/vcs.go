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
