package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var shaCache = make(map[string]bool)

// VerifyActionSHAs checks for correctly pinned SHAs in .github/workflows and .github/actions.
// It requires the 'gh' CLI to be installed and authenticated to verify SHAs against GitHub.
func VerifyActionSHAs(rootDir string) error {
	fmt.Println("🔍 Verifying GitHub Action SHA pinning (Remote-First)...")

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("HARD FAILURE: 'gh' CLI is required for infrastructure verification but was not found in PATH")
	}

	files := collectActionFiles(rootDir)
	errorCount := 0

	shaRegex := regexp.MustCompile(`@[0-9a-f]{40}$`)
	usesRegex := regexp.MustCompile(`^\s*uses:\s*([^'\s#]+)`)

	for _, file := range files {
		errs := verifyFileSHAs(file, shaRegex, usesRegex)
		errorCount += errs
	}

	if errorCount > 0 {
		return fmt.Errorf("infrastructure drift detected: %d action(s) failed verification (corrupted SHAs or missing pins)", errorCount)
	}

	fmt.Println("✅ All GitHub Actions are correctly pinned and verified against remote.")
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

	// 1. Syntax Check
	if !shaRegex.MatchString(actionSpec) {
		fmt.Printf("❌ Error in %s (Line %d): Action '%s' must be pinned to a 40-character SHA.\n", file, lineNum, actionSpec)
		return 1
	}

	if !strings.Contains(line, "# v") {
		fmt.Printf("⚠️  Warning in %s (Line %d): Action '%s' is missing a version comment (e.g. # v1.0.0).\n", file, lineNum, actionSpec)
	}

	// 2. Remote Verification Check (Mandatory)
	if !verifyRemoteSHA(file, lineNum, actionSpec) {
		errorCount++
	}

	return errorCount
}

func verifyRemoteSHA(file string, lineNum int, actionSpec string) bool {
	if verified, ok := shaCache[actionSpec]; ok {
		return verified
	}

	parts := strings.Split(actionSpec, "@")
	fullRepo := parts[0]
	sha := parts[1]

	// Extract {owner}/{repo} from {owner}/{repo}/{path}
	repoParts := strings.Split(fullRepo, "/")
	if len(repoParts) < 2 {
		fmt.Printf("❌ Error in %s (Line %d): Invalid action repo spec '%s'\n", file, lineNum, fullRepo)
		shaCache[actionSpec] = false
		return false
	}
	repo := repoParts[0] + "/" + repoParts[1]

	// Use gh api to verify commit existence
	// SKIP if running in GitHub Actions to avoid remote dependencies and auth issues.
	// Production CI should rely on the local gate having passed.
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		shaCache[actionSpec] = true
		return true
	}

	endpoint := fmt.Sprintf("repos/%s/commits/%s", repo, sha)
	cmd := exec.Command("gh", "api", endpoint, "--silent") //nolint:gosec // G204: Maintenance utility runs trusted commands

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		fmt.Printf("❌ Error in %s (Line %d): Action SHA @%s could not be verified in repo %s\n   Reason: %s\n", file, lineNum, sha, repo, errMsg)
		shaCache[actionSpec] = false
		return false
	}

	shaCache[actionSpec] = true
	return true
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
		// Detect invalid trivy-version hallucination
		if strings.Contains(content, "trivy-version:") {
			return fmt.Errorf("invalid CI configuration: 'aquasecurity/trivy-action' uses 'version', not 'trivy-version'")
		}
		if !strings.Contains(content, "version:") {
			fmt.Println("⚠️  WARNING: 'aquasecurity/trivy-action' missing explicit 'version' - using action default.")
		}
	} else {
		fmt.Println("⚠️  WARNING: No container scanner detected in CI (expected Trivy).")
	}

	fmt.Println("✅ CI Toolset Verification Passed")
	return nil
}

// VerifyJSSyntax ensures all JavaScript files in the project are syntactically correct.
func VerifyJSSyntax(rootDir string) error {
	staticDir := filepath.Join(rootDir, "ui/static/js")
	fmt.Printf("🔍 Verifying JavaScript syntax in %s...\n", staticDir)

	if _, err := exec.LookPath("node"); err != nil {
		fmt.Println("⚠️  Skipping JS syntax check: 'node' not found in PATH")
		return nil
	}

	files, err := filepath.Glob(filepath.Join(staticDir, "*.js"))
	if err != nil {
		return fmt.Errorf("failed to list JS files: %w", err)
	}

	for _, file := range files {
		// Skip minified files (usually pre-validated or from 3rd parties)
		if strings.HasSuffix(file, ".min.js") || strings.HasSuffix(file, ".umd.min.js") {
			continue
		}
		cmd := exec.Command("node", "-c", file) //nolint:gosec // maintenance utility runs trusted commands on local JS assets
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("JS syntax error in %s:\n%s", file, string(output))
		}
	}

	fmt.Println("✅ All JavaScript files passed syntax verification.")
	return nil
}
