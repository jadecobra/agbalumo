package maintenance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ExtractRoutes parses Go source files and extracts Echo routes using AST analysis.
func ExtractRoutes(paths ...string) ([]Route, error) {
	fset := token.NewFileSet()
	goFiles := collectGoFiles(paths)
	allFiles := parseGoFiles(fset, goFiles)

	if len(allFiles) == 0 {
		return nil, fmt.Errorf("no valid Go files found in provided paths")
	}

	groupPaths := extractGroupPaths(allFiles)
	routes := extractRouteDefinitions(allFiles, groupPaths)

	return uniqueAndSort(routes), nil
}

func collectGoFiles(paths []string) []string {
	var goFiles []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			goFiles = append(goFiles, walkGoDir(p)...)
		} else {
			goFiles = append(goFiles, p)
		}
	}
	return goFiles
}

func walkGoDir(dirPath string) []string {
	var files []string
	_ = filepath.Walk(dirPath, func(path string, i os.FileInfo, e error) error {
		if e == nil && !i.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func parseGoFiles(fset *token.FileSet, filePaths []string) []*ast.File {
	var allFiles []*ast.File
	for _, filePath := range filePaths {
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err == nil {
			allFiles = append(allFiles, node)
		}
	}
	return allFiles
}

func extractGroupPaths(files []*ast.File) map[string]string {
	groupPaths := make(map[string]string)
	groupPaths["e"] = ""
	for _, node := range files {
		ast.Inspect(node, func(n ast.Node) bool {
			ident, receiver, prefix := parseGroupCall(n)
			if ident != "" {
				groupPaths[ident] = groupPaths[receiver] + prefix
			}
			return true
		})
	}
	return groupPaths
}

func parseGroupCall(n ast.Node) (ident, receiver, prefix string) {
	assignStmt, ok := n.(*ast.AssignStmt)
	if !ok || len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return "", "", ""
	}

	lhs, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return "", "", ""
	}

	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok || len(callExpr.Args) == 0 {
		return "", "", ""
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Group" {
		return "", "", ""
	}

	recv, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return "", "", ""
	}

	argLit, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || argLit.Kind != token.STRING {
		return "", "", ""
	}

	return lhs.Name, recv.Name, strings.Trim(argLit.Value, "\"")
}

func extractRouteDefinitions(files []*ast.File, groupPaths map[string]string) []Route {
	var routes []Route
	for _, node := range files {
		ast.Inspect(node, func(n ast.Node) bool {
			method, path := parseRouteCall(n, groupPaths)
			if method != "" {
				routes = append(routes, NewRoute(method, path))
			}
			return true
		})
	}
	return routes
}

func parseRouteCall(n ast.Node, groupPaths map[string]string) (method, path string) {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok || len(callExpr.Args) == 0 {
		return "", ""
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || !isHttpMethod(selExpr.Sel.Name) {
		return "", ""
	}

	receiver, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return "", ""
	}

	argLit, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || argLit.Kind != token.STRING {
		return "", ""
	}

	suffix := strings.Trim(argLit.Value, "\"")
	return selExpr.Sel.Name, groupPaths[receiver.Name] + suffix
}

func isHttpMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD":
		return true
	}
	return false
}
