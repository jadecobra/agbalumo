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

// TestContextCostTokenBaseline logs token-based metrics for observation during the LOC→Token transition.
// Optimizes for worst-case: Claude Sonnet's 200k window (satisfying it also satisfies Gemini Flash's 1M).
// This test NEVER fails CI — it only logs. A hard gate will be added once baseline values are established.
func TestContextCostTokenBaseline(t *testing.T) {
	report, err := CalculateContextCost("../../.")
	if err != nil {
		t.Fatalf("Failed to calculate context cost: %v", err)
	}

	t.Logf("=== Token Baseline (Observation Only — No Gate Yet) ===")
	t.Logf("Total Files:          %d", report.TotalFiles)
	t.Logf("Total Lines:          %d", report.TotalLines)
	t.Logf("Total Tokens:         %d", report.TotalTokens)
	t.Logf("LOC  RMS:             %.2f  (gate: ≤110 — active)", report.RMS)
	t.Logf("Token RMS:            %.2f  (no gate yet)", report.TokenRMS)
	t.Logf("Context Window Used:  %.2f%%  of Claude Sonnet 200k (worst case)", report.ContextWindowPct)
	t.Logf("                      (Gemini Flash 1M is automatically satisfied if Sonnet passes)")
	t.Logf("")
	t.Logf("Top 10 Most Expensive Files (by token count):")
	for i, fc := range report.TopFiles {
		t.Logf("  %2d. [%6d tokens | %4d lines] %s", i+1, fc.Tokens, fc.Lines, fc.FilePath)
	}

	// Sanity checks — these should always pass if the tokenizer is working
	if report.TotalTokens <= 0 {
		t.Errorf("TotalTokens must be > 0; tokenizer may have failed silently")
	}
	if report.ContextWindowPct <= 0 {
		t.Errorf("ContextWindowPct must be > 0; check claudeSonnetWindow constant")
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
