package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// ParseCoverageProfile reads a standard go test -coverprofile output
// and calculates the percentage coverage for each package.
func ParseCoverageProfile(r io.Reader) (map[string]float64, error) {
	scanner := bufio.NewScanner(r)
	
	packageStmts := make(map[string]int)
	packageCovered := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines or the mode line
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}

		// Format: name:line.col,line.col numStmts count
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		filePart := parts[0]
		numStmtsStr := parts[1]
		countStr := parts[2]

		// Extract package from file path
		// e.g. github.com/jadecobra/agbalumo/internal/handler/auth.go:...
		idx := strings.LastIndex(filePart, "/")
		if idx == -1 {
			continue
		}
		
		pkgPath := filePart[:idx]
		
		numStmts, err := strconv.Atoi(numStmtsStr)
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			continue
		}

		packageStmts[pkgPath] += numStmts
		if count > 0 {
			packageCovered[pkgPath] += numStmts
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	coverageByPkg := make(map[string]float64)
	for pkgPath, total := range packageStmts {
		if total == 0 {
			coverageByPkg[pkgPath] = 0.0
		} else {
			covered := packageCovered[pkgPath]
			coverageByPkg[pkgPath] = (float64(covered) / float64(total)) * 100.0
		}
	}

	return coverageByPkg, nil
}

// ParseThresholds parses a JSON mapping of package paths to minimum coverage thresholds.
// It supports a special "default" key.
func ParseThresholds(data []byte) (map[string]float64, error) {
	if len(data) == 0 {
		return map[string]float64{"default": 90.0}, nil
	}

	var thresholds map[string]float64
	if err := json.Unmarshal(data, &thresholds); err != nil {
		return nil, fmt.Errorf("failed to parse coverage thresholds: %w", err)
	}

	return thresholds, nil
}

// EnforceCoverage checks the measured coverage against expected thresholds.
// Returns a list of formatted violation strings, sorted alphabetically.
func EnforceCoverage(coverage map[string]float64, thresholds map[string]float64) []string {
	defaultThreshold, hasDefault := thresholds["default"]
	
	var violations []string

	for pkg, pct := range coverage {
		expected, ok := thresholds[pkg]
		if ok {
			if pct < expected {
				violations = append(violations, fmt.Sprintf("%s: coverage %.1f%% is below explicit threshold of %.1f%%", pkg, pct, expected))
			}
		} else if hasDefault {
			if pct < defaultThreshold {
				violations = append(violations, fmt.Sprintf("%s: coverage %.1f%% is below default threshold of %.1f%%", pkg, pct, defaultThreshold))
			}
		}
	}

	sort.Strings(violations)
	return violations
}
