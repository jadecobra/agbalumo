package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSandboxParity(t *testing.T) {
	tests := []struct {
		name           string
		components     string
		sandbox        string
		wantViolations int
	}{
		{
			name: "Match - no violations",
			components: `
				{{ define "btn_close" }}<button>Close</button>{{ end }}
				{{ define "card_listing" }}<div>Listing</div>{{ end }}
			`,
			sandbox: `
				{{ template "btn_close" . }}
				{{ template "card_listing" . }}
			`,
			wantViolations: 0,
		},
		{
			name: "Missing reference - violation",
			components: `
				{{ define "btn_close" }}<button>Close</button>{{ end }}
				{{ define "new_badge" }}<span>New</span>{{ end }}
			`,
			sandbox: `
				{{ template "btn_close" . }}
			`,
			wantViolations: 1,
		},
		{
			name:           "Empty files - no violations",
			components:     "",
			sandbox:        "",
			wantViolations: 0,
		},
		{
			name: "Layout blocks - excluded from check",
			components: `
				{{ define "base.html" }}{{ end }}
				{{ define "title" }}{{ end }}
				{{ define "content" }}{{ end }}
				{{ define "head" }}{{ end }}
				{{ define "scripts" }}{{ end }}
				{{ define "real_component" }}{{ end }}
			`,
			sandbox: `
				{{ template "real_component" . }}
			`,
			wantViolations: 0,
		},
		{
			name:       "Raw button without sandbox launch testid - violation",
			components: "",
			sandbox: `
				<button class="px-6 py-3 bg-earth-dark text-white">Click me</button>
			`,
			wantViolations: 1,
		},
		{
			name:       "Raw listing-card class - violation",
			components: "",
			sandbox: `
				<div class="listing-card card-juicy">Some Listing</div>
			`,
			wantViolations: 1,
		},
		{
			name:       "Allowed sandbox launch button - no violations",
			components: "",
			sandbox: `
				<button data-testid="sandbox-launch-feedback">Feedback</button>
				<button data-testid="sandbox-launch-login-prompt">Login</button>
			`,
			wantViolations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, partialsDir, templatesDir := prepareTestFS(t)

			writeTestFile(t, partialsDir, "ui_components.html", tt.components)
			writeTestFile(t, templatesDir, "sandbox.html", tt.sandbox)

			violations, _ := CheckSandboxParity(tmpDir)
			if len(violations) != tt.wantViolations {
				t.Errorf("CheckSandboxParity() got %d violations, want %d (violations: %+v)", len(violations), tt.wantViolations, violations)
			}
		})
	}
}

func TestCheckSandboxParity_Overhaul_MultiScan(t *testing.T) {
	tmpDir, partialsDir, templatesDir := prepareTestFS(t)

	// Create components directory
	componentsDir := filepath.Join(tmpDir, "ui", "templates", "components")
	if err := os.MkdirAll(componentsDir, 0700); err != nil {
		t.Fatalf("failed to create components dir: %v", err)
	}

	// 1. Multi-file scan in partials
	writeTestFile(t, partialsDir, "partial1.html", `{{ define "comp_partial1" }}<div>1</div>{{ end }}`)
	writeTestFile(t, partialsDir, "partial2.html", `{{ define "comp_partial2" }}<div>2</div>{{ end }}`)

	// 2. Components directory scan
	writeTestFile(t, componentsDir, "comp1.html", `{{ define "comp_comp1" }}<div>3</div>{{ end }}`)

	// 3. Excluded component (should be ignored)
	writeTestFile(t, partialsDir, "exclude.html", `{{ define "navigation" }}<nav></nav>{{ end }}`)

	// Sandbox only references comp_partial1
	sandboxContent := `
		{{ template "comp_partial1" . }}
	`
	writeTestFile(t, templatesDir, "sandbox.html", sandboxContent)

	violations, err := CheckSandboxParity(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// We expect two violations: for comp_partial2 and comp_comp1.
	// navigation is excluded.
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d: %+v", len(violations), violations)
	}

	foundPartial2 := false
	foundComp1 := false
	for _, v := range violations {
		if v.Component == "comp_partial2" {
			foundPartial2 = true
		}
		if v.Component == "comp_comp1" {
			foundComp1 = true
		}
	}
	if !foundPartial2 {
		t.Errorf("expected violation for 'comp_partial2' not found")
	}
	if !foundComp1 {
		t.Errorf("expected violation for 'comp_comp1' not found")
	}
}

func TestCheckSandboxParity_Overhaul_HiddenDiv(t *testing.T) {
	tmpDir, partialsDir, templatesDir := prepareTestFS(t)

	writeTestFile(t, partialsDir, "ui_components.html", `{{ define "my_hidden_component" }}<div>Hidden</div>{{ end }}`)

	// Template is referenced inside a hidden div
	sandboxContent := `
		<div class="hidden">
			{{ template "my_hidden_component" . }}
		</div>
	`
	writeTestFile(t, templatesDir, "sandbox.html", sandboxContent)

	violations, err := CheckSandboxParity(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// We expect 1 violation because the component is referenced inside a hidden div
	foundHiddenViolation := false
	for _, v := range violations {
		if v.Component == "my_hidden_component" && strings.Contains(v.Message, "hidden") {
			foundHiddenViolation = true
		}
	}
	if !foundHiddenViolation {
		t.Errorf("expected violation for my_hidden_component in hidden div, violations: %+v", violations)
	}
}

func prepareTestFS(t *testing.T) (string, string, string) {
	tmpDir := t.TempDir()
	partialsDir := filepath.Join(tmpDir, "ui", "templates", "partials")
	templatesDir := filepath.Join(tmpDir, "ui", "templates")

	if err := os.MkdirAll(partialsDir, 0700); err != nil {
		t.Fatalf("failed to create partials dir: %v", err)
	}
	if err := os.MkdirAll(templatesDir, 0700); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}
	return tmpDir, partialsDir, templatesDir
}
