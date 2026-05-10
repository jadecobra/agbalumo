package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SnapshotParityViolation struct {
	File    string
	Message string
}

// CheckSnapshotParity ensures that every -darwin.png snapshot has a corresponding -linux.png snapshot
// and that they are not bitwise identical (which implies they were likely copied without being generated).
func CheckSnapshotParity(rootDir string) ([]SnapshotParityViolation, error) {
	snapshotDir := filepath.Join(rootDir, "tests/e2e/visual.spec.ts-snapshots")
	if _, err := os.Stat(snapshotDir); os.IsNotExist(err) {
		return nil, nil
	}

	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot directory: %w", err)
	}

	darwin, linux := mapSnapshots(files)
	return findSnapshotViolations(snapshotDir, darwin, linux), nil
}

func mapSnapshots(files []os.DirEntry) (map[string]bool, map[string]bool) {
	darwin := make(map[string]bool)
	linux := make(map[string]bool)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, "-darwin.png") {
			darwin[strings.TrimSuffix(name, "-darwin.png")] = true
		} else if strings.HasSuffix(name, "-linux.png") {
			linux[strings.TrimSuffix(name, "-linux.png")] = true
		}
	}
	return darwin, linux
}

func findSnapshotViolations(snapshotDir string, darwin, linux map[string]bool) []SnapshotParityViolation {
	var violations []SnapshotParityViolation
	for base := range darwin {
		if !linux[base] {
			violations = append(violations, SnapshotParityViolation{
				File:    base + "-darwin.png",
				Message: fmt.Sprintf("Snapshot %s-darwin.png is missing its linux counterpart (%s-linux.png)", base, base),
			})
			continue
		}

		// Check for bitwise identity
		darwinPath := filepath.Join(snapshotDir, base+"-darwin.png")
		linuxPath := filepath.Join(snapshotDir, base+"-linux.png")

		// #nosec G304
		dContent, errD := os.ReadFile(darwinPath)
		// #nosec G304
		lContent, errL := os.ReadFile(linuxPath)

		if errD == nil && errL == nil {
			if string(dContent) == string(lContent) {
				violations = append(violations, SnapshotParityViolation{
					File:    base + "-linux.png",
					Message: "Snapshot is bitwise identical to its darwin counterpart; it was likely copied rather than generated on Linux",
				})
			}
		}
	}
	return violations
}
