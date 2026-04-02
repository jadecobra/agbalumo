package sectypes

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateEntropy(t *testing.T) {
	t.Run("LowEntropy", func(t *testing.T) {
		assert.Equal(t, 0.0, CalculateEntropy("aaaaa"))
	})
	t.Run("HighEntropy", func(t *testing.T) {
		e := CalculateEntropy("abc123XYZ!@#")
		assert.True(t, e > 3.0)
	})
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, 0.0, CalculateEntropy(""))
	})
}

func TestDeduplicateViolations(t *testing.T) {
	v1 := SecurityViolation{File: "f1", Line: 1, Message: "m1"}
	v2 := SecurityViolation{File: "f1", Line: 1, Message: "m1"}
	v3 := SecurityViolation{File: "f2", Line: 2, Message: "m2"}

	list := []SecurityViolation{v1, v2, v3}
	deduped := DeduplicateViolations(list)

	assert.Len(t, deduped, 2)
}

func TestIsUnsafeString(t *testing.T) {
	fset := token.NewFileSet()

	t.Run("SafeLiteral", func(t *testing.T) {
		expr, _ := parser.ParseExpr("\"fixed\"")
		assert.False(t, IsUnsafeString(expr))
	})

	t.Run("UnsafeConcat", func(t *testing.T) {
		// Go strings like `abc` + "def" are BinaryExpr
		expr, _ := parser.ParseExpr("`abc` + \"def\"")
		assert.True(t, IsUnsafeString(expr))
	})

	t.Run("UnsafeSprintf", func(t *testing.T) {
		expr, _ := parser.ParseExpr("fmt.Sprintf(\"SELECT * FROM %s\", table)")
		assert.True(t, IsUnsafeString(expr))
	})

	t.Run("UnsafeVar", func(t *testing.T) {
		expr, _ := parser.ParseExpr("query")
		assert.True(t, IsUnsafeString(expr))
	})

	_ = fset
}

func TestIsIgnored(t *testing.T) {
	t.Run("IsIgnoredRaw", func(t *testing.T) {
		assert.False(t, IsIgnoredRaw("var x = 1"))
		assert.True(t, IsIgnoredRaw("// #nosec - test rationale"))
		assert.True(t, IsIgnoredRaw("// antigravity:allow"))
	})

	t.Run("IsIgnoredAST", func(t *testing.T) {
		fset := token.NewFileSet()
		// Try a simpler way to trigger #nosec
		src := "package p\n// #nosec - test rationale\nfunc main() {}"
		f, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

		var node ast.Node
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				node = fn
				break
			}
		}

		assert.True(t, IsIgnored(node, f, fset))

		src2 := "package p\nfunc main() {}"
		f2, _ := parser.ParseFile(fset, "test2.go", src2, parser.ParseComments)
		node2 := f2.Decls[0]
		assert.False(t, IsIgnored(node2, f2, fset))
	})
}

func TestSecurityViolationString(t *testing.T) {
	v := SecurityViolation{File: "test.go", Line: 10, Type: "SQLi", Message: "unsafe query"}
	s := v.String()
	assert.Contains(t, s, "test.go:10")
	assert.Contains(t, s, "SQLi")
	assert.Contains(t, s, "unsafe query")
}

func TestLargeFiles(t *testing.T) {
	// Simple test for ignoring large or invalid strings
	assert.False(t, IsIgnoredRaw(strings.Repeat("a", 2000)))
}
