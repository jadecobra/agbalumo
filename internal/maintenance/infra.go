package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyActionSHAs checks for correctly pinned SHAs in .github/workflows and .github/actions.
func VerifyActionSHAs(rootDir string) error {
	fmt.Println("🔍 Verifying GitHub Action SHA pinning...")

	var errorCount int
	// Find all workflow and action files
	wfDir := filepath.Join(rootDir, ".github/workflows")
	actDir := filepath.Join(rootDir, ".github/actions")

	var files []string
	_ = filepath.Walk(wfDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
			files = append(files, path)
		}
		return nil
	})
	_ = filepath.Walk(actDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && (strings.HasSuffix(path, "action.yml") || strings.HasSuffix(path, "action.yaml")) {
			files = append(files, path)
		}
		return nil
	})

	shaRegex := regexp.MustCompile(`@[0-9a-f]{40}$`)
	usesRegex := regexp.MustCompile(`^\s*uses:\s*([^'\s#]+)`)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("❌ Error opening %s: %v\n", file, err)
			errorCount++
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			// Basic check for uses: while excluding local actions (./) and comments
			if strings.Contains(line, "uses:") && !strings.Contains(line, "uses: ./") {
				trimmedLine := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLine, "#") {
					continue
				}

				// Extract action spec
				match := usesRegex.FindStringSubmatch(line)
				if len(match) < 2 {
					continue
				}
				actionSpec := match[1]
				// Clean quotes
				actionSpec = strings.Trim(actionSpec, "\"'")

				if !shaRegex.MatchString(actionSpec) {
					fmt.Printf("❌ Error in %s (Line %d): Action '%s' must be pinned to a 40-character SHA.\n", file, lineNum, actionSpec)
					errorCount++
				}

				// Check for version comment (e.g. # v1.0.0)
				if !strings.Contains(line, "# v") {
					fmt.Printf("⚠️  Warning in %s (Line %d): Action '%s' is missing a version comment (e.g. # v1.0.0).\n", file, lineNum, actionSpec)
				}
			}
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("infrastructure drift detected: %d action(s) are not pinned to a SHA", errorCount)
	}

	fmt.Println("✅ All GitHub Actions are correctly pinned to SHAs.")
	return nil
}

// VerifyCITools ensures only approved CI tools are used in ci.yml.
func VerifyCITools(rootDir string) error {
	ciFile := filepath.Join(rootDir, ".github/workflows/ci.yml")
	fmt.Printf("🔍 Verifying CI tools in %s...\n", ciFile)

	data, err := os.ReadFile(ciFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", ciFile, err)
	}
	content := string(data)

	// Check for Docker Scout (known to fail without entitlement/auth in our environment)
	if strings.Contains(content, "docker/scout-action") {
		return fmt.Errorf("proprietary tool 'docker/scout-action' found without confirmed authentication")
	}

	// Confirm Trivy is used (our preferred open-source alternative)
	if strings.Contains(content, "aquasecurity/trivy-action") {
		fmt.Println("✅ PASS: Using Trivy for container scanning (Open Source, local-friendly).")
	} else {
		fmt.Println("⚠️  WARNING: No container scanner detected in CI (expected Trivy).")
	}

	fmt.Println("✅ CI Toolset Verification Passed")
	return nil
}
