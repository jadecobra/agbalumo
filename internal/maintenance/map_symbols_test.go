package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSymbolMap(t *testing.T) {
	t.Run("Types and Funcs", func(t *testing.T) {
		content := `package test
type MyStruct struct {}
func MyFunc() {}
func privateFunc() {}
`
		expected := []string{"[Type] MyStruct -> test.go", "[Func] MyFunc -> test.go"}
		verifySymbols(t, content, expected, "privateFunc")
	})

	t.Run("Method Receivers", func(t *testing.T) {
		content := `package test
type Handler struct {}
func (h *Handler) ServeHTTP() {}
func (h Handler) Get() {}
`
		expected := []string{"[Type] Handler -> test.go", "[Method] (Handler).ServeHTTP -> test.go", "[Method] (Handler).Get -> test.go"}
		verifySymbols(t, content, expected, "")
	})

	t.Run("Interfaces", func(t *testing.T) {
		content := `package test
type Store interface { Save() error }
`
		expected := []string{"[Interface] Store -> test.go"}
		verifySymbols(t, content, expected, "")
	})
}

func verifySymbols(t *testing.T, content string, expected []string, unexpected string) {
	t.Helper()
	tmpDir := t.TempDir()
	goPath := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got := GenerateSymbolMap(tmpDir)

	for _, want := range expected {
		if !strings.Contains(got, want) {
			t.Errorf("GenerateSymbolMap() missing symbol %q\nGot:\n%s", want, got)
		}
	}

	if unexpected != "" && strings.Contains(got, unexpected) {
		t.Errorf("GenerateSymbolMap() should not contain %q", unexpected)
	}
}
