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

// GenerateSymbolMap scans all .go files in rootDir and returns a formatted symbol map.
func GenerateSymbolMap(rootDir string) string {
	var symbols []string
	fset := token.NewFileSet()
	absRoot, _ := filepath.Abs(rootDir)

	_ = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, _ := filepath.Rel(absRoot, path)
		extracted := extractSymbolsFromFile(fset, path, rel)
		symbols = append(symbols, extracted...)
		return nil
	})

	sort.Strings(symbols)
	return strings.Join(symbols, "\n")
}

func extractSymbolsFromFile(fset *token.FileSet, path, rel string) []string {
	var extracted []string
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.IsExported() {
				extracted = append(extracted, fmt.Sprintf("[Type] %s -> %s", x.Name.Name, rel))
			}
		case *ast.FuncDecl:
			if x.Name.IsExported() {
				extracted = append(extracted, fmt.Sprintf("[Func] %s -> %s", x.Name.Name, rel))
			}
		}
		return true
	})
	return extracted
}

// GenerateTemplateMap calls the existing BuildTemplateGraph and returns a formatted template map.
func GenerateTemplateMap(rootDir string) string {
	graph, err := BuildTemplateGraph(rootDir)
	if err != nil {
		return fmt.Sprintf("Error building template graph: %v", err)
	}

	var templates []string
	for name, node := range graph {
		rel := node.DefinedIn
		if absRoot, err := filepath.Abs(rootDir); err == nil {
			if r, err := filepath.Rel(absRoot, node.DefinedIn); err == nil {
				rel = r
			}
		}
		templates = append(templates, fmt.Sprintf("[Template] %s -> %s", name, rel))
	}

	sort.Strings(templates)
	return strings.Join(templates, "\n")
}
