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

// PipeDrainViolation represents a synchronous os.Pipe read without concurrent drain.
type PipeDrainViolation struct {
	File    string
	Message string
	Line    int
}

// CheckPipeDrain inspects all *_test.go files in rootDir for os.Pipe calls,
// verifying that reading from the pipe reader occurs in a separate goroutine
// before closing or waiting on execution.
func CheckPipeDrain(rootDir string) ([]PipeDrainViolation, error) {
	var violations []PipeDrainViolation

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return handleDirWalk(path, rootDir, info)
		}
		if !isTargetTestFile(path, info.Name()) {
			return nil
		}

		fileViolations, scanErr := checkFilePipeDrain(path)
		if scanErr != nil {
			return scanErr
		}
		violations = append(violations, fileViolations...)
		return nil
	}

	if err := filepath.Walk(rootDir, walkFn); err != nil {
		return nil, fmt.Errorf("failed scanning for pipe drain violations: %w", err)
	}

	return violations, nil
}

func isTargetTestFile(path, name string) bool {
	if filepath.Ext(path) != extGo || !strings.HasSuffix(name, "_test.go") {
		return false
	}
	return name != "pipe_drain_test.go"
}

func checkFilePipeDrain(filePath string) ([]PipeDrainViolation, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	var violations []PipeDrainViolation

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		violations = append(violations, checkFuncBlock(fset, filePath, funcDecl.Body)...)
		return false
	})

	return violations, nil
}

func checkFuncBlock(fset *token.FileSet, filePath string, body *ast.BlockStmt) []PipeDrainViolation {
	var violations []PipeDrainViolation

	for _, stmt := range body.List {
		pipeCalls := findPipeAssignments(stmt)
		for _, pipe := range pipeCalls {
			if !isReaderDrainedConcurrently(body, pipe.readerName) {
				pos := fset.Position(pipe.callPos)
				violations = append(violations, PipeDrainViolation{
					File:    filePath,
					Line:    pos.Line,
					Message: fmt.Sprintf("pipe reader %q is not drained in a concurrent goroutine; synchronous reads block when output exceeds OS pipe buffer limit (use testutil.CaptureStdout or a concurrent goroutine drain)", pipe.readerName),
				})
			}
		}
	}

	return violations
}

type pipeAssignment struct {
	readerName string
	callPos    token.Pos
}

func findPipeAssignments(stmt ast.Stmt) []pipeAssignment {
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	var results []pipeAssignment
	for i, rhs := range assignStmt.Rhs {
		call, ok := rhs.(*ast.CallExpr)
		if !ok || !isOSPipeCall(call) {
			continue
		}

		readerName := extractAssigneeName(assignStmt.Lhs, i)
		results = append(results, pipeAssignment{
			readerName: readerName,
			callPos:    call.Pos(),
		})
	}

	return results
}

func extractAssigneeName(lhs []ast.Expr, index int) string {
	if index >= len(lhs) {
		return "_"
	}
	if id, ok := lhs[index].(*ast.Ident); ok {
		return id.Name
	}
	return "_"
}

func isOSPipeCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkgIdent.Name == "os" && sel.Sel.Name == "Pipe"
}

func isReaderDrainedConcurrently(body *ast.BlockStmt, readerName string) bool {
	if readerName == "" || readerName == "_" {
		return false
	}

	foundConcurrentDrain := false

	ast.Inspect(body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}

		ast.Inspect(goStmt, func(inner ast.Node) bool {
			id, ok := inner.(*ast.Ident)
			if ok && id.Name == readerName {
				foundConcurrentDrain = true
				return false
			}
			return true
		})

		return true
	})

	return foundConcurrentDrain
}
