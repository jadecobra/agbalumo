package maintenance

import (
	"fmt"
	"os"
	"strings"
)

// RootWhitelist is the set of files and directories allowed in the project root.
var RootWhitelist = map[string]bool{
	"AGENTS.md":            true,
	"CONTRIBUTING.md":      true,
	"Dockerfile":           true,
	"README.md":            true,
	"docker-compose.yml":   true,
	"fly.toml":             true,
	"go.mod":               true,
	"go.sum":               true,
	"main.go":              true,
	"package-lock.json":    true,
	"package.json":         true,
	"pnpm-lock.yaml":       true,
	"pnpm-workspace.yaml":  true,
	"tailwind.config.js":   true,
	"playwright.config.ts": true,
}

// DirWhitelist is the set of directories allowed in the project root.
var DirWhitelist = map[string]bool{
	"artifacts":         true,
	"bin":               true,
	"certs":             true,
	"cmd":               true,
	"config":            true,
	"docs":              true,
	"etc":               true,
	"internal":          true,
	"logs":              true,
	"node_modules":      true,
	"playwright-report": true,
	"scripts":           true,
	"test-results":      true,
	"tests":             true,
	"tools":             true,
	"ui":                true,
}

// VerifyRootHygiene checks for unexpected files in the root directory.
func VerifyRootHygiene(rootDir string) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return fmt.Errorf("failed to read root directory: %v", err)
	}

	var violations []string
	for _, entry := range entries {
		if !isAllowed(entry) {
			if entry.IsDir() {
				violations = append(violations, fmt.Sprintf("Directory: %s/", entry.Name()))
			} else {
				violations = append(violations, fmt.Sprintf("File: %s", entry.Name()))
			}
		}
	}

	if len(violations) > 0 {
		return fmt.Errorf("root directory clutter detected (use 'verify janitor' or move to bin/):\n  - %s", strings.Join(violations, "\n  - "))
	}

	return nil
}

func isAllowed(entry os.DirEntry) bool {
	name := entry.Name()

	// Ignore hidden files (starts with dot)
	if strings.HasPrefix(name, ".") {
		return true
	}

	if entry.IsDir() {
		return DirWhitelist[name]
	}

	return RootWhitelist[name]
}
