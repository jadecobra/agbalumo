package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyPlaywrightVersionParity(t *testing.T) {
	tests := []struct {
		name        string
		packageJSON string
		ciRunner    string
		wantErr     bool
	}{
		{
			name: "exact match",
			packageJSON: `{
				"devDependencies": {
					"@playwright/test": "1.59.1"
				}
			}`,
			ciRunner: `
				"mcr.microsoft.com/playwright:v1.59.1-noble"
			`,
			wantErr: false,
		},
		{
			name: "caret match",
			packageJSON: `{
				"devDependencies": {
					"@playwright/test": "^1.59.1"
				}
			}`,
			ciRunner: `
				"mcr.microsoft.com/playwright:v1.59.1-noble"
			`,
			wantErr: false,
		},
		{
			name: "mismatch major minor",
			packageJSON: `{
				"devDependencies": {
					"@playwright/test": "1.59.1"
				}
			}`,
			ciRunner: `
				"mcr.microsoft.com/playwright:v1.52.0-noble"
			`,
			wantErr: true,
		},
		{
			name:        "missing package.json",
			packageJSON: "",
			ciRunner:    "any",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			setupPlaywrightVersionFixtures(t, tmpDir, tt.packageJSON, tt.ciRunner)

			err := VerifyPlaywrightVersionParity(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPlaywrightVersionParity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupPlaywrightVersionFixtures(t *testing.T, tmpDir, packageJSON, ciRunner string) {
	t.Helper()
	if packageJSON != "" {
		err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0600)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create internal/maintenance/ci_runner.go relative to tmpDir
	maintenanceDir := filepath.Join(tmpDir, "internal", "maintenance")
	err := os.MkdirAll(maintenanceDir, 0750)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(maintenanceDir, "ci_runner.go"), []byte(ciRunner), 0600)
	if err != nil {
		t.Fatal(err)
	}
}
