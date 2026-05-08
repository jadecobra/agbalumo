package maintenance

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// DeprecatedViolation represents a violation of the deprecated pattern rule.
type DeprecatedViolation struct {
	File       string
	Pattern    string
	Suggestion string
	Line       int
}

// CheckDeprecatedPatterns scans the codebase for deprecated patterns.
func CheckDeprecatedPatterns(rootDir string) ([]DeprecatedViolation, error) {
	var violations []DeprecatedViolation

	// #nosec G122 -- local maintenance utility
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		v, err := scanFileForDeprecated(path)
		if err != nil {
			return err
		}
		violations = append(violations, v...)
		return nil
	})

	return violations, err
}

func scanFileForDeprecated(path string) ([]DeprecatedViolation, error) {
	// Check for map[string]interface{} in internal/module and internal/handler
	isModuleOrHandler := strings.Contains(path, "internal/module") || strings.Contains(path, "internal/handler")

	// #nosec G304 -- local maintenance utility
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var violations []DeprecatedViolation
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if isModuleOrHandler && strings.Contains(line, "map[string]interface{}") {
			violations = append(violations, DeprecatedViolation{
				File:       path,
				Line:       lineNum,
				Pattern:    "map[string]interface{}",
				Suggestion: "Use typed ViewModel struct",
			})
		}

		if !strings.Contains(path, "internal/maintenance/") && strings.Contains(line, "RenderWithBaseContext") {
			violations = append(violations, DeprecatedViolation{
				File:       path,
				Line:       lineNum,
				Pattern:    "RenderWithBaseContext",
				Suggestion: "Use RenderTyped with typed ViewModel",
			})
		}
	}

	return violations, scanner.Err()
}
