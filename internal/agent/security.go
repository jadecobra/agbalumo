package agent

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

var (
	// internalOpen is a hook for testing file operations.
	internalOpen = util.SafeOpen
)

// SecurityViolation represents a potential security issue found in the code.
type SecurityViolation struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
}

func (v SecurityViolation) String() string {
	return fmt.Sprintf("%s:%d:%d: [%s] %s", v.File, v.Line, v.Column, v.Type, v.Message)
}

var (
	// Common secret patterns
	secretPatterns = map[string]*regexp.Regexp{
		"AWS Access Key": regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"AWS Secret Key": regexp.MustCompile(`(?i)aws(.{0,20})?['"][0-9a-zA-Z/+]{40}['"]`),
		"Slack Webhook":  regexp.MustCompile(`https://hooks\.slack\.com/services/T[a-zA-Z0-9_]+/B[a-zA-Z0-9_]+/[a-zA-Z0-9_]+`),
		"Private Key":    regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY-----`),
		"Generic Secret": regexp.MustCompile(`(?i)(password|secret|key|token|access_token|authorization|auth)\s*[:=]\s*["'][^"']{4,}["']`),
	}

	// Structural security patterns (Check for insecure coding practices)
	// Some of these are split into AST-based checks for .go files to improve precision.
	structuralPatterns = map[string]*regexp.Regexp{
		"Insecure Handler":  regexp.MustCompile(`onclick\s*=`),
		"Dangerous JS":      regexp.MustCompile(`(eval\(|Function\(|innerHTML\s*=)`),
		"Forbidden CDN":     regexp.MustCompile(`https?://(unpkg\.com|cdn\.jsdelivr\.net|cdn\.tailwindcss\.com|jsdelivr\.net)`),
		"Hardcoded OAuth":   regexp.MustCompile(`GetAuthCodeURL\(["'][^"']+["']`),
		"Gosec NoRationale": regexp.MustCompile(`//\s*#nosec\s*($|\n|[^a-zA-Z0-9 ])`),
	}

	// patterns for AST-based checks in Go string literals
	insecureGoPatterns = map[string]*regexp.Regexp{
		"Insecure Handler": regexp.MustCompile(`onclick\s*=`),
		"Dangerous JS":     regexp.MustCompile(`(eval\(|Function\(|innerHTML\s*=)`),
		"Forbidden CDN":    regexp.MustCompile(`https?://(unpkg\.com|cdn\.jsdelivr\.net|cdn\.tailwindcss\.com|jsdelivr\.net)`),
	}

	// File-type specific patterns
	htmlPatterns = map[string]*regexp.Regexp{
		"Inline Script": regexp.MustCompile(`(?i)<script`),
	}
)

// VerifySecurityStatic runs static analysis checkers for security vulnerabilities on multiple targets.
func VerifySecurityStatic(targets ...string) ([]SecurityViolation, error) {
	var allViolations []SecurityViolation

	for _, target := range targets {
		info, err := util.SafeStat(target)
		if err != nil {
			// If target doesn't exist, just skip it (might happen with deleted files in staged list)
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("failed to stat target %s: %w", target, err)
		}

		// If it's a file, check it directly and skip walking
		if !info.IsDir() {
			violations, err := checkFile(target)
			if err != nil {
				return nil, fmt.Errorf("failed to check file %s: %w", target, err)
			}
			allViolations = append(allViolations, violations...)
			continue
		}

		// If it's a directory, walk it
		err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories, vendor, node_modules, and hidden dirs (except .)
			if info.IsDir() {
				if strings.HasPrefix(info.Name(), ".") && info.Name() != "." && info.Name() != target {
					return filepath.SkipDir
				}
				return nil
			}

			if strings.Contains(path, "/vendor/") || strings.Contains(path, "/node_modules/") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			// Skip binary files and large assets
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".pdf" || ext == ".exe" || ext == ".bin" {
				return nil
			}

			violations, err := checkFile(path)
			if err != nil {
				// If it's not a Go file and failed to check, just continue
				if !strings.HasSuffix(path, ".go") {
					return nil
				}
				return fmt.Errorf("failed to check file %s: %w", path, err)
			}

			allViolations = append(allViolations, violations...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return deduplicateViolations(allViolations), nil
}

func checkFile(path string) ([]SecurityViolation, error) {
	var violations []SecurityViolation

	// Always run secret and structural scanning on raw text first
	rawViolations, err := checkSecretsRaw(path)
	if err == nil {
		violations = append(violations, rawViolations...)
	}

	structuralViolations, err := checkStructuralRaw(path)
	if err == nil {
		violations = append(violations, structuralViolations...)
	}

	// For Go files, run AST-based checks
	if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
		fset := token.NewFileSet()
		f, err := internalOpen(path)
		if err != nil {
			return violations, err
		}
		defer func() { _ = f.Close() }()

		node, err := parser.ParseFile(fset, path, f, parser.ParseComments)
		if err != nil {
			return violations, err // Return what we found + error
		}

		violations = append(violations, checkSQLi(node, fset)...)
		violations = append(violations, checkXSS(node, fset)...)
		violations = append(violations, checkInsecurePatternsGo(node, fset)...)
		violations = append(violations, checkEntropyGo(node, fset)...)
	}

	// Deduplicate violations on the same line and file
	return deduplicateViolations(violations), nil
}

func deduplicateViolations(violations []SecurityViolation) []SecurityViolation {
	seen := make(map[string]SecurityViolation)
	var lines []string
	for _, v := range violations {
		key := fmt.Sprintf("%s:%d", v.File, v.Line)
		existing, ok := seen[key]
		if !ok {
			seen[key] = v
			lines = append(lines, key)
		} else if existing.Type == "Entropy" && v.Type == "Secret" {
			// Prioritize "Secret" over "Entropy" for the same line
			seen[key] = v
		}
	}
	var unique []SecurityViolation
	for _, key := range lines {
		unique = append(unique, seen[key])
	}
	return unique
}

// checkSecretsRaw scans a file's raw content for regex-based secrets.
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

// checkEntropyGo scans AST string literals for high entropy.
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

// calculateEntropy returns the Shannon entropy of a string.
func calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}

	var entropy float64
	for _, count := range counts {
		p := float64(count) / float64(len(s))
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// checkSQLi inspects a file for potential SQL injection vulnerabilities.
func checkSQLi(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Look for db.Query, db.Exec, etc.
		method := sel.Sel.Name
		if method == "Query" || method == "Exec" || method == "QueryRow" || method == "Prepare" {
			if len(call.Args) > 0 {
				arg := call.Args[0]
				if isUnsafeString(arg) && !isIgnored(n, node, fset) {
					pos := fset.Position(arg.Pos())
					violations = append(violations, SecurityViolation{
						File:    pos.Filename,
						Line:    pos.Line,
						Column:  pos.Column,
						Type:    "SQLi",
						Message: fmt.Sprintf("Potential SQL injection: unsafe string construction in %s", method),
					})
				}
			}
		}

		return true
	})

	return violations
}

// checkXSS inspects a file for potential cross-site scripting vulnerabilities.
func checkXSS(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Case 1: template.HTML(variable)
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "template" && sel.Sel.Name == "HTML" {
				if len(call.Args) > 0 {
					arg := call.Args[0]
					if _, ok := arg.(*ast.BasicLit); !ok && !isIgnored(n, node, fset) {
						pos := fset.Position(arg.Pos())
						violations = append(violations, SecurityViolation{
							File:    pos.Filename,
							Line:    pos.Line,
							Column:  pos.Column,
							Type:    "XSS",
							Message: "Potential XSS: template.HTML used with non-literal value",
						})
					}
				}
			}

			// Case 2: c.HTML(code, html)
			if sel.Sel.Name == "HTML" {
				// Check if it's likely an Echo context (conventionally 'c')
				if x, ok := sel.X.(*ast.Ident); ok && (x.Name == "c" || x.Name == "ctx") {
					if len(call.Args) >= 2 {
						arg := call.Args[1]
						if isUnsafeString(arg) && !isIgnored(n, node, fset) {
							pos := fset.Position(arg.Pos())
							violations = append(violations, SecurityViolation{
								File:    pos.Filename,
								Line:    pos.Line,
								Column:  pos.Column,
								Type:    "XSS",
								Message: "Potential XSS: unsafe HTML construction in c.HTML",
							})
						}
					}
				}
			}
		}

		return true
	})

	return violations
}

