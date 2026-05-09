package maintenance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CompareCoverageThreshold ensures the current threshold isn't lower than the HEAD version.
func CompareCoverageThreshold(path string) error {
	currentVal, err := getThresholdValue(path, "")
	if err != nil {
		return fmt.Errorf("failed to get current threshold: %w", err)
	}

	// Get the threshold from the previous commit (HEAD)
	cmd := exec.Command("git", "show", "HEAD:"+path) //nolint:gosec // maintenance utility
	previousData, err := cmd.Output()
	if err != nil {
		// If the file doesn't exist in HEAD (new file), it's fine.
		return nil
	}

	// Create a temp file to use getThresholdValue on the previous data
	tmpFile, err := os.CreateTemp("", "prev-coverage-*.json")
	if err != nil {
		return nil
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	if _, errWrite := tmpFile.Write(previousData); errWrite != nil {
		return nil
	}
	_ = tmpFile.Close()

	previousVal, err := getThresholdValue(tmpFile.Name(), "")
	if err != nil {
		previousVal = 0.0
	}

	if currentVal < previousVal {
		return fmt.Errorf("coverage threshold cannot be lowered: %.2f < %.2f", currentVal, previousVal)
	}

	return nil
}

func getThresholdValue(path, pkg string) (float64, error) {
	data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return 0, fmt.Errorf("empty threshold file")
	}

	// Try parsing as float first (legacy)
	if val, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return val, nil
	}

	// Try parsing as JSON map
	var thresholds map[string]float64
	if err := json.Unmarshal(data, &thresholds); err != nil {
		return 0, fmt.Errorf("invalid threshold format (not a float or JSON map): %w", err)
	}

	// Look for specific package, then default
	if val, ok := thresholds[pkg]; ok {
		return val, nil
	}
	if val, ok := thresholds["default"]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("no threshold found for package %s or 'default'", pkg)
}
