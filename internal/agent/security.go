package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

var (
	// internalOpen is a hook for testing file operations.
	internalOpen = util.SafeOpen
)

// VerifySecurityStatic runs static analysis checkers for security vulnerabilities on multiple targets.
func VerifySecurityStatic(targets ...string) ([]SecurityViolation, error) {
	var allViolations []SecurityViolation

	for _, target := range targets {
		info, err := util.SafeStat(target)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("failed to stat target %s: %w", target, err)
		}

		if !info.IsDir() {
			if strings.HasSuffix(target, "_test.go") {
				continue
			}
			violations, fErr := checkFile(target)
			if fErr != nil {
				return nil, fmt.Errorf("failed to check file %s: %w", target, fErr)
			}
			allViolations = append(allViolations, violations...)
			continue
		}

		err = filepath.Walk(target, func(path string, info os.FileInfo, wErr error) error {
			if wErr != nil {
				return wErr
			}

			if info.IsDir() {
				if strings.HasPrefix(info.Name(), ".") && info.Name() != "." && info.Name() != target {
					return filepath.SkipDir
				}
				return nil
			}

			if strings.Contains(path, "/vendor/") || strings.Contains(path, "/node_modules/") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".pdf" || ext == ".exe" || ext == ".bin" {
				return nil
			}

			violations, fErr := checkFile(path)
			if fErr != nil {
				if !strings.HasSuffix(path, ".go") {
					return nil
				}
				return fmt.Errorf("failed to check file %s: %w", path, fErr)
			}

			allViolations = append(allViolations, violations...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return deduplicateViolations(allViolations), nil
}
