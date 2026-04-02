package sql

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/jadecobra/agbalumo/internal/agent/security/sectypes"
)

type SQLHandler struct{}

// CheckSQLi inspects a file for potential SQL injection vulnerabilities.
func CheckSQLi(node *ast.File, fset *token.FileSet) []sectypes.SecurityViolation {
	var violations []sectypes.SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Look for db.Query, db.Exec, etc.
		method := sel.Sel.Name
		if method == "Query" || method == "Exec" || method == "QueryRow" || method == "Prepare" {
			if len(call.Args) > 0 {
				arg := call.Args[0]
				if sectypes.IsUnsafeString(arg) && !sectypes.IsIgnored(n, node, fset) {
					pos := fset.Position(arg.Pos())
					violations = append(violations, sectypes.SecurityViolation{
						File:    pos.Filename,
						Line:    pos.Line,
						Column:  pos.Column,
						Type:    "SQLi",
						Message: fmt.Sprintf("Potential SQL injection: unsafe string construction in %s", method),
					})
				}
			}
		}

		return true
	})

	return violations
}