// isUnsafeString checks if an expression is a concatenation or a Sprintf call.
func isUnsafeString(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return true
		}
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "fmt" && sel.Sel.Name == "Sprintf" {
				return true
			}
		}
	}
	return false
}

func isIgnored(n ast.Node, file *ast.File, fset *token.FileSet) bool {
	pos := fset.Position(n.Pos())
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			cpos := fset.Position(c.Pos())
			// Ignore if comment is on the same line or the line directly above
			if cpos.Line == pos.Line || cpos.Line == pos.Line-1 {
				if strings.Contains(c.Text, "#nosec") || strings.Contains(c.Text, "antigravity:allow") {
					return true
				}
			}
		}
	}
	return false
}

func isIgnoredRaw(line string) bool {
	return strings.Contains(line, "#nosec") || strings.Contains(line, "antigravity:allow")
}

// checkInsecurePatternsGo scans Go string literals for insecure patterns using AST.
func checkInsecurePatternsGo(node *ast.File, fset *token.FileSet) []SecurityViolation {
	var violations []SecurityViolation

	ast.Inspect(node, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		// Clean string literal (handle both regular and backticks)
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

// checkStructuralRaw scans a file for structural security issues using patterns.
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

		// Check global structural patterns
		for name, re := range structuralPatterns {
			// Skip specific checks for .go files if they are handled by AST
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

		// Check HTML-specific patterns
		if ext == ".html" {
			for name, re := range htmlPatterns {
				if re.MatchString(line) {
					// Special handling for Inline Script: Check if it's NOT an external script
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
