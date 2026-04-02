package agent

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

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
