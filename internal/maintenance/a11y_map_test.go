package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMapA11yViolations(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Setup test-results directory and a mock error-context.md
	resultsDir := filepath.Join(tmpDir, "test-results")
	err := os.MkdirAll(resultsDir, 0750)
	if err != nil {
		t.Fatal(err)
	}

	// Mock Axe JSON output
	axeJSON := `[{"id":"color-contrast","impact":"serious","nodes":[{"html":"<button class='bg-red-500'>Click</button>","failureSummary":"Fix contrast"}]}]`
	errorContent := "Some other text\nA11y Violations: " + axeJSON + "\nTrailing text"

	errFile := filepath.Join(resultsDir, "error-context.md")
	if err = os.WriteFile(errFile, []byte(errorContent), 0600); err != nil {
		t.Fatal(err)
	}

	// 2. Setup mock template file
	templatesDir := filepath.Join(tmpDir, "ui", "templates")
	if err = os.MkdirAll(templatesDir, 0750); err != nil {
		t.Fatal(err)
	}
	templateFile := filepath.Join(templatesDir, "button.html")
	templateContent := `<div>
	<button class='bg-red-500'>Click</button>
</div>`
	if err = os.WriteFile(templateFile, []byte(templateContent), 0600); err != nil {
		t.Fatal(err)
	}

	// 3. Run MapA11yViolations
	// Note: This will fail to compile initially because MapA11yViolations is not defined.
	// But TDD requires a red state.
	violations, err := MapA11yViolations(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(violations) != 1 {
		t.Fatalf("Expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.ViolationID != "color-contrast" {
		t.Errorf("Expected violation color-contrast, got %s", v.ViolationID)
	}
	if v.TemplateFile != "ui/templates/button.html" {
		t.Errorf("Expected template button.html, got %s", v.TemplateFile)
	}
	if v.Line != 2 {
		t.Errorf("Expected line 2, got %d", v.Line)
	}
}
