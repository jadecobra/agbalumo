package agent

import (
	"os"
	"testing"
)

func TestVerifyTemplateDrift(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.MkdirAll(tmpDir+"/internal/ui", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/ui/templates", 0755)
	if err != nil {
		t.Fatal(err)
	}

	rendererContent := `package ui
func NewRenderer() {
	funcMap := template.FuncMap{
		"func1": func() {},
		"func2": func() {},
	}
}`
	_ = os.WriteFile(tmpDir+"/internal/ui/renderer.go", []byte(rendererContent), 0644)

	templateContent := `<div>{{ func1 . }}</div>`
	_ = os.WriteFile(tmpDir+"/ui/templates/index.html", []byte(templateContent), 0644)

	// Change dir to TempDir
	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Case 1: Success
	if !VerifyTemplateDrift() {
		t.Error("VerifyTemplateDrift failed on valid sync")
	}

	// Case 2: Drift
	templateContentDrift := `<div>{{ unknownFunc . }}</div>`
	_ = os.WriteFile(tmpDir+"/ui/templates/index.html", []byte(templateContentDrift), 0644)
	if VerifyTemplateDrift() {
		t.Error("VerifyTemplateDrift passed on drift")
	}
}

func TestExtractTemplateFunctionCalls_Complex(t *testing.T) {
	tmpDir := t.TempDir()
	content := "\n\t\t<div>{{ upper .Name }}</div>\n\t\t<div>{{ lower (trim .Value) }}</div>\n\t\t<div>{{ unknown . }}</div>\n\t"
	_ = os.WriteFile(tmpDir+"/test.html", []byte(content), 0644)

	calls, err := ExtractTemplateFunctionCalls(tmpDir)
	if err != nil {
		t.Fatalf("ExtractTemplateFunctionCalls failed: %v", err)
	}

	expected := map[string]bool{
		"upper":   true,
		"lower":   true,
		"unknown": true,
	}

	// Check that we found at least these
	found := make(map[string]bool)
	for _, c := range calls {
		found[c] = true
	}

	for exp := range expected {
		if !found[exp] {
			t.Errorf("expected call %s not found", exp)
		}
	}
}
