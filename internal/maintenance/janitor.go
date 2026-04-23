package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
)

// RunJanitor moves stale root-level artifacts to .tester/
func RunJanitor(rootDir string) error {
	patterns := []string{
		"critique_full.txt",
		"*_ci.yml",
		"*_ci_job.log",
		"precommit_check.txt",
		"index_test.html",
		"server.log",
	}

	testerDir := filepath.Join(rootDir, ".tester")

	// Find matches for each pattern
	var matches []string
	for _, p := range patterns {
		m, err := filepath.Glob(filepath.Join(rootDir, p))
		if err != nil {
			return fmt.Errorf("failed to glob pattern %s: %v", p, err)
		}
		matches = append(matches, m...)
	}

	if len(matches) == 0 {
		return nil
	}

	// Ensure .tester/ exists
	if err := os.MkdirAll(testerDir, 0750); err != nil {
		return fmt.Errorf("failed to create .tester directory: %v", err)
	}

	// Move each match to .tester/
	for _, oldPath := range matches {
		filename := filepath.Base(oldPath)
		newPath := filepath.Join(testerDir, filename)

		fmt.Printf("Cleaning up: %s -> %s\n", filename, filepath.Join(".tester", filename))
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to move %s to %s: %v", oldPath, newPath, err)
		}
	}

	return nil
}
