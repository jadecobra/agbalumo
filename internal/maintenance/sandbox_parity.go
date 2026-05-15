package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// SandboxParityViolation represents a component that is defined but not used in the sandbox.
type SandboxParityViolation struct {
	Component string
	File      string
	Message   string
}

// CheckSandboxParity verifies that all components defined in ui_components.html are documented in sandbox.html.
func CheckSandboxParity(rootDir string) ([]SandboxParityViolation, error) {
	componentsFile := filepath.Join(rootDir, "ui", "templates", "partials", "ui_components.html")
	sandboxFile := filepath.Join(rootDir, "ui", "templates", "sandbox.html")

	defined, err := extractBlocks(componentsFile, `\{\{\s*define\s+"([^"]+)"`)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	referenced, err := extractBlocks(sandboxFile, `\{\{\s*template\s+"([^"]+)"`)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return findMissingComponents(defined, referenced), nil
}

func findMissingComponents(defined, referenced []string) []SandboxParityViolation {
	exclude := map[string]bool{
		"base.html": true,
		"title":     true,
		"content":   true,
		"head":      true,
		"scripts":   true,
	}

	refMap := make(map[string]bool)
	for _, r := range referenced {
		refMap[r] = true
	}

	var violations []SandboxParityViolation
	for _, d := range defined {
		if exclude[d] || refMap[d] {
			continue
		}
		violations = append(violations, SandboxParityViolation{
			Component: d,
			File:      "ui/templates/partials/ui_components.html",
			Message:   fmt.Sprintf("Component %q is defined but not documented in sandbox.html", d),
		})
	}
	return violations
}

func extractBlocks(path string, pattern string) ([]string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	re := regexp.MustCompile(pattern)
	var names []string
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
			if len(match) > 1 && !seen[match[1]] {
				names = append(names, match[1])
				seen[match[1]] = true
			}
		}
	}

	return names, scanner.Err()
}
