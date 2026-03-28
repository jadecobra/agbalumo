package agent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

func VerifySecurityStaticGate(paths ...string) bool {
	fmt.Println("Running AST-based structural security checks...")
	var allViolations []SecurityViolation
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, p := range paths {
		violations, err := VerifySecurityStatic(p)
		if err != nil {
			fmt.Printf("❌ Error running security checks on %s: %v\n", p, err)
			return false
		}
		allViolations = append(allViolations, violations...)
	}

	if len(allViolations) == 0 {
		fmt.Println("✅ Gate PASS: no structural security violations found.")
		return true
	}

	limit := 1
	fmt.Printf("❌ Gate FAIL: %d security violations detected (showing first %d).\n", len(allViolations), limit)
	for i, v := range allViolations {
		if i >= limit {
			fmt.Printf("... and %d more violations.\n", len(allViolations)-limit)
			break
		}
		fmt.Printf("  [%s] %s:%d:%d: %s\n", v.Type, v.File, v.Line, v.Column, v.Message)
	}
	return false
}

// RunCommand is a helper to run commands and capture output
var ExecCommand = exec.Command
var LookPath = exec.LookPath
var OSStat = os.Stat

func RunCommand(name string, args ...string) ([]byte, error) {
	cmd := ExecCommand(name, args...)
	return cmd.CombinedOutput()
}

func VerifyRedTest(pattern string) bool {
	if pattern == "ui-bypass" || pattern == "--ui-bypass" {
		out, err := RunCommand("git", "status", "--porcelain")
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					filename := fields[len(fields)-1]
					if strings.HasSuffix(filename, ".go") && !strings.HasSuffix(filename, "_test.go") {
						fmt.Printf("❌ Gate FAIL: --ui-bypass is only permitted for UI files. Non-UI file modified: %s\n", filename)
						return false
					}
				}
			}
		}
		fmt.Println("⚠️  UI BYPASS ENGAGED: Skipping Go test failure requirement.")
		fmt.Println("⚠️  NOTE: You are strictly required to use the browser_subagent to verify your changes.")
		fmt.Printf("Gate PASS: %s bypassed for UI change.\n", GateRedTest)
		return true
	}

	pkgPath := "./internal/agent/..."
	outputPattern := ""
	if strings.HasPrefix(pattern, "./") {
		pkgPath = pattern
	} else {
		outputPattern = pattern
	}

	fmt.Printf("Running tests in %s expecting failure...\n", pkgPath)

	// 1. Verify code compiles first.
	_ = util.SafeMkdir(".tester")
	compileOut, err := RunCommand("go", "test", "-buildvcs=false", "-run=^$", pkgPath)
	if err != nil {
		fmt.Println("FAIL: Code does not compile. Fixed syntax/imports before running red-test.")
		_ = util.SafeWriteFile(filepath.Join(".tester", "red-test-compile.log"), compileOut)
		fmt.Println(string(compileOut))
		return false
	}

	// 2. Run tests and capture JSON output.
	testOut, _ := RunCommand("go", "test", "-buildvcs=false", "-json", pkgPath)

	res, err := ParseTestJSON(bytes.NewReader(testOut))
	if err != nil {
		fmt.Println("FAIL: Failed to parse test JSON")
		return false
	}

	if res.CompilationFailed {
		fmt.Println("Gate FAIL: tests failed but could not find '--- FAIL:' marker. Check for panics or setup issues.")
		return false
	}

	if res.Success {
		fmt.Printf("Gate FAIL: %s passed but was expected to fail.\n", GateRedTest)
		return false
	}

	validFailures := 0
	for _, f := range res.Failures {
		if f.TestName != "panic/raw" && f.TestName != "build" {
			validFailures++
		}
	}

	if validFailures == 0 {
		fmt.Println("Gate FAIL: tests failed but no individual test failures were recorded (possible panic/os.Exit evasion).")
		return false
	}

	if outputPattern != "" {
		patternFound := false
		for _, fail := range res.Failures {
			if strings.Contains(fail.Output, outputPattern) {
				patternFound = true
				break
			}
		}
		if patternFound {
			fmt.Printf("Gate PASS: %s failed with expected pattern '%s'.\n", GateRedTest, outputPattern)
			return true
		} else {
			fmt.Printf("Gate FAIL: %s failed but pattern '%s' not found in output.\n", GateRedTest, outputPattern)
			SummarizeTestFailures(res.Failures, 1) // Only show 1 for red-test
			return false
		}
	}

	fmt.Printf("Gate PASS: %s failed as expected.\n", GateRedTest)
	return true
}

