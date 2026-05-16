package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMapA11yViolations(t *testing.T) {
	tests := []struct {
		name           string
		errorContent   string
		templateHTML   string
		wantViolation  string
		wantFile       string
		wantSuggestion string
		wantCount      int
		wantLine       int
		wantErr        bool
	}{
		{
			name:           "parse success - maps violation to template",
			errorContent:   "Some text\nA11y Violations: " + `[{"id":"color-contrast","impact":"serious","nodes":[{"html":"<button class='bg-red-500'>Click</button>","failureSummary":"Fix contrast"}]}]` + "\nTrailing text",
			templateHTML:   "<div>\n\t<button class='bg-red-500'>Click</button>\n</div>",
			wantCount:      1,
			wantViolation:  "color-contrast",
			wantFile:       "ui/templates/button.html",
			wantLine:       2,
			wantSuggestion: "Run `go run ./cmd/verify design` for static analysis.",
		},
		{
			name:           "parse success - label violation with fix suggestion",
			errorContent:   `A11y Violations: [{"id":"label","impact":"critical","nodes":[{"html":"<input type='text'>","failureSummary":"No label"}]}]`,
			templateHTML:   "<form>\n<input type='text'>\n</form>",
			wantCount:      1,
			wantViolation:  "label",
			wantFile:       "ui/templates/button.html",
			wantLine:       2,
			wantSuggestion: "Use <label for='ID'> or aria-label to associate text with the input.",
		},
		{
			name:         "parse failure - malformed JSON returns error",
			errorContent: "A11y Violations: [{broken json",
			wantErr:      true,
		},
		{
			name:         "missing error-context.md - returns nil",
			errorContent: "", // signals: don't create the file
			wantCount:    0,
		},
		{
			name:         "no violations marker - returns nil",
			errorContent: "Some random test output with no Axe data",
			wantCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			resultsDir := filepath.Join(tmpDir, "test-results")
			templatesDir := filepath.Join(tmpDir, "ui", "templates")

			if err := os.MkdirAll(resultsDir, 0750); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(templatesDir, 0750); err != nil {
				t.Fatal(err)
			}

			// Only create error-context.md if errorContent is non-empty
			if tt.errorContent != "" {
				errFile := filepath.Join(resultsDir, "error-context.md")
				if err := os.WriteFile(errFile, []byte(tt.errorContent), 0600); err != nil {
					t.Fatal(err)
				}
			}

			// Create template file if provided
			if tt.templateHTML != "" {
				templateFile := filepath.Join(templatesDir, "button.html")
				if err := os.WriteFile(templateFile, []byte(tt.templateHTML), 0600); err != nil {
					t.Fatal(err)
				}
			}

			violations, err := MapA11yViolations(tmpDir)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(violations) != tt.wantCount {
				t.Fatalf("got %d violations, want %d", len(violations), tt.wantCount)
			}

			if tt.wantCount > 0 {
				v := violations[0]
				if v.ViolationID != tt.wantViolation {
					t.Errorf("ViolationID = %q, want %q", v.ViolationID, tt.wantViolation)
				}
				if v.TemplateFile != tt.wantFile {
					t.Errorf("TemplateFile = %q, want %q", v.TemplateFile, tt.wantFile)
				}
				if v.Line != tt.wantLine {
					t.Errorf("Line = %d, want %d", v.Line, tt.wantLine)
				}
				if v.FixSuggestion != tt.wantSuggestion {
					t.Errorf("FixSuggestion = %q, want %q", v.FixSuggestion, tt.wantSuggestion)
				}
			}
		})
	}
}
