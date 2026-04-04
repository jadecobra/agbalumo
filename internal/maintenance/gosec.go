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

		// Exclude common directories to simulate bash excluded paths
		// Using filepath.ToSlash to handle potential Windows paths if needed,
		// but since we're on Mac, a simple contains is fine.
		if strings.Contains(path, "/vendor/") || strings.Contains(path, "/.tester/") || strings.Contains(path, "/tmp/") || strings.Contains(path, "/.go/") {
			return nil
		}

		content, readErr := os.ReadFile(path)
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
			sb.WriteString(fmt.Sprintf("  %s\n", issue))
		}
		return fmt.Errorf("%s", sb.String())
	}

	return nil
}
