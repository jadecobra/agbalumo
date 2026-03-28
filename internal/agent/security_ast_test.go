package agent

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInsecurePatternsGo(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Insecure onclick in backticks",
			code: `package main
func main() {
	html := ` + "`" + `<button onclick="alert(1)">Click</button>` + "`" + `
}`,
			expected: 1,
		},
		{
			name: "Dangerous JS eval in string literal",
			code: `package main
func main() {
	js := "eval('window.x = 1')"
}`,
			expected: 1,
		},
		{
			name: "Forbidden CDN in string literal",
			code: `package main
func main() {
	url := "https://unpkg.com/htmx.org"
}`,
			expected: 1,
		},
		{
			name: "Ignored onclick with comment",
			code: `package main
func main() {
	// #nosec - testing AST exclusion
	html := ` + "`" + `<button onclick="alert(1)">Click</button>` + "`" + `
}`,
			expected: 0,
		},
		{
			name: "Safe string with no patterns",
			code: `package main
func main() {
	msg := "Hello World"
}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}
			violations := checkInsecurePatternsGo(f, fset)
			assert.Equal(t, tt.expected, len(violations), "Should find %d violations in %s", tt.expected, tt.name)
		})
	}
}
