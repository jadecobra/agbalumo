package maintenance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExtractRoutes parses Go source files and extracts Echo routes using AST analysis.
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
			_ = filepath.Walk(p, func(path string, i os.FileInfo, e error) error {
				if e == nil && !i.IsDir() && strings.HasSuffix(path, ".go") {
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

	groupPaths := make(map[string]string)
	groupPaths["e"] = "" // Standard Echo instance naming

	var routes []Route
	muList := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	isHttpMethod := func(method string) bool {
		for _, m := range muList {
			if m == method {
				return true
			}
		}
		return false
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
				basePath = ""
			}
			if len(callExpr.Args) > 0 {
				if argLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && argLit.Kind == token.STRING {
					pathSuffix := strings.Trim(argLit.Value, "\"")
					if pathSuffix == "" || strings.HasPrefix(pathSuffix, "/") {
						fullPath := NormalizePath(basePath + pathSuffix)
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

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

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
