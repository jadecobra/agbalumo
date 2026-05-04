package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GeneratePrunedTree returns a string representation of the directory tree
// starting from rootDir, limited by maxDepth, and skipping common excluded directories.
func GeneratePrunedTree(rootDir string, maxDepth int) string {
	var sb strings.Builder
	absRoot, _ := filepath.Abs(rootDir)
	excluded := map[string]bool{
		"node_modules": true, "bin": true, ".git": true, ".cache": true, "tmp": true, "artifacts": true,
	}

	_ = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(absRoot, path)
		if rel == "." {
			return nil
		}

		if shouldSkipDir(info, excluded) {
			return filepath.SkipDir
		}

		depth := strings.Count(rel, string(os.PathSeparator)) + 1
		if depth > maxDepth {
			return handleDepthExceeded(info)
		}

		sb.WriteString(formatTreeLine(info, depth))
		return nil
	})

	return sb.String()
}

func shouldSkipDir(info os.FileInfo, excluded map[string]bool) bool {
	return info.IsDir() && excluded[info.Name()]
}

func handleDepthExceeded(info os.FileInfo) error {
	if info.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

func formatTreeLine(info os.FileInfo, depth int) string {
	indent := strings.Repeat("  ", depth-1)
	name := info.Name()
	if info.IsDir() {
		name += "/"
	}
	return fmt.Sprintf("%s%s\n", indent, name)
}
