package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDirs(t *testing.T, tmpDir string) (string, string) {
	t.Helper()
	uiDir := filepath.Join(tmpDir, "ui", "templates")
	if err := os.MkdirAll(uiDir, 0750); err != nil {
		t.Fatal(err)
	}
	moduleDir := filepath.Join(tmpDir, "internal", "module")
	if err := os.MkdirAll(moduleDir, 0750); err != nil {
		t.Fatal(err)
	}
	return uiDir, moduleDir
}

type testCase struct {
	name       string
	template   string
	goCode     string
	violations int
}

func runContractTest(t *testing.T, tt testCase, tmpDir, uiDir, moduleDir string) {
	_ = os.RemoveAll(uiDir)
	_ = os.RemoveAll(moduleDir)
	setupTestDirs(t, tmpDir)

	if err := os.WriteFile(filepath.Join(uiDir, "test.html"), []byte(tt.template), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "types.go"), []byte(tt.goCode), 0600); err != nil {
		t.Fatal(err)
	}

	violations, err := CheckTemplateContracts(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(violations) != tt.violations {
		t.Errorf("expected %d violations, got %d: %v", tt.violations, len(violations), violations)
	}
}

func TestCheckTemplateContracts(t *testing.T) {
	tmpDir := t.TempDir()
	uiDir, moduleDir := setupTestDirs(t, tmpDir)

	tests := []testCase{
		{
			name:       "matching contract",
			template:   `{{ define "test" }}<h1>{{ .Title }}</h1>{{ end }}`,
			goCode:     "package module\ntype TestViewModel struct {\n\tTitle string\n}",
			violations: 0,
		},
		{
			name:       "missing field",
			template:   `{{ define "test" }}<h1>{{ .NonExistent }}</h1>{{ end }}`,
			goCode:     "package module\ntype TestViewModel struct {\n\tTitle string\n}",
			violations: 1,
		},
		{
			name:       "inherited field",
			template:   `{{ define "test" }}<h1>{{ .Env }}</h1>{{ end }}`,
			goCode:     "package module\ntype BaseViewData struct {\n\tEnv string\n}\ntype TestViewModel struct {\n\tBaseViewData\n\tTitle string\n}",
			violations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runContractTest(t, tt, tmpDir, uiDir, moduleDir)
		})
	}
}
