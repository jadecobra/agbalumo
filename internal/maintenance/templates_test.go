package maintenance

import (
	"path/filepath"
	"testing"
)

func TestTemplateUtilities(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t, "template_utils")
	defer cleanup()

	tests := []struct {
		testFn   func(string) ([]string, error)
		expected map[string]bool
		name     string
		filename string
		content  string
	}{
		{
			name:     "ExtractRendererFunctions",
			filename: "renderer.go",
			content: `
package ui
func GetFuncs() {
	funcs := map[string]interface{}{
		"formatDate": fmt.Sprintf,
		"truncate":   nil,
	}
}
`,
			testFn:   ExtractRendererFunctions,
			expected: map[string]bool{"formatDate": true, "truncate": true},
		},
		{
			name:     "ExtractTemplateFunctionCalls",
			filename: "index.html",
			content: `
<div>{{ formatDate .Date }}</div>
<div>{{ range .Items | filterItems }}</div>
`,
			testFn: func(path string) ([]string, error) {
				return ExtractTemplateFunctionCalls(filepath.Dir(path))
			},
			expected: map[string]bool{"formatDate": true, "filterItems": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestFile(t, tmpDir, tt.filename, tt.content)
			res, err := tt.testFn(path)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			assertStringsMatch(t, tt.name, res, tt.expected)
		})
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
