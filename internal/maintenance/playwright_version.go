package maintenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyPlaywrightVersionParity checks if the Playwright version in package.json
// matches the Docker image version used in ci_runner.go.
// Replaces Strict Lesson: Playwright Local Linux Gate
func VerifyPlaywrightVersionParity(cwd string) error {
	packageJSONPath := filepath.Join(cwd, "package.json")
	data, err := os.ReadFile(packageJSONPath) //nolint:gosec // G304: Maintenance utility reads project config
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err = json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	playwrightVer, ok := pkg.DevDependencies["@playwright/test"]
	if !ok {
		return fmt.Errorf("@playwright/test not found in devDependencies")
	}

	// Strip caret/tilde
	playwrightVer = strings.TrimLeft(playwrightVer, "^~")

	ciRunnerPath := filepath.Join(cwd, "internal", "maintenance", "ci_runner.go")
	ciData, err := os.ReadFile(ciRunnerPath) //nolint:gosec // G304: Maintenance utility reads source code for version parity
	if err != nil {
		return fmt.Errorf("failed to read ci_runner.go: %w", err)
	}

	re := regexp.MustCompile(`playwright:v([0-9.]+)-`)
	match := re.FindStringSubmatch(string(ciData))
	if len(match) < 2 {
		return fmt.Errorf("failed to find Playwright Docker image tag in ci_runner.go")
	}

	dockerVer := match[1]

	// Compare major.minor
	pParts := strings.Split(playwrightVer, ".")
	dParts := strings.Split(dockerVer, ".")

	if len(pParts) < 2 || len(dParts) < 2 {
		return fmt.Errorf("invalid version format: package=%s, docker=%s", playwrightVer, dockerVer)
	}

	if pParts[0] != dParts[0] || pParts[1] != dParts[1] {
		return fmt.Errorf("Playwright version mismatch: package.json has %s, but ci_runner.go uses %s (Docker). Please sync them.", playwrightVer, dockerVer)
	}

	return nil
}
