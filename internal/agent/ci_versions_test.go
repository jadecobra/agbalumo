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
		sha      string
		critical bool
	}{
		// Main Workflow
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/checkout", minVers: "v6.0.2", sha: "de0fac2e4500dabe0009e67214ff5f5447ce83dd"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "golangci/golangci-lint-action", minVers: "v9.2.0", sha: "1e7e51e771db61008b38414a730f564565cf7c20"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/upload-artifact", minVers: "v7.0.0", sha: "bbbca2ddaa5d8feaa63e36b76fdaad77386f024f"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "actions/setup-node", minVers: "v6.3.0", sha: "53b83947a5a98c8d113130e565377fae1a50d02f"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "docker/setup-buildx-action", minVers: "v4.0.0", sha: "4d04d5d9486b7bd6fa91e7baf45bbb4f8b9deedd"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "docker/build-push-action", minVers: "v7.0.0", sha: "d08e5c354a6adb9ed34480a06d141179aa583294"},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "aquasecurity/trivy-action", minVers: "v0.35.0", sha: "57a97c7e7821a5776cebc9bb87c984fa69cba8f1", critical: true},
		{file: filepath.Join(workflowsDir, "ci.yml"), action: "superfly/flyctl-actions/setup-flyctl", minVers: "v1.5", sha: "fc53c09e1bc3be6f54706524e3b82c4f462f77be"},

		// Custom Action
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "actions/setup-go", minVers: "v6.4.0", sha: "4a3601121dd01d1626a1e23e37211e3254c1c06c"},
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "actions/cache", minVers: "v5.0.4", sha: "668228422ae6a00e4ad889ee87cd7109ec5666a7"},
		{file: filepath.Join(customActionsDir, "setup-task-with-cache/action.yml"), action: "arduino/setup-task", minVers: "v2.0.0", sha: "b91d5d2c96a56797b48ac1e0e89220bf64044611"},
	}

	for _, tc := range tests {
		t.Run(filepath.Base(tc.file)+"/"+tc.action, func(t *testing.T) {
			content, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tc.file, err)
			}

			lines := strings.Split(string(content), "\n")
			occurrences := 0
			for i, line := range lines {
				if strings.Contains(line, "uses: "+tc.action) {
					occurrences++
					// 1. Verify SHA part
					if !strings.Contains(line, tc.sha) {
						t.Errorf("%s:L%d - Action %s uses invalid SHA (expected %s)", tc.file, i+1, tc.action, tc.sha)
					}
					// 2. Verify comment version (if present)
					if !strings.Contains(line, tc.minVers) {
						t.Errorf("%s:L%d - Action %s comment mismatch (expected %s)", tc.file, i+1, tc.action, tc.minVers)
					}
				}
			}

			if occurrences == 0 {
				t.Errorf("Action %s not found in %s", tc.action, tc.file)
			}
		})
	}
}

func TestCIEnvFlagsNode24(t *testing.T) {
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
	flagFound := false
	for _, line := range lines {
		if strings.Contains(line, "FORCE_JAVASCRIPT_ACTIONS_TO_NODE24:") {
			flagFound = true
			if !strings.Contains(line, "true") {
				t.Errorf("FORCE_JAVASCRIPT_ACTIONS_TO_NODE24 must be set to true in %s", ciFile)
			}
		}
	}

	if !flagFound {
		t.Errorf("FORCE_JAVASCRIPT_ACTIONS_TO_NODE24 flag not found in %s", ciFile)
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