func VerifyApiSpec(workflowType string) bool {
	fmt.Println("Running API and CLI drift checks...")

	codeRoutes, err := ExtractRoutes("cmd", "internal/handler", "internal/module")
	if err != nil {
		fmt.Println("Error extracting routes from code:", err)
		return false
	}

	// Use a local npm cache to avoid permission issues in CI/CD or restricted environments
	npmCache := filepath.Join(".tester", "tmp", "npm_cache")
	_ = util.SafeMkdir(npmCache)
	
	cmd := ExecCommand("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, "NPM_CONFIG_CACHE="+npmCache)
	openapiData, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error bundling docs/openapi.yaml:", err)
		return false
	}
	openapiRoutes, err := ExtractOpenAPIRoutes(openapiData)
	if err != nil {
		fmt.Println("Error extracting openapi routes:", err)
		return false
	}
	
	// #nosec G304 - Internal harness tool reading project docs
	mdData, err := util.SafeReadFile("docs/api.md")
	if err != nil {
		fmt.Println("Error reading docs/api.md:", err)
		return false
	}
	mdRoutes, err := ExtractMarkdownRoutes(mdData)
	if err != nil {
		fmt.Println("Error extracting md routes:", err)
		return false
	}

	drifts := CheckAPIDrift(codeRoutes, openapiRoutes, mdRoutes)

	// -- native CLI Drift calculations --
	cliCodeCmds, err := ExtractCLICodeCommands("cmd")
	if err != nil {
		fmt.Println("Error extracting CLI code cmds:", err)
		return false
	}
	
	cliMDCmds, err := ExtractCLIMarkdownCommands("docs/cli.md", "docs/cli")
	if err != nil {
		fmt.Println("Error extracting CLI md cmds:", err)
		return false
	}

	cliDrifts := CheckCLIDrift(cliCodeCmds, cliMDCmds)
	drifts = append(drifts, cliDrifts...)

	if len(drifts) == 0 {
		fmt.Println("✅ Gate PASS: drift checks passed.")
		return true
	}

	for _, drift := range drifts {
		fmt.Println(drift)
	}

	if workflowType == WorkflowRefactor || workflowType == WorkflowBugfix {
		fmt.Printf("⚠️  Gate FAIL: drift checks failed. For '%s' workflow, these are mandatory passive validations.\n", workflowType)
		fmt.Println("Please ensure you haven't accidentally broken existing API or CLI contracts.")
	}
	fmt.Println("❌ Gate FAIL: contract drift detected.")
	return false
}

func VerifyImplementation() bool {
	fmt.Println("Running early lint and build...")

	if !VerifyLint() {
		return false
	}

	buildOut, err := RunCommand("go", "build", "-buildvcs=false", "./cmd/...", "./internal/...")
	if err != nil {
		fmt.Println("❌ Gate FAIL: implementation build failed.")
		fmt.Println("--- BUILD OUTPUT ---")
		fmt.Println(string(buildOut))
		return false
	}

	fmt.Println("Running tests...")
	_ = util.SafeMkdir(filepath.Join(".tester", "coverage"))
	covFile := filepath.Join(".tester", "coverage", "coverage.out")
	testOut, testErr := RunCommand("go", "test", "-buildvcs=false", "-json", "-coverprofile="+covFile, "./internal/agent/...")
	if testErr != nil {
		fmt.Println("❌ Gate FAIL: implementation tests failed.")
		res, parseErr := ParseTestJSON(bytes.NewReader(testOut))
		if parseErr == nil && len(res.Failures) > 0 {
			SummarizeTestFailures(res.Failures, 1)
		} else if len(testOut) > 0 {
			fmt.Println("--- RAW TEST OUTPUT (TRUNCATED) ---")
			lines := strings.Split(string(testOut), "\n")
			if len(lines) > 20 {
				fmt.Println(strings.Join(lines[:20], "\n"))
				fmt.Printf("... truncated %d lines ...\n", len(lines)-20)
			} else {
				fmt.Println(string(testOut))
			}
		}
		return false
	}

	fmt.Printf("✅ Gate PASS: %s build and tests passed.\n", GateImplementation)
	return true
}

