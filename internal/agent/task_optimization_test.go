package agent

import (
	"os"
	"strings"
	"testing"
)

func TestTaskfileToolingOptimization(t *testing.T) {
	// Find Taskfile.yml by searching up from the current package dir
	filePath := "Taskfile.yml"
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(filePath); err == nil {
			break
		}
		filePath = "../" + filePath
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read Taskfile.yml at %s: %v", filePath, err)
	}
	taskfile := string(content)

	tests := []struct {
		name           string
		expectedStatus string
		expectedCmd    string
	}{
		{
			name:           "vulncheck",
			expectedStatus: "test -f ./.tester/tmp/go/bin/govulncheck",
			expectedCmd:    "{{.TASKFILE_DIR}}/.tester/tmp/go/bin/govulncheck ./...",
		},
		{
			name:           "lint",
			expectedStatus: "test -f ./.tester/tmp/go/bin/golangci-lint",
			expectedCmd:    "{{.TASKFILE_DIR}}/.tester/tmp/go/bin/golangci-lint run",
		},
		{
			name:           "gitleaks",
			expectedStatus: "test -f ./.tester/tmp/go/bin/gitleaks",
			expectedCmd:    "{{.TASKFILE_DIR}}/.tester/tmp/go/bin/gitleaks",
		},
		{
			name:           "gocognit",
			expectedStatus: "test -f ./.tester/tmp/go/bin/gocognit",
			expectedCmd:    "go install github.com/uudashr/gocognit/cmd/gocognit@v1.0.1",
		},
		{
			name:           "dupl",
			expectedStatus: "test -f ./.tester/tmp/go/bin/dupl",
			expectedCmd:    "go install github.com/mibk/dupl@v1.0.0",
		},
		{
			name:           "goconst",
			expectedStatus: "test -f ./.tester/tmp/go/bin/goconst",
			expectedCmd:    "go install github.com/jgautheron/goconst/cmd/goconst@v1.7.1",
		},
		{
			name:           "fieldalignment",
			expectedStatus: "test -f ./.tester/tmp/go/bin/fieldalignment",
			expectedCmd:    "go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.30.0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(taskfile, tc.expectedStatus) {
				t.Errorf("Taskfile.yml does not contain expected status check: %s", tc.expectedStatus)
			}
			if !strings.Contains(taskfile, tc.expectedCmd) {
				t.Errorf("Taskfile.yml does not contain expected command: %s", tc.expectedCmd)
			}
		})
	}
}
