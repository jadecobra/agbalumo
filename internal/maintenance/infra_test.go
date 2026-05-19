package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyCITools(t *testing.T) {
	tests := []struct {
		name        string
		ymlContent  string
		errorMsg    string
		expectError bool
	}{
		{
			name: "Valid Trivy and Flyctl with retry loop",
			ymlContent: `
name: CI/CD Pipeline
jobs:
  deploy:
    steps:
      - name: Setup Flyctl
        uses: superfly/flyctl-actions/setup-flyctl@v1
      - name: Deploy
        run: |
          for i in {1..3}; do
            flyctl deploy --remote-only && exit 0
            sleep 15
          done
          exit 1
      - name: Run Trivy
        uses: aquasecurity/trivy-action@v1
        with:
          version: 'v0.70.0'
`,
			errorMsg:    "",
			expectError: false,
		},
		{
			name: "Invalid Trivy input trivy-version",
			ymlContent: `
name: CI/CD Pipeline
jobs:
  test:
    steps:
      - name: Run Trivy
        uses: aquasecurity/trivy-action@v1
        with:
          trivy-version: 'v0.70.0'
`,
			errorMsg:    "invalid CI configuration: 'aquasecurity/trivy-action' uses 'version', not 'trivy-version'",
			expectError: true,
		},
		{
			name: "Flyctl deploy without retry loop",
			ymlContent: `
name: CI/CD Pipeline
jobs:
  deploy:
    steps:
      - name: Setup Flyctl
        uses: superfly/flyctl-actions/setup-flyctl@v1
      - name: Deploy
        run: flyctl deploy --remote-only
`,
			errorMsg:    "invalid CI configuration: 'flyctl deploy' must implement automated retry logic with exponential backoff to handle platform 503 errors",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runVerifyCIToolsTest(t, tt.ymlContent, tt.expectError, tt.errorMsg)
		})
	}
}

func runVerifyCIToolsTest(t *testing.T, ymlContent string, expectError bool, errorMsg string) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "verify_ci_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	wfDir := filepath.Join(tempDir, ".github/workflows")
	if errMk := os.MkdirAll(wfDir, 0700); errMk != nil {
		t.Fatalf("failed to create workflows dir: %v", errMk)
	}

	ciFile := filepath.Join(wfDir, "ci.yml")
	if errWrite := os.WriteFile(ciFile, []byte(ymlContent), 0600); errWrite != nil {
		t.Fatalf("failed to write ci.yml: %v", errWrite)
	}

	err = VerifyCITools(tempDir)
	if expectError {
		if err == nil {
			t.Errorf("expected error but got none")
		} else if errorMsg != "" && err.Error() != errorMsg {
			t.Errorf("expected error %q, got %q", errorMsg, err.Error())
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
