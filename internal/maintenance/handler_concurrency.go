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

// HandlerConcurrencyViolation represents a violation of handler concurrency bounds.
type HandlerConcurrencyViolation struct {
	File    string
	Func    string
	Message string
	Line    int
}

// HandlerConcurrencyOptions configures the handler concurrency guard.
type HandlerConcurrencyOptions struct {
	MaxGoroutinesPerFunc int
}

// DefaultHandlerConcurrencyOptions provides recommended safety defaults.
func DefaultHandlerConcurrencyOptions() HandlerConcurrencyOptions {
	return HandlerConcurrencyOptions{
		MaxGoroutinesPerFunc: 4,
	}
}

// CheckHandlerConcurrency scans all HTTP handler files in internal/module/
// for unbounded goroutine loops or goroutine fan-outs exceeding the concurrency budget.
func CheckHandlerConcurrency(rootDir string, opts HandlerConcurrencyOptions) ([]HandlerConcurrencyViolation, error) {
	if opts.MaxGoroutinesPerFunc <= 0 {
		opts.MaxGoroutinesPerFunc = 4
	}

	moduleDir := filepath.Join(rootDir, "internal", "module")
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return nil, nil
	}

	var violations []HandlerConcurrencyViolation
	walkErr := filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipConcurrencyScan(path, info) {
			return nil
		}

		fileViolations, scanErr := inspectHandlerConcurrencyFile(path, opts)
		if scanErr != nil {
			return scanErr
		}
		violations = append(violations, fileViolations...)
		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("failed to scan module directory: %w", walkErr)
	}

	return violations, nil
}

func shouldSkipConcurrencyScan(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}
	return !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")
}

func inspectHandlerConcurrencyFile(filePath string, opts HandlerConcurrencyOptions) ([]HandlerConcurrencyViolation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	var violations []HandlerConcurrencyViolation

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		funcName := resolveFullFuncName(funcDecl)
		v := checkFunctionBodyConcurrency(fset, filePath, funcName, funcDecl.Body, opts)
		violations = append(violations, v...)
		return false
	})

	return violations, nil
}

func resolveFullFuncName(funcDecl *ast.FuncDecl) string {
	funcName := funcDecl.Name.Name
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		recvType := extractReceiverTypeName(funcDecl.Recv.List[0].Type)
		if recvType != "" {
			return recvType + "." + funcName
		}
	}
	return funcName
}

func extractReceiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return extractReceiverTypeName(t.X)
	case *ast.Ident:
		return t.Name
	default:
		return ""
	}
}

func checkFunctionBodyConcurrency(fset *token.FileSet, filePath, funcName string, body *ast.BlockStmt, opts HandlerConcurrencyOptions) []HandlerConcurrencyViolation {
	var violations []HandlerConcurrencyViolation
	var goStmtCount int

	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch stmt := n.(type) {
		case *ast.ForStmt:
			violations = append(violations, inspectLoopForGoStmt(fset, filePath, funcName, stmt.Body)...)
		case *ast.RangeStmt:
			violations = append(violations, inspectLoopForGoStmt(fset, filePath, funcName, stmt.Body)...)
		case *ast.GoStmt:
			goStmtCount++
		}

		return true
	})

	if goStmtCount > opts.MaxGoroutinesPerFunc {
		pos := fset.Position(body.Pos())
		violations = append(violations, HandlerConcurrencyViolation{
			File:    filePath,
			Line:    pos.Line,
			Func:    funcName,
			Message: fmt.Sprintf("goroutine fan-out of %d exceeds concurrency limit of %d", goStmtCount, opts.MaxGoroutinesPerFunc),
		})
	}

	return violations
}

func inspectLoopForGoStmt(fset *token.FileSet, filePath, funcName string, body *ast.BlockStmt) []HandlerConcurrencyViolation {
	var violations []HandlerConcurrencyViolation

	ast.Inspect(body, func(n ast.Node) bool {
		if goStmt, ok := n.(*ast.GoStmt); ok {
			pos := fset.Position(goStmt.Pos())
			violations = append(violations, HandlerConcurrencyViolation{
				File:    filePath,
				Line:    pos.Line,
				Func:    funcName,
				Message: "unbounded goroutine spawn in loop inside HTTP handler",
			})
		}
		return true
	})

	return violations
}
