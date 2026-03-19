package agent

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Route represents an extracted HTTP route.
type Route struct {
	Method string
	Path   string
}

// ExtractRoutes parses Go source files and extracts Echo routes.
func ExtractRoutes(paths ...string) ([]Route, error) {
	fset := token.NewFileSet()
	var allFiles []*ast.File

	var goFiles []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			_ = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() && strings.HasSuffix(path, ".go") {
					goFiles = append(goFiles, path)
				}
				return nil
			})
		} else {
			goFiles = append(goFiles, p)
		}
	}

	for _, filePath := range goFiles {
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err == nil {
			allFiles = append(allFiles, node)
		}
	}

	if len(allFiles) == 0 {
		return nil, fmt.Errorf("no valid Go files found in provided paths")
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
	for _, node := range allFiles {
		ast.Inspect(node, func(n ast.Node) bool {
			assignStmt, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}

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
	}

	// Second pass: Find route definitions
	for _, node := range allFiles {
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
				// If the receiver isn't a known group, guess it's a top-level route if it starts with "/"
				// We relax constraints so distributed handlers don't strictly get missed if they use `e.GET` directly
				// without tracking `e` assignment.
				basePath = ""
			}

			if len(callExpr.Args) > 0 {
				if argLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && argLit.Kind == token.STRING {
					pathSuffix := strings.Trim(argLit.Value, "\"")
					
					// Only keep if it looks like a valid route definition (must be relative path starting with / or empty string on an existing group)
					if pathSuffix == "" || strings.HasPrefix(pathSuffix, "/") {
						fullPath := normalizePath(basePath + pathSuffix)
						routes = append(routes, Route{
							Method: httpMethod,
							Path:   fullPath,
						})
					}
				}
			}

			return true
		})
	}

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
