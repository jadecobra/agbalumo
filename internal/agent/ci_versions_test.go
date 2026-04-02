package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCIVersionsNode24(t *testing.T) {
	// Find project root
	root := "../../"
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		root = "./" // Fallback for local run
	}

	workflowsDir := filepath.Join(root, ".github/workflows")
	customActionsDir := filepath.Join(root, ".github/actions")

	tests := []struct {
		file     string
		action   string
		minVers  string
		critical bool
	}{
		// Main Workflow
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/checkout", minVers: "v6.0.2"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/setup-go", minVers: "v6.4.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "golangci/golangci-lint-action", minVers: "v9.2.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/upload-artifact", minVers: "v6.0.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/setup-node", minVers: "v6.3.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "docker/setup-buildx-action", minVers: "v4.0.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "docker/build-push-action", minVers: "v7.0.0"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "aquasecurity/trivy-action", minVers: "v0.35.0", critical: true},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "superfly/flyctl-actions/setup-flyctl", minVers: "v1.5"},

		// Custom Action
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "actions/setup-go", minVers: "v6.4.0"},
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "actions/cache", minVers: "v5.0.4"},
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "arduino/setup-task", minVers: "v2.0.0"},
	}

	for _, tc := range tests {
		t.Run(filepath.Base(tc.file)+"/"+tc.action, func(t *testing.T) {
			content, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tc.file, err)
			}

			// Simple check for "uses: action@minVers" or later
			// Since we use SHAs, we also check for the tag in the comment if present
			lines := strings.Split(string(content), "\n")
			found := false
			for _, line := range lines {
				if strings.Contains(line, "uses: "+tc.action) {
					// Check if the current version in the line is at least minVers
					// For simplicity in RED phase, we just check if it contains the exact minVers or higher
					// If it contains a lower version (e.g. v6.0.2 vs v6.1.0), it should fail.
					if strings.Contains(line, tc.minVers) {
						found = true
						break
					}
					// If it contains a higher version, we should also pass (not implemented here for RED simplicity)
				}
			}

			if !found {
				t.Errorf("Action %s in %s does not meet minimum version %s", tc.action, tc.file, tc.minVers)
			}
		})
	}
}
