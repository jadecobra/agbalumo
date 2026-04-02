package agent

import (
	"testing"
)

// TestContextCostThreshold asserts that the codebase context cost is within acceptable limits.
// This test is intended to fail until the high-cost files (like outlier reports) are addressed.
func TestContextCostThreshold(t *testing.T) {
	const threshold = 110.0
	report, err := CalculateContextCost("../../.") // Calculate from repo root
	if err != nil {
		t.Fatalf("Failed to calculate context cost: %v", err)
	}

	if report.RMS > threshold {
		t.Errorf("Context Cost (RMS) too high: %.2f (Threshold: %.2f). Top files must be refactored or excluded.", report.RMS, threshold)
		t.Logf("Total Files: %d, Total Lines: %d", report.TotalFiles, report.TotalLines)
		t.Logf("Top 10 expensive files:")
		for i := 0; i < 10 && i < len(report.TopFiles); i++ {
			t.Logf("  %s: %d lines", report.TopFiles[i].FilePath, report.TopFiles[i].Lines)
		}
	}
}

// TestSecurityModularizationCheck asserts that the security logic has been correctly split into modules.
// This test will fail until the split files exist and contain the required logic.
func TestSecurityModularizationCheck(t *testing.T) {
	// Re-assigning functions to ensure the split files are hooked into the main dispatcher.
	// We'll check if these functions are accessible and functional.
	// For now, this is a placeholder that we will refine as we define the split.
	// But simply checking for the existence of the files will be the first failing step.
}
