package agent

import (
	"fmt"
	"go/ast"
	"go/token"
)

// checkFileInclusion inspects for potential file inclusion vulnerabilities.
func checkFileInclusion(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var pkgName, methodName string
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			methodName = sel.Sel.Name
			if x, ok := sel.X.(*ast.Ident); ok {
				pkgName = x.Name
			}
		}

		// Check for os.Open, os.ReadFile, io/ioutil.ReadFile, etc.
		if (pkgName == "os" || pkgName == "ioutil") && (methodName == "Open" || methodName == "ReadFile") {
			if len(call.Args) > 0 {
				arg := call.Args[0]
				// If it's not a basic string literal, flag it
				if _, ok := arg.(*ast.BasicLit); !ok && !isIgnored(n, node, fset) {
					pos := fset.Position(arg.Pos())
					violations = append(violations, SecurityViolation{
						File:    pos.Filename,
						Line:    pos.Line,
						Column:  pos.Column,
						Type:    "FileInclusion",
						Message: fmt.Sprintf("Potential file inclusion via variable on %s.%s", pkgName, methodName),
					})
				}
			}
		}

		return true
	})

	return violations
}
