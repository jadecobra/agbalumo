package sectypes

import (
	"go/ast"
	"go/token"
	"math"
	"strings"
)

func CalculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}

	var entropy float64
	for _, count := range counts {
		p := float64(count) / float64(len(s))
		entropy -= p * math.Log2(p)
	}

	return entropy
}

func IsIgnored(n ast.Node, file *ast.File, fset *token.FileSet) bool {
	pos := fset.Position(n.Pos())
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			cpos := fset.Position(c.Pos())
			if cpos.Line == pos.Line || cpos.Line == pos.Line-1 {
				if strings.Contains(c.Text, "#nosec") || strings.Contains(c.Text, "antigravity:allow") {
					return true
				}
			}
		}
	}
	return false
}

func IsIgnoredRaw(line string) bool {
	return strings.Contains(line, "#nosec") || strings.Contains(line, "antigravity:allow")
}

func IsUnsafeString(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return false // Literals are safe
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return true
		}
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "fmt" && sel.Sel.Name == "Sprintf" {
				return true
			}
		}
	case *ast.Ident, *ast.SelectorExpr:
		return true // Dynamic values are unsafe
	}
	return false
}
