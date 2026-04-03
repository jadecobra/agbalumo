package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractRendererFunctions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "renderer_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	rendererCode := `
package ui
func GetFuncs() {
	funcs := map[string]interface{}{
		"formatDate": fmt.Sprintf,
		"truncate":   nil,
	}
}
`
	rendererPath := filepath.Join(tmpDir, "renderer.go")
	_ = os.WriteFile(rendererPath, []byte(rendererCode), 0644)

	funcs, err := ExtractRendererFunctions(rendererPath)
	if err != nil {
		t.Fatalf("ExtractRendererFunctions failed: %v", err)
	}

	expected := map[string]bool{"formatDate": true, "truncate": true}
	for _, f := range funcs {
		delete(expected, f)
	}

	if len(expected) > 0 {
		t.Errorf("missing expected functions: %v", expected)
	}
}

func TestExtractTemplateFunctionCalls(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templateCode := `
<div>{{ formatDate .Date }}</div>
<div>{{ range .Items | filterItems }}</div>
`
	_ = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(templateCode), 0644)

	used, err := ExtractTemplateFunctionCalls(tmpDir)
	if err != nil {
		t.Fatalf("ExtractTemplateFunctionCalls failed: %v", err)
	}

	expected := map[string]bool{"formatDate": true, "filterItems": true}
	for _, u := range used {
		delete(expected, u)
	}

	if len(expected) > 0 {
		t.Errorf("missing expected template function calls: %v", expected)
	}
}

func TestCheckTemplateDrift(t *testing.T) {
	defined := []string{"formatDate"}
	used := []string{"formatDate", "unknownFunc", "ExportedType"}

	drifts := CheckTemplateDrift(defined, used)

	if len(drifts) != 1 {
		t.Errorf("expected 1 drift, got %d", len(drifts))
	}

	if len(drifts) > 0 && drifts[0] != "Undefined template function used: 'unknownFunc'" {
		t.Errorf("unexpected drift message: %s", drifts[0])
	}
}
