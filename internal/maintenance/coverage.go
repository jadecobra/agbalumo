package maintenance

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CompareCoverageThreshold ensures the current threshold isn't lower than the HEAD version.
func CompareCoverageThreshold(path string) error {
	// G304: Maintenance utility reads coverage threshold file
	currentData, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return fmt.Errorf("failed to read current threshold: %w", err)
	}
	currentVal, err := strconv.ParseFloat(strings.TrimSpace(string(currentData)), 64)
	if err != nil {
		return fmt.Errorf("invalid current threshold value: %w", err)
	}

	// Get the threshold from the previous commit (HEAD)
	// G204: Maintenance utility uses git show to compare thresholds
	cmd := exec.Command("git", "show", "HEAD:"+path) //nolint:gosec // maintenance utility
	previousData, err := cmd.Output()
	if err != nil {
		// If the file doesn't exist in HEAD (new file), it's fine.
		return nil
	}

	previousVal, err := strconv.ParseFloat(strings.TrimSpace(string(previousData)), 64)
	if err != nil {
		// If we can't parse the previous value, assume 0.0
		previousVal = 0.0
	}

	if currentVal < previousVal {
		return fmt.Errorf("coverage threshold cannot be lowered: %.2f < %.2f", currentVal, previousVal)
	}

	return nil
}
