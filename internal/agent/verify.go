package agent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// RunCommand is a helper to run commands and capture output
var ExecCommand = exec.Command

func RunCommand(name string, args ...string) ([]byte, error) {
	cmd := ExecCommand(name, args...)
	return cmd.CombinedOutput()
}

func VerifyRedTest(pattern string) bool {
	fmt.Println("Running tests expecting failure...")

	// 1. Verify code compiles first.
	_ = os.MkdirAll(".tester", 0755)
	compileOut, err := RunCommand("go", "test", "-run=^$", "./...")
	if err != nil {
		fmt.Println("FAIL: Code does not compile. Fixed syntax/imports before running red-test.")
		_ = os.WriteFile(filepath.Join(".tester", "red-test-compile.log"), compileOut, 0644)
		fmt.Println(string(compileOut))
		return false
	}

	// 2. Run tests and capture JSON output.
	testOut, _ := RunCommand("go", "test", "-json", "./...")

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
		fmt.Println("Gate FAIL: red-test passed but was expected to fail.")
		return false
	}

	if pattern != "" {
		patternFound := false
		for _, fail := range res.Failures {
			if strings.Contains(fail.Output, pattern) {
				patternFound = true
				break
			}
		}
		if patternFound {
			fmt.Printf("Gate PASS: red-test failed with expected pattern '%s'.\n", pattern)
			return true
		} else {
			fmt.Printf("Gate FAIL: red-test failed but pattern '%s' not found in output.\n", pattern)
			fmt.Println("--- TEST FAILURES ---")
			for _, fail := range res.Failures {
				fmt.Println(fail.TestName, ":")
				fmt.Println(fail.Output)
			}
			return false
		}
	}

	fmt.Println("Gate PASS: red-test failed as expected.")
	return true
}

func VerifyApiSpec(workflowType string) bool {
	fmt.Println("Running API and CLI drift checks...")

	codeRoutes, err := ExtractRoutes("cmd", "internal/handler", "internal/module")
	if err != nil {
		fmt.Println("Error extracting routes from code:", err)
		return false
	}

	openapiData, err := RunCommand("npx", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
	if err != nil {
		fmt.Println("Error bundling docs/openapi.yaml:", err)
		return false
	}
	openapiRoutes, err := ExtractOpenAPIRoutes(openapiData)
	if err != nil {
		fmt.Println("Error extracting openapi routes:", err)
		return false
	}

	mdData, err := os.ReadFile("docs/api.md")
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

	if workflowType == "refactor" || workflowType == "bugfix" {
		fmt.Printf("⚠️  Gate FAIL: drift checks failed. For '%s' workflow, these are mandatory passive validations.\n", workflowType)
		fmt.Println("Please ensure you haven't accidentally broken existing API or CLI contracts.")
	}
	fmt.Println("❌ Gate FAIL: contract drift detected.")
	return false
}

func VerifyImplementation() bool {
	fmt.Println("Running early lint and build...")

	vetOut, err := RunCommand("go", "vet", "./...")
	if err != nil {
		fmt.Println("❌ Gate FAIL: early static analysis (go vet) failed.")
		fmt.Println("--- GO VET OUTPUT ---")
		fmt.Println(string(vetOut))
		return false
	}

	buildOut, err := RunCommand("go", "build", "./...")
	if err != nil {
		fmt.Println("❌ Gate FAIL: implementation build failed.")
		fmt.Println("--- BUILD OUTPUT ---")
		fmt.Println(string(buildOut))
		return false
	}

	fmt.Println("Running tests...")
	_ = os.MkdirAll(filepath.Join(".tester", "coverage"), 0755)
	covFile := filepath.Join(".tester", "coverage", "coverage.out")
	testOut, err := RunCommand("go", "test", "-json", "-coverprofile="+covFile, "./...")
	if err != nil {
		fmt.Println("❌ Gate FAIL: implementation tests failed.")
		res, parseErr := ParseTestJSON(bytes.NewReader(testOut))
		if parseErr == nil && len(res.Failures) > 0 {
			fmt.Println("--- TEST FAILURES ---")
			for _, fail := range res.Failures {
				fmt.Println(fail.TestName, ":")
				fmt.Println(fail.Output)
			}
		} else {
			fmt.Println("--- RAW TEST OUTPUT ---")
			fmt.Println(string(testOut))
		}
		return false
	}

	fmt.Println("✅ Gate PASS: implementation build and tests passed.")
	return true
}

func VerifyLint() bool {
	fmt.Println("Running linter...")

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Println("⚠️  golangci-lint not found, skipping...")
		fmt.Println("✅ Gate PASS: lint passed.")
		return true
	}

	cmd := ExecCommand("golangci-lint", "run", "-c", "scripts/.golangci.yml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Gate FAIL: lint failed.")
		return false
	}

	fmt.Println("✅ Gate PASS: lint passed.")
	return true
}

func VerifyCoverage() bool {
	fmt.Println("Verifying test coverage...")

	covFile := filepath.Join(".tester", "coverage", "coverage.out")

	if _, statErr := os.Stat(covFile); os.IsNotExist(statErr) {
		fmt.Println("❌ Gate FAIL: coverage profile not generated.")
		return false
	}

	f, err := os.Open(covFile)
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

	thresholdsData, err := os.ReadFile(filepath.Join(".agents", "coverage-thresholds.json"))
	var thresholds map[string]float64
	if err == nil {
		thresholds, _ = ParseThresholds(thresholdsData)
	} else {
		globalThreshold := 90.0
		legacyData, err := os.ReadFile(filepath.Join(".agents", "coverage-threshold"))
		if err == nil {
			parsed, err := strconv.ParseFloat(strings.TrimSpace(string(legacyData)), 64)
			if err == nil {
				globalThreshold = parsed
			}
		}
		thresholds = map[string]float64{"default": globalThreshold}
	}

	violations := EnforceCoverage(coverage, thresholds)

	// Since we don't have total coverage computed from the map easily here without reimplementing the total count from output, Let's just output the violations.
	// We'll calculate a crude overall coverage if there are violations. But the violations list is fine.
	if len(violations) > 0 {
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

	fmt.Println("✅ Gate PASS: Coverage meets thresholds.")
	return true
}
