package agent

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/agent/security/sql"
	"github.com/jadecobra/agbalumo/internal/agent/security/web"
)

func checkFile(path string) ([]SecurityViolation, error) {
	var violations []SecurityViolation

	rawViolations, err := checkSecretsRaw(path)
	if err == nil {
		violations = append(violations, rawViolations...)
	}

	structuralViolations, err := checkStructuralRaw(path)
	if err == nil {
		violations = append(violations, structuralViolations...)
	}

	if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
		fset := token.NewFileSet()
		f, err := internalOpen(path)
		if err != nil {
			return violations, err
		}
		defer func() { _ = f.Close() }()

		node, err := parser.ParseFile(fset, path, f, parser.ParseComments)
		if err != nil {
			return violations, err
		}

		violations = append(violations, sql.CheckSQLi(node, fset)...)
		violations = append(violations, web.CheckXSS(node, fset)...)
		violations = append(violations, web.CheckSSRF(node, fset)...)
		violations = append(violations, checkFileInclusion(node, fset)...)
		violations = append(violations, checkInsecurePatternsGo(node, fset)...)
		violations = append(violations, checkEntropyGo(node, fset)...)
	}

	return deduplicateViolations(violations), nil
}

func checkSecretsRaw(path string) ([]SecurityViolation, error) {
	file, err := internalOpen(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var violations []SecurityViolation
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if isIgnoredRaw(line) {
			lineNum++
			continue
		}

		for name, re := range secretPatterns {
			if re.MatchString(line) {
				violations = append(violations, SecurityViolation{
					File:    path,
					Line:    lineNum,
					Column:  1,
					Type:    "Secret",
					Message: fmt.Sprintf("Potential secret found (%s)", name),
				})
			}
		}
		lineNum++
	}

	return violations, scanner.Err()
}

func checkEntropyGo(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		val := strings.Trim(lit.Value, "\"`")
		trimmed := strings.TrimSpace(val)
		upper := strings.ToUpper(trimmed)
		isSQL := strings.HasPrefix(upper, "SELECT") ||
			strings.HasPrefix(upper, "INSERT") ||
			strings.HasPrefix(upper, "UPDATE") ||
			strings.HasPrefix(upper, "DELETE") ||
			strings.HasPrefix(upper, "CREATE") ||
			strings.HasPrefix(upper, "DROP") ||
			strings.HasPrefix(upper, "ALTER") ||
			strings.HasPrefix(upper, "WITH") ||
			strings.HasPrefix(upper, "PRAGMA") ||
			strings.Contains(upper, "TABLE") ||
			strings.Contains(upper, "JOIN") ||
			strings.Contains(upper, "GROUP BY")
		isURL := strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
		isSentence := strings.Count(trimmed, " ") > 5
		if len(val) >= 40 && calculateEntropy(val) > 5.0 && !isSQL && !isURL && !isSentence {
			if !isIgnored(n, node, fset) {
				pos := fset.Position(lit.Pos())
				violations = append(violations, SecurityViolation{
					File:    pos.Filename,
					Line:    pos.Line,
					Column:  pos.Column,
					Type:    "Entropy",
					Message: "High entropy string detected (potential secret)",
				})
			}
		}
		return true
	})

	return violations
}

func checkInsecurePatternsGo(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		val := lit.Value
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '`' && val[len(val)-1] == '`') {
				val = val[1 : len(val)-1]
			}
		}

		for name, re := range insecureGoPatterns {
			if re.MatchString(val) && !isIgnored(n, node, fset) {
				pos := fset.Position(lit.Pos())
				violations = append(violations, SecurityViolation{
					File:    pos.Filename,
					Line:    pos.Line,
					Column:  pos.Column,
					Type:    "Structural",
					Message: fmt.Sprintf("Insecure pattern found in Go string: %s", name),
				})
			}
		}
		return true
	})

	return violations
}

func checkStructuralRaw(path string) ([]SecurityViolation, error) {
	file, err := internalOpen(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var violations []SecurityViolation
	scanner := bufio.NewScanner(file)
	lineNum := 1
	ext := strings.ToLower(filepath.Ext(path))

	for scanner.Scan() {
		line := scanner.Text()
		if isIgnoredRaw(line) {
			lineNum++
			continue
		}

		for name, re := range structuralPatterns {
			if ext == ".go" {
				if _, exists := insecureGoPatterns[name]; exists {
					continue
				}
			}

			if re.MatchString(line) {
				violations = append(violations, SecurityViolation{
					File:    path,
					Line:    lineNum,
					Column:  1,
					Type:    "Structural",
					Message: fmt.Sprintf("Insecure pattern found: %s", name),
				})
			}
		}

		if ext == ".html" {
			for name, re := range web.HTMLPatterns {
				if re.MatchString(line) {
					if name == "Inline Script" {
						if strings.Contains(strings.ToLower(line), "src=") {
							continue
						}
					}
					violations = append(violations, SecurityViolation{
						File:    path,
						Line:    lineNum,
						Column:  1,
						Type:    "Structural",
						Message: fmt.Sprintf("HTML security issue: %s", name),
					})
				}
			}
		}

		lineNum++
	}

	return violations, scanner.Err()
}
