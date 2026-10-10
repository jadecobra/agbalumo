package maintenance

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ErrorSwallowViolation represents a location where an error return value was swallowed by assigning to blank identifier.
type ErrorSwallowViolation struct {
	File    string
	Message string
	Line    int
}

var errInterface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func implementsError(t types.Type) bool {
	if t == nil {
		return false
	}
	return types.Implements(t, errInterface)
}

func unparen(expr ast.Expr) ast.Expr {
	for {
		p, ok := expr.(*ast.ParenExpr)
		if !ok {
			return expr
		}
		expr = p.X
	}
}

func isExemptMethodCall(call *ast.CallExpr) bool {
	if call == nil {
		return false
	}
	fun := unparen(call.Fun)
	var name string
	switch f := fun.(type) {
	case *ast.SelectorExpr:
		name = f.Sel.Name
	case *ast.Ident:
		name = f.Name
	default:
		return false
	}
	return name == "Close" || name == "Rollback" || name == "CloseWithError"
}

func isCommentExempt(lineComments map[int][]string, startLine, endLine int) bool {
	for l := startLine - 1; l <= endLine; l++ {
		for _, text := range lineComments[l] {
			if strings.Contains(text, "nolint:error-swallow") || strings.Contains(text, "explicit-ignore") {
				return true
			}
		}
	}
	return false
}

type swallowVisitor struct {
	pkg          *packages.Package
	violations   *[]ErrorSwallowViolation
	lineComments map[int][]string
	inDefer      bool
}

func (v *swallowVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	if _, ok := node.(*ast.DeferStmt); ok {
		child := *v
		child.inDefer = true
		return &child
	}

	if assign, ok := node.(*ast.AssignStmt); ok {
		if !v.inDefer {
			v.checkAssign(assign)
		}
		return v
	}

	return v
}

func (v *swallowVisitor) checkAssign(assign *ast.AssignStmt) {
	startPos := v.pkg.Fset.Position(assign.Pos())
	endPos := v.pkg.Fset.Position(assign.End())
	if isCommentExempt(v.lineComments, startPos.Line, endPos.Line) {
		return
	}

	if len(assign.Rhs) == 1 && len(assign.Lhs) > 1 {
		v.checkMultiReturn(assign)
		return
	}

	if len(assign.Lhs) == len(assign.Rhs) {
		v.checkSingleReturn(assign)
	}
}

func (v *swallowVisitor) checkMultiReturn(assign *ast.AssignStmt) {
	rhsExpr := assign.Rhs[0]
	if call, isCall := unparen(rhsExpr).(*ast.CallExpr); isCall && isExemptMethodCall(call) {
		return
	}

	tv := v.resolveType(rhsExpr)
	if tv.Type == nil {
		return
	}

	tuple, ok := tv.Type.(*types.Tuple)
	if !ok {
		return
	}

	v.checkTupleAssign(tuple, assign.Lhs)
}

func (v *swallowVisitor) resolveType(expr ast.Expr) types.TypeAndValue {
	if tv, ok := v.pkg.TypesInfo.Types[expr]; ok {
		return tv
	}
	return v.pkg.TypesInfo.Types[unparen(expr)]
}

func (v *swallowVisitor) checkTupleAssign(tuple *types.Tuple, lhs []ast.Expr) {
	for i := 0; i < len(lhs) && i < tuple.Len(); i++ {
		ident, ok := lhs[i].(*ast.Ident)
		if !ok || ident.Name != "_" {
			continue
		}

		if implementsError(tuple.At(i).Type()) {
			pos := v.pkg.Fset.Position(lhs[i].Pos())
			*v.violations = append(*v.violations, ErrorSwallowViolation{
				File:    pos.Filename,
				Message: "error return value discarded with blank identifier in multi-return assignment",
				Line:    pos.Line,
			})
		}
	}
}

func (v *swallowVisitor) checkSingleReturn(assign *ast.AssignStmt) {
	for i := 0; i < len(assign.Lhs); i++ {
		v.checkSingleAssignElement(assign.Lhs[i], assign.Rhs[i])
	}
}

func (v *swallowVisitor) checkSingleAssignElement(lhs, rhs ast.Expr) {
	ident, ok := lhs.(*ast.Ident)
	if !ok || ident.Name != "_" {
		return
	}

	call, isCall := unparen(rhs).(*ast.CallExpr)
	if !isCall || isExemptMethodCall(call) {
		return
	}

	tv := v.resolveType(rhs)
	if tv.Type == nil || !implementsError(tv.Type) {
		return
	}

	pos := v.pkg.Fset.Position(lhs.Pos())
	*v.violations = append(*v.violations, ErrorSwallowViolation{
		File:    pos.Filename,
		Message: "error return value discarded with blank identifier in call assignment",
		Line:    pos.Line,
	})
}

// CheckErrorSwallow inspects Go source files in targetDirs under rootDir for discarded error returns.
func CheckErrorSwallow(rootDir string, targetDirs ...string) ([]ErrorSwallowViolation, error) {
	patterns := resolveTargetPatterns(rootDir, targetDirs)
	if len(patterns) == 0 {
		return nil, nil
	}

	cfg := &packages.Config{
		Dir:  rootDir,
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	return scanPackages(pkgs), nil
}

func resolveTargetPatterns(rootDir string, targetDirs []string) []string {
	if len(targetDirs) == 0 {
		targetDirs = []string{"internal/repository", "internal/service"}
	}

	var patterns []string
	for _, td := range targetDirs {
		p := normalizePattern(rootDir, td)
		if p != "" {
			patterns = append(patterns, p)
		}
	}
	return patterns
}

func normalizePattern(rootDir, td string) string {
	tdClean := filepath.Clean(td)
	fullPath := filepath.Join(rootDir, tdClean)
	info, err := os.Stat(fullPath)
	if err != nil || !info.IsDir() {
		return ""
	}

	if tdClean == "." {
		return "./..."
	}
	p := filepath.ToSlash(tdClean)
	if !strings.HasPrefix(p, "./") {
		p = "./" + p
	}
	if !strings.HasSuffix(p, "/...") {
		p = strings.TrimSuffix(p, "/") + "/..."
	}
	return p
}

func scanPackages(pkgs []*packages.Package) []ErrorSwallowViolation {
	var violations []ErrorSwallowViolation
	for _, pkg := range pkgs {
		if len(pkg.Syntax) == 0 || pkg.TypesInfo == nil {
			continue
		}
		violations = append(violations, scanFiles(pkg)...)
	}
	return violations
}

func scanFiles(pkg *packages.Package) []ErrorSwallowViolation {
	var violations []ErrorSwallowViolation
	for _, file := range pkg.Syntax {
		if !file.Pos().IsValid() {
			continue
		}
		filename := pkg.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}
		violations = append(violations, scanFile(pkg, file)...)
	}
	return violations
}

func scanFile(pkg *packages.Package, file *ast.File) []ErrorSwallowViolation {
	var violations []ErrorSwallowViolation
	visitor := &swallowVisitor{
		pkg:          pkg,
		violations:   &violations,
		lineComments: extractLineComments(pkg, file),
	}
	ast.Walk(visitor, file)
	return violations
}

func extractLineComments(pkg *packages.Package, file *ast.File) map[int][]string {
	lineComments := make(map[int][]string)
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			start := pkg.Fset.Position(c.Pos()).Line
			end := pkg.Fset.Position(c.End()).Line
			for l := start; l <= end; l++ {
				lineComments[l] = append(lineComments[l], c.Text)
			}
		}
	}
	return lineComments
}
