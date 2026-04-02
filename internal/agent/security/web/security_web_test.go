package web

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckXSS(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Potential XSS in template.HTML",
			code: `package main
func main() {
	template.HTML(unsafe_var)
}`,
			expected: 1,
		},
		{
			name: "Safe template.HTML with literal",
			code: `package main
func main() {
	template.HTML("<b>Safe</b>")
}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, "test.go", tt.code, 0)
			violations := CheckXSS(f, fset)
			assert.Equal(t, tt.expected, len(violations))
		})
	}
}

func TestCheckSSRF(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Potential SSRF in http.Get",
			code: `package main
func main() {
	http.Get("http://example.com/" + path)
}`,
			expected: 1,
		},
		{
			name: "Safe http.Get with literal",
			code: `package main
func main() {
	http.Get("https://example.com")
}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, "test.go", tt.code, 0)
			violations := CheckSSRF(f, fset)
			assert.Equal(t, tt.expected, len(violations))
		})
	}
}
