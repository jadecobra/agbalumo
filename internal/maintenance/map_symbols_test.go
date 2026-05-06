package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSymbolMap(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "types and funcs",
			content: `package test
type MyStruct struct {}
func MyFunc() {}
func privateFunc() {}
`,
			expected: []string{
				"[Type] MyStruct -> test.go",
				"[Func] MyFunc -> test.go",
			},
		},
		{
			name: "method pointer receiver",
			content: `package test
type Handler struct {}
func (h *Handler) ServeHTTP() {}
`,
			expected: []string{
				"[Type] Handler -> test.go",
				"[Method] (Handler).ServeHTTP -> test.go",
			},
		},
		{
			name: "method value receiver",
			content: `package test
type Handler struct {}
func (h Handler) Get() {}
`,
			expected: []string{
				"[Type] Handler -> test.go",
				"[Method] (Handler).Get -> test.go",
			},
		},
		{
			name: "interface type",
			content: `package test
type Store interface {
	Save() error
}
`,
			expected: []string{
				"[Interface] Store -> test.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			goPath := filepath.Join(tmpDir, "test.go")
			if err := os.WriteFile(goPath, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}

			got := GenerateSymbolMap(tmpDir)

			for _, want := range tt.expected {
				if !strings.Contains(got, want) {
					t.Errorf("GenerateSymbolMap() missing symbol %q\nGot:\n%s", want, got)
				}
			}

			if strings.Contains(got, "privateFunc") {
				t.Errorf("GenerateSymbolMap() should not contain privateFunc")
			}
		})
	}
}
