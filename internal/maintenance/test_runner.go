package maintenance

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// RunTests executes the Go test suite with race detection and coverage.
// It also enforces the coverage threshold defined in .metrics/coverage or .agents/coverage-threshold.
func RunTests(pkg string, race bool, thresholdPath string, short bool, parallel int) error {
	args := []string{"test", "-v", "-coverprofile=.tester/coverage/coverage.out"}
	if race {
		args = append(args, "-race")
	}
	if short {
		args = append(args, "-short")
	}
	if parallel > 0 {
		args = append(args, "-parallel", strconv.Itoa(parallel))
	}
	args = append(args, pkg)

	// G204: Maintenance utility executes go test
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("🧪 Running tests for %s (race=%v, short=%v, parallel=%d)...\n", pkg, race, short, parallel)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	// Calculate total coverage
	totalCoverage, err := calculateTotalCoverage(".tester/coverage/coverage.out")
	if err != nil {
		return fmt.Errorf("failed to calculate coverage: %w", err)
	}

	fmt.Printf("📊 Total Coverage: %.2f%%\n", totalCoverage)

	// Enforce threshold
	if thresholdPath != "" {
		threshold, err := readThreshold(thresholdPath)
		if err != nil {
			return err
		}
		if totalCoverage < threshold {
			return fmt.Errorf("coverage threshold not met: %.2f%% < %.2f%%", totalCoverage, threshold)
		}
		fmt.Printf("✅ Coverage threshold check passed (%.2f%% >= %.2f%%)\n", totalCoverage, threshold)
	}

	return nil
}

func calculateTotalCoverage(profilePath string) (float64, error) {
	// G204: Maintenance utility executes go tool cover
	cmd := exec.Command("go", "tool", "cover", "-func="+profilePath) //nolint:gosec // maintenance utility
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to run go tool cover: %w", err)
	}

	// Find the total line: "total: (statements) 86.1%"
	re := regexp.MustCompile(`total:\s+\(statements\)\s+(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return 0, fmt.Errorf("failed to parse coverage output")
	}

	return strconv.ParseFloat(matches[1], 64)
}

func readThreshold(path string) (float64, error) {
	data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		if os.IsNotExist(err) {
			return 80.0, nil // Default to 80% if file missing
		}
		return 0, fmt.Errorf("failed to read threshold file: %w", err)
	}
	return strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
}
