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

func TestCINodeRuntime24(t *testing.T) {
	// Find project root
	root := "../../"
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		root = "./"
	}

	ciFile := filepath.Join(root, ".github/workflows/ci.yml")
	content, err := os.ReadFile(ciFile)
	if err != nil {
		t.Fatalf("failed to read %s: %v", ciFile, err)
	}

	lines := strings.Split(string(content), "\n")
	nodeVersionLineFound := false
	envNodeVersion := ""

	// 1. Extract top-level NODE_VERSION from env block
	for _, line := range lines {
		if strings.Contains(line, "NODE_VERSION:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				envNodeVersion = strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			}
		}
	}

	// 2. Verify all node-version usages
	for _, line := range lines {
		if strings.Contains(line, "node-version:") {
			nodeVersionLineFound = true
			trimmed := strings.TrimSpace(line)

			// Support both hardcoded '24' and interpolated expression
			if strings.Contains(trimmed, "${{ env.NODE_VERSION }}") {
				if envNodeVersion != "24" {
					t.Errorf("In %s, node-version uses expression but global NODE_VERSION is '%s' (expected '24')", ciFile, envNodeVersion)
				}
			} else if !strings.Contains(trimmed, "'24'") && !strings.Contains(trimmed, "\"24\"") && !strings.Contains(trimmed, ": 24") {
				t.Errorf("In %s, found invalid node-version: %s (expected '24' or expression)", ciFile, trimmed)
			}
		}
	}

	if !nodeVersionLineFound {
		t.Errorf("In %s, no node-version found", ciFile)
	}
}

func TestPackageNodeEngine24(t *testing.T) {
	// Find project root
	root := "../../"
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		root = "./"
	}

	pkgFile := filepath.Join(root, "package.json")
	content, err := os.ReadFile(pkgFile)
	if err != nil {
		t.Fatalf("failed to read %s: %v", pkgFile, err)
	}

	if !strings.Contains(string(content), "\"node\": \">=24\"") &&
		!strings.Contains(string(content), "\"node\": \"^24") &&
		!strings.Contains(string(content), "\"node\": \"24") {
		t.Errorf("In %s, engines.node does not specify version 24", pkgFile)
	}
}
