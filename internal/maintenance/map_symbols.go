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
			if extractedType := handleTypeSpec(x, rel); extractedType != "" {
				extracted = append(extracted, extractedType)
			}
		case *ast.FuncDecl:
			if extractedFunc := handleFuncDecl(x, rel); extractedFunc != "" {
				extracted = append(extracted, extractedFunc)
			}
		}
		return true
	})
	return extracted
}

func handleTypeSpec(x *ast.TypeSpec, rel string) string {
	if !x.Name.IsExported() {
		return ""
	}
	label := "[Type]"
	if _, ok := x.Type.(*ast.InterfaceType); ok {
		label = "[Interface]"
	}
	return fmt.Sprintf("%s %s -> %s", label, x.Name.Name, rel)
}

func handleFuncDecl(x *ast.FuncDecl, rel string) string {
	if !x.Name.IsExported() {
		return ""
	}
	if x.Recv != nil && len(x.Recv.List) > 0 {
		recvType := extractReceiverType(x.Recv.List[0].Type)
		if recvType != "" {
			return fmt.Sprintf("[Method] (%s).%s -> %s", recvType, x.Name.Name, rel)
		}
	}
	return fmt.Sprintf("[Func] %s -> %s", x.Name.Name, rel)
}

func extractReceiverType(t ast.Expr) string {
	if star, ok := t.(*ast.StarExpr); ok {
		if ident, ok := star.X.(*ast.Ident); ok {
			return ident.Name
		}
	} else if ident, ok := t.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}