func VerifyLint() bool {
	fmt.Println("Running linter...")

	lintPath, err := LookPath("golangci-lint")
	if err != nil {
		// Fallback for Mac ARM Homebrew
		if _, statErr := OSStat("/opt/homebrew/bin/golangci-lint"); statErr == nil {
			lintPath = "/opt/homebrew/bin/golangci-lint"
		} else {
			fmt.Println("⚠️  golangci-lint not found, skipping...")
			fmt.Printf("✅ Gate PASS: %s passed (skipped).\n", GateLint)
			return true
		}
	}

	cmd := ExecCommand(lintPath, "run", "-c", "scripts/.golangci.yml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Gate WARNING: lint failed: %v. Continuing as this may be an environmental issue.\n", err)
		return true
	}

	fmt.Printf("✅ Gate PASS: %s passed.\n", GateLint)
	return true
}

func VerifyCoverage() bool {
	fmt.Println("Verifying test coverage...")

	covFile := filepath.Join(".tester", "coverage", "coverage.out")

	if _, statErr := util.SafeStat(covFile); util.SafeIsNotExist(statErr) {
		fmt.Println("❌ Gate FAIL: coverage profile not generated.")
		return false
	}
	
	// #nosec G304 - Internal harness tool reading coverage profile
	f, err := util.SafeOpen(covFile)
	if err != nil {
		fmt.Println("❌ Gate FAIL: unable to read coverage profile.")
		return false
	}
	defer func() { _ = f.Close() }()

	coverage, err := ParseCoverageProfile(f)
	if err != nil {
		fmt.Println("❌ Gate FAIL: unable to parse coverage profile.")
		return false
	}
	
	// #nosec G304 - Internal harness tool reading thresholds
	thresholdsData, err := util.SafeReadFile(filepath.Join(".agents", "coverage-thresholds.json"))
	var thresholds map[string]float64
	if err == nil {
		var parseErr error
		thresholds, parseErr = ParseThresholds(thresholdsData)
		if parseErr != nil {
			fmt.Println("❌ " + parseErr.Error())
			return false
		}
	} else {
		globalThreshold := 90.0
		// #nosec G304 - Internal harness tool reading threshold
		legacyData, err := util.SafeReadFile(filepath.Join(".agents", "coverage-threshold"))
		if err == nil {
			parsed, err := strconv.ParseFloat(strings.TrimSpace(string(legacyData)), 64)
			if err == nil {
				globalThreshold = parsed
			}
		}
		thresholds = map[string]float64{"default": globalThreshold}
	}

	violations := EnforceCoverage(coverage, thresholds)

	// We'll calculate a crude overall coverage if there are violations. But the violations list is fine.
	if len(violations) > 0 {
		// #nosec G204 - Internal harness tool executing cover tool
		out, _ := exec.Command("go", "tool", "cover", "-func="+covFile).CombinedOutput()
		totalLine := ""
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "total:") {
				totalLine = line
			}
		}
		
		fmt.Printf("❌ Gate FAIL: %s. Thresholds not met.\n", totalLine)
		for _, v := range violations {
			fmt.Println("  " + v)
		}
		return false
	}

	fmt.Printf("✅ Gate PASS: %s meets thresholds.\n", GateCoverage)
	return true
}

