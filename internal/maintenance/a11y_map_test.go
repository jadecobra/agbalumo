package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func setupA11yTest(t *testing.T, errorContent, templateHTML string) string {
	t.Helper()
	tmpDir := t.TempDir()
	resultsDir := filepath.Join(tmpDir, "test-results")
	templatesDir := filepath.Join(tmpDir, "ui", "templates")

	if err := os.MkdirAll(resultsDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(templatesDir, 0750); err != nil {
		t.Fatal(err)
	}

	if errorContent != "" {
		errFile := filepath.Join(resultsDir, "error-context.md")
		if err := os.WriteFile(errFile, []byte(errorContent), 0600); err != nil {
			t.Fatal(err)
		}
	}

	if templateHTML != "" {
		templateFile := filepath.Join(templatesDir, "button.html")
		if err := os.WriteFile(templateFile, []byte(templateHTML), 0600); err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func assertViolation(t *testing.T, v A11yViolationMap, wantViolation, wantFile, wantSuggestion string, wantLine int) {
	t.Helper()
	if v.ViolationID != wantViolation {
		t.Errorf("ViolationID = %q, want %q", v.ViolationID, wantViolation)
	}
	if v.TemplateFile != wantFile {
		t.Errorf("TemplateFile = %q, want %q", v.TemplateFile, wantFile)
	}
	if v.Line != wantLine {
		t.Errorf("Line = %d, want %d", v.Line, wantLine)
	}
	if v.FixSuggestion != wantSuggestion {
		t.Errorf("FixSuggestion = %q, want %q", v.FixSuggestion, wantSuggestion)
	}
}

func TestMapA11yViolations_Success(t *testing.T) {
	tests := []struct {
		name           string
		errorContent   string
		templateHTML   string
		wantViolation  string
		wantFile       string
		wantSuggestion string
		wantLine       int
	}{
		{
			name:           "maps violation to template",
			errorContent:   "A11y Violations: " + `[{"id":"color-contrast","impact":"serious","nodes":[{"html":"<button class='bg-red-500'>Click</button>","failureSummary":"Fix contrast"}]}]`,
			templateHTML:   "<div>\n\t<button class='bg-red-500'>Click</button>\n</div>",
			wantViolation:  "color-contrast",
			wantFile:       "ui/templates/button.html",
			wantLine:       2,
			wantSuggestion: "Run `go run ./cmd/verify design` for static analysis.",
		},
		{
			name:           "label violation with fix suggestion",
			errorContent:   `A11y Violations: [{"id":"label","impact":"critical","nodes":[{"html":"<input type='text'>","failureSummary":"No label"}]}]`,
			templateHTML:   "<form>\n<input type='text'>\n</form>",
			wantViolation:  "label",
			wantFile:       "ui/templates/button.html",
			wantLine:       2,
			wantSuggestion: "Use <label for='ID'> or aria-label to associate text with the input.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupA11yTest(t, tt.errorContent, tt.templateHTML)
			violations, err := MapA11yViolations(tmpDir)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(violations) != 1 {
				t.Fatalf("got %d violations, want 1", len(violations))
			}

			assertViolation(t, violations[0], tt.wantViolation, tt.wantFile, tt.wantSuggestion, tt.wantLine)
		})
	}
}

func TestMapA11yViolations_Errors(t *testing.T) {
	tests := []struct {
		name         string
		errorContent string
	}{
		{
			name:         "malformed JSON",
			errorContent: "A11y Violations: [{broken json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupA11yTest(t, tt.errorContent, "")
			_, err := MapA11yViolations(tmpDir)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestMapA11yViolations_Empty(t *testing.T) {
	tests := []struct {
		name         string
		errorContent string
	}{
		{
			name:         "missing error-context.md",
			errorContent: "",
		},
		{
			name:         "no violations marker",
			errorContent: "Some random test output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupA11yTest(t, tt.errorContent, "")
			violations, err := MapA11yViolations(tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(violations) != 0 {
				t.Errorf("got %d violations, want 0", len(violations))
			}
		})
	}
}
