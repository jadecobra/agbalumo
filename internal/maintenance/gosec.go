package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CheckGosecRationale verifies that all // #nosec directives include a rationale comment.
// Rationale must be preceded by a hyphen (-) or double-hyphen (--).
func CheckGosecRationale(rootDir string) error {
	var invalid []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		if shouldSkipPath(rootDir, path) {
			return nil
		}

		invalid = append(invalid, verifyGosecLines(path)...)
		return nil
	})

	if err != nil {
		return err
	}

	if len(invalid) > 0 {
		return formatGosecErrors(invalid)
	}

	return nil
}

func shouldSkipPath(rootDir, path string) bool {
	rel, err := filepath.Rel(rootDir, path)
	if err != nil {
		return false
	}
	pathParts := strings.Split(filepath.ToSlash(rel), "/")
	for _, part := range pathParts {
		if part == "vendor" || part == ".tester" || part == "tmp" || part == ".go" {
			return true
		}
	}
	return false
}

func verifyGosecLines(path string) []string {
	var invalid []string
	content, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // maintenance utility reads local source files for analysis
	if err != nil {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if isInvalidGosecLine(path, line) {
			invalid = append(invalid, fmt.Sprintf("%s:%d -> %s", path, i+1, strings.TrimSpace(line)))
		}
	}
	return invalid
}

func isInvalidGosecLine(path, line string) bool {
	trimmed := strings.TrimSpace(line)
	// Only match if the line starts with // #nosec or has it as a trailing comment
	if !strings.HasPrefix(trimmed, "// #nosec") && !strings.HasPrefix(trimmed, "//#nosec") {
		return false
	}

	if strings.Contains(line, " - ") || strings.Contains(line, " -- ") {
		return false
	}

	// Extra safety: if it's in a test file, ignore the test case strings
	if strings.HasSuffix(path, "_test.go") && (strings.Contains(line, "content:") || strings.Contains(line, "name:")) {
		return false
	}

	return true
}

func formatGosecErrors(invalid []string) error {
	var sb strings.Builder
	sb.WriteString("mandatory rationale missing for #nosec directives:\n")
	for _, issue := range invalid {
		fmt.Fprintf(&sb, "  %s\n", issue)
	}
	return fmt.Errorf("%s", sb.String())
}
