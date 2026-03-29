package agent

import (
	"strings"
	"testing"
)

func TestVerifySecurityStaticGate_Basic(t *testing.T) {
	// Simple test to ensure it doesn't panic and handles empty paths
	// We can't fully test truncation here without complex filesystem mocks
	// but we can at least invoke it.
	VerifySecurityStaticGate()
}

func TestSummarizeTestFailures_Truncation(t *testing.T) {
	failures := []TestFailure{
		{TestName: "Test1", Output: "Error 1"},
		{TestName: "Test2", Output: "Error 2"},
		{TestName: "Test3", Output: "Error 3"},
	}

	output := captureStdout(t, func() {
		SummarizeTestFailures(failures, 1)
	})

	if !strings.Contains(output, "Test1") {
		t.Errorf("expected Test1 in output, got: %s", output)
	}
	if strings.Contains(output, "Test2") {
		t.Error("did not expect Test2 in output (truncated)")
	}
	if !strings.Contains(output, "... and 2 more failures.") {
		t.Errorf("expected summary count, got: %s", output)
	}
}
