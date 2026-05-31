package maintenance

import (
	"os"
	"path/filepath"
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
