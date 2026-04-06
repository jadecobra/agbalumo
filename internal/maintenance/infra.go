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

	files := collectActionFiles(rootDir)
	errorCount := 0

	shaRegex := regexp.MustCompile(`@[0-9a-f]{40}$`)
	usesRegex := regexp.MustCompile(`^\s*uses:\s*([^'\s#]+)`)

	for _, file := range files {
		errs := verifyFileSHAs(file, shaRegex, usesRegex)
		errorCount += errs
	}

	if errorCount > 0 {
		return fmt.Errorf("infrastructure drift detected: %d action(s) are not pinned to a SHA", errorCount)
	}

	fmt.Println("✅ All GitHub Actions are correctly pinned to SHAs.")
	return nil
}

func collectActionFiles(rootDir string) []string {
	var files []string
	wfDir := filepath.Join(rootDir, ".github/workflows")
	actDir := filepath.Join(rootDir, ".github/actions")

	walk := func(dir string, suffix string) {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(path, suffix) {
				files = append(files, path)
			}
			return nil
		})
	}

	walk(wfDir, ".yml")
	walk(wfDir, ".yaml")
	walk(actDir, "action.yml")
	walk(actDir, "action.yaml")
	return files
}

func verifyFileSHAs(file string, shaRegex, usesRegex *regexp.Regexp) int {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		fmt.Printf("❌ Error opening %s: %v\n", file, err)
		return 1
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	errorCount := 0
	for scanner.Scan() {
		lineNum++
		if errs := verifyLineSHA(file, lineNum, scanner.Text(), shaRegex, usesRegex); errs > 0 {
			errorCount += errs
		}
	}
	return errorCount
}

func verifyLineSHA(file string, lineNum int, line string, shaRegex, usesRegex *regexp.Regexp) int {
	if !strings.Contains(line, "uses:") || strings.Contains(line, "uses: ./") {
		return 0
	}

	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "#") {
		return 0
	}

	match := usesRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return 0
	}

	actionSpec := strings.Trim(match[1], "\"'")
	errorCount := 0

	if !shaRegex.MatchString(actionSpec) {
		fmt.Printf("❌ Error in %s (Line %d): Action '%s' must be pinned to a 40-character SHA.\n", file, lineNum, actionSpec)
		errorCount++
	}

	if !strings.Contains(line, "# v") {
		fmt.Printf("⚠️  Warning in %s (Line %d): Action '%s' is missing a version comment (e.g. # v1.0.0).\n", file, lineNum, actionSpec)
	}

	return errorCount
}

// VerifyCITools ensures only approved CI tools are used in ci.yml.
func VerifyCITools(rootDir string) error {
	ciFile := filepath.Join(rootDir, ".github/workflows/ci.yml")
	fmt.Printf("🔍 Verifying CI tools in %s...\n", ciFile)

	data, err := os.ReadFile(filepath.Clean(ciFile))
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", ciFile, err)
	}
	content := string(data)

	if strings.Contains(content, "docker/scout-action") {
		return fmt.Errorf("proprietary tool 'docker/scout-action' found without confirmed authentication")
	}

	if strings.Contains(content, "aquasecurity/trivy-action") {
		fmt.Println("✅ PASS: Using Trivy for container scanning (Open Source, local-friendly).")
	} else {
		fmt.Println("⚠️  WARNING: No container scanner detected in CI (expected Trivy).")
	}

	fmt.Println("✅ CI Toolset Verification Passed")
	return nil
}
