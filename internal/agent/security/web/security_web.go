package web

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

	"github.com/jadecobra/agbalumo/internal/agent/security/sectypes"
)

type WebHandler struct{}

var (
	// File-type specific patterns for web/HTML
	HTMLPatterns = map[string]*regexp.Regexp{
		"Inline Script": regexp.MustCompile(`(?i)<script`),
	}
)

// CheckXSS inspects a file for potential cross-site scripting vulnerabilities.
func CheckXSS(node *ast.File, fset *token.FileSet) []sectypes.SecurityViolation {
	var violations []sectypes.SecurityViolation

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
					if _, ok := arg.(*ast.BasicLit); !ok && !sectypes.IsIgnored(n, node, fset) {
						pos := fset.Position(arg.Pos())
						violations = append(violations, sectypes.SecurityViolation{
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
						if sectypes.IsUnsafeString(arg) && !sectypes.IsIgnored(n, node, fset) {
							pos := fset.Position(arg.Pos())
							violations = append(violations, sectypes.SecurityViolation{
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

// CheckSSRF inspects for potential SSRF vulnerabilities.
func CheckSSRF(node *ast.File, fset *token.FileSet) []sectypes.SecurityViolation {
	var violations []sectypes.SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var methodName string
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			methodName = sel.Sel.Name
			// Check for http.Get, http.Post, etc. or client.Do(req)
			if x, ok := sel.X.(*ast.Ident); ok {
				if x.Name == "http" && (methodName == "Get" || methodName == "Post" || methodName == "PostForm" || methodName == "Head") {
					if len(call.Args) > 0 && sectypes.IsUnsafeString(call.Args[0]) && !sectypes.IsIgnored(n, node, fset) {
						pos := fset.Position(call.Args[0].Pos())
						violations = append(violations, sectypes.SecurityViolation{
							File:    pos.Filename,
							Line:    pos.Line,
							Column:  pos.Column,
							Type:    "SSRF",
							Message: fmt.Sprintf("Potential SSRF: unsafe URL construction in http.%s", methodName),
						})
					}
				}
				if methodName == "Do" {
					// Potential Client.Do(req). If req was constructed unsafely.
					// This is a bit complex for a simple AST check, but we can look for 'req' as argument
					if len(call.Args) > 0 && !sectypes.IsIgnored(n, node, fset) {
						// For now, let's flag all .Do(req) where the URL might be tainted.
						pos := fset.Position(call.Pos())
						violations = append(violations, sectypes.SecurityViolation{
							File:    pos.Filename,
							Line:    pos.Line,
							Column:  pos.Column,
							Type:    "SSRF",
							Message: "Potential SSRF via taint analysis on http.DefaultClient.Do(req)",
						})
					}
				}
			}
		}

		return true
	})

	return violations
}
