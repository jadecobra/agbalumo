package agent

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
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
			return nil, fmt.Errorf("malformed coverage line: %s", line)
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

// CoverageConfig holds the coverage requirements securely.
type CoverageConfig struct {
	Thresholds map[string]float64 `json:"thresholds"`
	Signature  string             `json:"signature"`
}

func calculateCoverageSignature(c *CoverageConfig) string {
	copy := *c
	copy.Signature = "" // exclude signature itself from hash

	// predictable hashing by marshalling
	b, _ := json.Marshal(copy)
	hash := sha256.Sum256(b)
	return fmt.Sprintf("%x", hash)
}

// ParseThresholds parses a JSON mapping of package paths to minimum coverage thresholds.
// It supports a special "default" key. It now enforces an Anti-Cheat signature.
func ParseThresholds(data []byte) (map[string]float64, error) {
	if len(data) == 0 {
		return map[string]float64{"default": 90.0}, nil
	}

	var config CoverageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse coverage thresholds: %w", err)
	}

	if config.Signature != "" {
		expected := calculateCoverageSignature(&config)
		if config.Signature != expected {
			return nil, fmt.Errorf("ANTI-CHEAT TRIGGERED: Manual modification of .agents/coverage-thresholds.json detected")
		}
	} else {
		// Missing signature is also spoofing/tampering, reject it
		return nil, fmt.Errorf("ANTI-CHEAT TRIGGERED: Manual modification of .agents/coverage-thresholds.json detected")
	}

	if config.Thresholds == nil {
		return map[string]float64{"default": 90.0}, nil
	}
	return config.Thresholds, nil
}

// SaveThresholds signs and persists the thresholds.
func SaveThresholds(path string, thresholds map[string]float64) error {
	config := CoverageConfig{
		Thresholds: thresholds,
	}
	config.Signature = calculateCoverageSignature(&config)

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return util.SafeWriteFile(path, b)
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
