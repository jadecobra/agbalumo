package maintenance

import (
	"os"
	"path/filepath"
	"strings"
)

var excludedDirs = map[string]bool{
	"vendor":       true,
	".git":         true,
	"testdata":     true,
	".tester":      true,
	"node_modules": true,
	"cmd/verify":   true,
}

// CheckAgentsCoverage walks the rootDir and finds Go packages missing AGENTS.md.
// A Go package is defined as a directory containing at least one .go file.
func CheckAgentsCoverage(rootDir string) ([]string, error) {
	var missing []string
	err := filepath.Walk(rootDir, walker(rootDir, &missing))
	return missing, err
}

func walker(rootDir string, missing *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		rel, err := filepath.Rel(rootDir, path)
		if err != nil || rel == "." {
			return err
		}

		if isExcluded(rel) {
			return filepath.SkipDir
		}

		return checkDir(path, rel, missing)
	}
}

func checkDir(path, rel string, missing *[]string) error {
	isPkg, hasAgents, err := analyzeDir(path)
	if err != nil {
		return err
	}

	if isPkg && !hasAgents {
		*missing = append(*missing, rel)
	}
	return nil
}

func isExcluded(rel string) bool {
	parts := strings.Split(rel, string(os.PathSeparator))
	for i := 1; i <= len(parts); i++ {
		subPath := filepath.Join(parts[:i]...)
		if excludedDirs[subPath] {
			return true
		}
	}
	return false
}

func analyzeDir(path string) (bool, bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, false, err
	}

	hasGo := false
	hasAgents := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".go") {
			hasGo = true
		}
		if entry.Name() == "AGENTS.md" {
			hasAgents = true
		}
	}
	return hasGo, hasAgents, nil
}
