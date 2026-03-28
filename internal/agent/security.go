package agent

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// SecurityViolation represents a potential security issue found in the code.
type SecurityViolation struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
}

func (v SecurityViolation) String() string {
	return fmt.Sprintf("%s:%d:%d: [%s] %s", v.File, v.Line, v.Column, v.Type, v.Message)
}

// VerifySecurityStatic runs static analysis checkers for security vulnerabilities.
func VerifySecurityStatic(root string) ([]SecurityViolation, error) {
	var allViolations []SecurityViolation

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, non-Go files, and hidden files/vendor
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.Contains(path, "/.") || strings.Contains(path, "/vendor/") {
			return nil
		}

		// Skip tests
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		violations, err := checkFile(path)
		if err != nil {
			return fmt.Errorf("failed to check file %s: %w", path, err)
		}

		allViolations = append(allViolations, violations...)
		return nil
	})

	return allViolations, err
}

func checkFile(path string) ([]SecurityViolation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var violations []SecurityViolation
	violations = append(violations, checkSQLi(node, fset)...)
	violations = append(violations, checkXSS(node, fset)...)

	return violations, nil
}

// checkSQLi inspects a file for potential SQL injection vulnerabilities.
func checkSQLi(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

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
				if isUnsafeString(arg) && !isIgnored(n, node, fset) {
					pos := fset.Position(arg.Pos())
					violations = append(violations, SecurityViolation{
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

// checkXSS inspects a file for potential cross-site scripting vulnerabilities.
func checkXSS(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Case 1: template.HTML(variable)
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "template" && sel.Sel.Name == "HTML" {
				if len(call.Args) > 0 {
					arg := call.Args[0]
					if _, ok := arg.(*ast.BasicLit); !ok && !isIgnored(n, node, fset) {
						pos := fset.Position(arg.Pos())
						violations = append(violations, SecurityViolation{
							File:    pos.Filename,
							Line:    pos.Line,
							Column:  pos.Column,
							Type:    "XSS",
							Message: "Potential XSS: template.HTML used with non-literal value",
						})
					}
				}
			}

			// Case 2: c.HTML(code, html)
			if sel.Sel.Name == "HTML" {
				// Check if it's likely an Echo context (conventionally 'c')
				if x, ok := sel.X.(*ast.Ident); ok && (x.Name == "c" || x.Name == "ctx") {
					if len(call.Args) >= 2 {
						arg := call.Args[1]
						if isUnsafeString(arg) && !isIgnored(n, node, fset) {
							pos := fset.Position(arg.Pos())
							violations = append(violations, SecurityViolation{
								File:    pos.Filename,
								Line:    pos.Line,
								Column:  pos.Column,
								Type:    "XSS",
								Message: "Potential XSS: unsafe HTML construction in c.HTML",
							})
						}
					}
				}
			}
		}

		return true
	})

	return violations
}

// isUnsafeString checks if an expression is a concatenation or a Sprintf call.
func isUnsafeString(expr ast.Expr) bool {
	switch e := expr.(type) {
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
	}
	return false
}

func isIgnored(n ast.Node, file *ast.File, fset *token.FileSet) bool {
	pos := fset.Position(n.Pos())
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			cpos := fset.Position(c.Pos())
			// Ignore if comment is on the same line or the line directly above
			if cpos.Line == pos.Line || cpos.Line == pos.Line-1 {
				if strings.Contains(c.Text, "#nosec") || strings.Contains(c.Text, "antigravity:allow") {
					return true
				}
			}
		}
	}
	return false
}
