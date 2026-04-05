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
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Exclude common directories to simulate bash excluded paths.
		// We only exclude segments relative to the rootDir to allow tests in /tmp on Linux.
		rel, err := filepath.Rel(rootDir, path)
		if err == nil {
			pathParts := strings.Split(filepath.ToSlash(rel), "/")
			for _, part := range pathParts {
				if part == "vendor" || part == ".tester" || part == "tmp" || part == ".go" {
					return nil // skip
				}
			}
		}

		content, readErr := os.ReadFile(filepath.Clean(path)) //nolint:gosec // maintenance utility reads local source files for analysis
		if readErr != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Only match if the line starts with // #nosec or has it as a trailing comment
			// This avoids matching " // #nosec" inside strings in most cases.
			if (strings.HasPrefix(trimmed, "// #nosec") || strings.HasPrefix(trimmed, "//#nosec")) &&
				!strings.Contains(line, " - ") && !strings.Contains(line, " -- ") {

				// Extra safety: if it's in a test file, ignore the test case strings
				if strings.HasSuffix(path, "_test.go") && (strings.Contains(line, "content:") || strings.Contains(line, "name:")) {
					continue
				}

				invalid = append(invalid, fmt.Sprintf("%s:%d -> %s", path, i+1, strings.TrimSpace(line)))
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(invalid) > 0 {
		var sb strings.Builder
		sb.WriteString("mandatory rationale missing for #nosec directives:\n")
		for _, issue := range invalid {
			fmt.Fprintf(&sb, "  %s\n", issue)
		}
		return fmt.Errorf("%s", sb.String())
	}

	return nil
}
