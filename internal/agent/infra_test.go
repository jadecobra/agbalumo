package agent

import (
	"os/exec"
	"testing"
)

func TestInfrastructureActionSHAs(t *testing.T) {
	// Root of the workspace
	cmd := exec.Command("../../scripts/ci/verify-action-shas.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Infrastructure SHA verification failed: %s\nOutput: %s", err, string(output))
	}
}