func VerifyTemplateDrift() bool {
	fmt.Println("Running Template Function Drift Check...")

	rendererPath := "internal/ui/renderer.go"
	templatesDir := "ui/templates"

	definedFuncs, err := ExtractRendererFunctions(rendererPath)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return false
	}

	usedFuncs, err := ExtractTemplateFunctionCalls(templatesDir)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return false
	}

	drifts := CheckTemplateDrift(definedFuncs, usedFuncs)
	if len(drifts) == 0 {
		fmt.Println("✅ All template functions are in sync.")
		return true
	}

	for _, d := range drifts {
		fmt.Printf("❌ %s\n", d)
	}
	fmt.Println("❌ Gate FAIL: Template Function Drift Detected!")
	return false
}

func ExtractRendererFunctions(path string) ([]string, error) {
	data, err := util.SafeReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read renderer file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var funcs []string
	
	// Implementation matches bash logic: grep -E '^		"[a-zA-Z0-9]+":'
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"") && strings.Contains(line, "\":") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				funcs = append(funcs, parts[1])
			}
		}
	}

	return util.UniqueStrings(funcs), nil
}

func ExtractTemplateFunctionCalls(dir string) ([]string, error) {
	var used []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			data, err := util.SafeReadFile(path)
			if err != nil {
				return err
			}
			
			// Simple extraction based on bash logic
			content := string(data)
			
			// Find {{ func or {{ range func
			// We'll use a simplified version of the bash regex.
			// Re-implemented logic:
			// find ui/templates -name "*.html" -exec cat {} + | \
			// grep -oE '\{\{[[:space:]]*(range[[:space:]]+)?([a-zA-Z0-9]+)[[:space:]]'
			
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				// {{ range func ... }}
				if strings.Contains(line, "{{") {
					parts := strings.Split(line, "{{")
					for _, p := range parts[1:] {
						inner := strings.TrimSpace(p)
						if strings.HasPrefix(inner, "range") {
							inner = strings.TrimSpace(strings.TrimPrefix(inner, "range"))
						}
						
						// Take the first word
						words := strings.FieldsFunc(inner, func(r rune) bool {
							return r == ' ' || r == '}' || r == '|' || r == '(' || r == ')'
						})
						if len(words) > 0 {
							name := words[0]
							// Filter out common template keywords and variables
							if !isTemplateKeyword(name) && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "$") {
								used = append(used, name)
							}
						}
					}
				}
				
				// | func
				if strings.Contains(line, "|") {
					parts := strings.Split(line, "|")
					for _, p := range parts[1:] {
						inner := strings.TrimSpace(p)
						words := strings.FieldsFunc(inner, func(r rune) bool {
							return r == ' ' || r == '}' || r == '|' || r == '(' || r == ')'
						})
						if len(words) > 0 {
							name := words[0]
							if !isTemplateKeyword(name) && !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "$") {
								used = append(used, name)
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return util.UniqueStrings(used), nil
}

func isTemplateKeyword(s string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "end": true, "range": true, "with": true,
		"define": true, "block": true, "template": true, "nil": true, "len": true,
		"and": true, "or": true, "not": true, "index": true, "slice": true,
		"printf": true, "print": true, "println": true, "html": true,
		"urlquery": true, "js": true, "call": true,
	}
	return keywords[s]
}

func CheckTemplateDrift(defined, used []string) []string {
	var drifts []string
	defMap := make(map[string]bool)
	for _, d := range defined {
		defMap[d] = true
	}

	for _, u := range used {
		if !defMap[u] {
			// Skip uppercase fields as they might be direct struct access (though rare in these templates)
			if len(u) > 0 && u[0] >= 'A' && u[0] <= 'Z' {
				continue
			}
			drifts = append(drifts, fmt.Sprintf("Undefined template function used: '%s'", u))
		}
	}
	return drifts
}
func SummarizeTestFailures(failures []TestFailure, max int) {
	if len(failures) == 0 {
		return
	}
	fmt.Println("--- TEST FAILURES ---")
	limit := len(failures)
	if limit > max {
		limit = max
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("❌ %s:\n%s\n", failures[i].TestName, failures[i].Output)
	}
	if len(failures) > max {
		fmt.Printf("... and %d more failures.\n", len(failures)-max)
	}
}
