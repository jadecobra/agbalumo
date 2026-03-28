package agent

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckSQLi(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Safe query with placeholder",
			code: `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = ?", id)
}`,
			expected: 0,
		},
		{
			name: "Unsafe query with concatenation",
			code: `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = " + id)
}`,
			expected: 1,
		},
		{
			name: "Unsafe query with fmt.Sprintf",
			code: `package main
import "fmt"
func main() {
	db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}
			violations := checkSQLi(f, fset)
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
			}
		})
	}
}

func TestCheckXSS(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Safe template.HTML with literal",
			code: `package main
import "html/template"
func main() {
	_ = template.HTML("<div>Static</div>")
}`,
			expected: 0,
		},
		{
			name: "Unsafe template.HTML with variable",
			code: `package main
import "html/template"
func main() {
	_ = template.HTML(userInput)
}`,
			expected: 1,
		},
		{
			name: "Unsafe Echo c.HTML with raw string",
			code: `package main
func Handler(c echo.Context) error {
	return c.HTML(200, "<div>" + userInput + "</div>")
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}
			violations := checkXSS(f, fset)
			if len(violations) != tt.expected {
				t.Errorf("Expected %d violations, got %d", tt.expected, len(violations))
			}
		})
	}
}
