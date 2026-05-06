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

type ContractViolation struct {
	Template   string
	Field      string
	ExpectedIn string
}

func CheckTemplateContracts(rootDir string) ([]ContractViolation, error) {
	templatesDir := filepath.Join(rootDir, "ui", "templates")
	moduleDir := filepath.Join(rootDir, "internal", "module")

	graph, err := BuildTemplateGraph(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to build template graph: %w", err)
	}

	allowedFields, err := extractViewModelFields(moduleDir)
	if err != nil {
		return nil, fmt.Errorf("failed to extract view model fields: %w", err)
	}

	// Also allow common template built-ins or special cases if needed
	// For now, stick to the prompt's strict requirements.

	var violations []ContractViolation
	for name, node := range graph {
		for _, ref := range node.References {
			if !allowedFields[ref] {
				violations = append(violations, ContractViolation{
					Template:   name,
					Field:      ref,
					ExpectedIn: "ViewModels/ViewData",
				})
			}
		}
	}

	return violations, nil
}

func extractViewModelFields(moduleDir string) (map[string]bool, error) {
	fields := make(map[string]bool)
	fset := token.NewFileSet()

	err := filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		extractFieldsFromNode(node, fields)
		return nil
	})

	return fields, err
}

func extractFieldsFromNode(node *ast.File, fields map[string]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if !isViewType(ts.Name.Name) {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		collectStructFields(st, fields)
		return true
	})
}

func isViewType(name string) bool {
	return strings.HasSuffix(name, "ViewModel") || strings.HasSuffix(name, "ViewData")
}

func collectStructFields(st *ast.StructType, fields map[string]bool) {
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			if ident, ok := field.Type.(*ast.Ident); ok {
				fields[ident.Name] = true
			}
		} else {
			for _, name := range field.Names {
				fields[name.Name] = true
			}
		}
	}
}
