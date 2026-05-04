package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSymbolMap(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dummy .go file
	goContent := `package test
type MyStruct struct {}
func MyFunc() {}
func (s *MyStruct) MyMethod() {}
func privateFunc() {}
`
	goPath := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a test file (should be ignored)
	testGoContent := `package test
func TestSomething(t *testing.T) {}
`
	testGoPath := filepath.Join(tmpDir, "test_test.go")
	if err := os.WriteFile(testGoPath, []byte(testGoContent), 0600); err != nil {
		t.Fatal(err)
	}

	got := GenerateSymbolMap(tmpDir)

	wantSymbols := []string{
		"[Type] MyStruct -> test.go",
		"[Func] MyFunc -> test.go",
		"[Func] MyMethod -> test.go",
	}

	for _, s := range wantSymbols {
		if !strings.Contains(got, s) {
			t.Errorf("GenerateSymbolMap() missing symbol %q\nGot:\n%s", s, got)
		}
	}

	if strings.Contains(got, "privateFunc") {
		t.Errorf("GenerateSymbolMap() should not contain privateFunc")
	}

	if strings.Contains(got, "TestSomething") {
		t.Errorf("GenerateSymbolMap() should not contain TestSomething")
	}
}

func TestGenerateTemplateMap(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dummy .html file
	htmlContent := `{{ define "my_partial" }}
  <div>Hello</div>
{{ end }}
`
	htmlPath := filepath.Join(tmpDir, "template.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0600); err != nil {
		t.Fatal(err)
	}

	got := GenerateTemplateMap(tmpDir)

	want := "[Template] my_partial -> template.html"
	if !strings.Contains(got, want) {
		t.Errorf("GenerateTemplateMap() missing %q\nGot:\n%s", want, got)
	}
}
