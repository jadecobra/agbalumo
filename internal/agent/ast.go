package agent

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

// Route represents an extracted HTTP route.
type Route struct {
	Method string
	Path   string
}

// ExtractRoutes parses a Go source file and extracts Echo routes.
func ExtractRoutes(filename string) ([]Route, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Map of variable names to their group prefix paths.
	// E.g. "adminGroup" -> "/admin"
	// "e" -> ""
	groupPaths := make(map[string]string)
	
	// Add default echo instance if we find it (often 'e' or 'app')
	// In our cmd/server.go, it's typically 'e'
	groupPaths["e"] = ""

	var routes []Route
	var muList = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	isHttpMethod := func(method string) bool {
		for _, m := range muList {
			if m == method {
				return true
			}
		}
		return false
	}

	normalizePath := func(p string) string {
		// Replace :id with {id}
		p = regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(p, "{$1}")
		// Remove trailing slashes (except root)
		if len(p) > 1 && strings.HasSuffix(p, "/") {
			p = strings.TrimSuffix(p, "/")
		}
		// Deduplicate slashes
		p = regexp.MustCompile(`//+`).ReplaceAllString(p, "/")
		if p == "" {
			p = "/"
		}
		return p
	}

	// First pass: Find group definitions
	ast.Inspect(node, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Look for `adminGroup := e.Group("/admin")` or `adminLoginGroup := adminGroup.Group("/login")`
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			return true
		}

		ident, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}

		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if selExpr.Sel.Name == "Group" {
			receiver, ok := selExpr.X.(*ast.Ident)
			if ok {
				if len(callExpr.Args) > 0 {
					if argLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && argLit.Kind == token.STRING {
						prefix := strings.Trim(argLit.Value, "\"")
						base := groupPaths[receiver.Name]
						groupPaths[ident.Name] = base + prefix
					}
				}
			}
		}

		return true
	})

	// Second pass: Find route definitions
	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		httpMethod := selExpr.Sel.Name
		if !isHttpMethod(httpMethod) {
			return true
		}

		receiver, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return true
		}

		basePath, exists := groupPaths[receiver.Name]
		if !exists {
			// Not a recognized router group/instance.
			return true
		}

		if len(callExpr.Args) > 0 {
			if argLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && argLit.Kind == token.STRING {
				pathSuffix := strings.Trim(argLit.Value, "\"")
				fullPath := normalizePath(basePath + pathSuffix)
				routes = append(routes, Route{
					Method: httpMethod,
					Path:   fullPath,
				})
			}
		}

		return true
	})

	// Sort routes for deterministic output
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	// Deduplicate
	var uniqueRoutes []Route
	seen := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		if !seen[key] {
			seen[key] = true
			uniqueRoutes = append(uniqueRoutes, r)
		}
	}

	return uniqueRoutes, nil
}
