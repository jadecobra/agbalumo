package sql

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSQLi(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "Potential SQLi in Query",
			code: `package main
func main() {
	db.Query("SELECT * FROM users WHERE id = " + id)
}`,
			expected: 1,
		},
		{
			name: "Safe Query with literal",
			code: `package main
func main() {
	db.Query("SELECT * FROM users")
}`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, "test.go", tt.code, 0)
			violations := CheckSQLi(f, fset)
			assert.Equal(t, tt.expected, len(violations))
		})
	}
}
