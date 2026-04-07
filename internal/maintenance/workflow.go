package maintenance

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorkflowPhase represents the current TDD/Agent workflow state.
type WorkflowPhase string

const (
	PhaseRed      WorkflowPhase = "RED"
	PhaseGreen    WorkflowPhase = "GREEN"
	PhaseRefactor WorkflowPhase = "REFACTOR"
	PhaseIdle     WorkflowPhase = "IDLE"
)

// InferCurrentPhase detects the workflow state from Git history and staged changes.
func InferCurrentPhase(rootDir string) (WorkflowPhase, error) {
	staged, err := getStagedFiles(rootDir)
	if err != nil {
		return PhaseIdle, err
	}
	if staged == "" {
		return PhaseIdle, nil
	}

	if isOnlyTestsStaged(staged) {
		return PhaseRed, nil
	}

	lastMsg, err := getLastCommitMsg(rootDir)
	if err != nil {
		return PhaseIdle, err
	}

	switch {
	case strings.HasPrefix(lastMsg, "test"):
		return PhaseGreen, nil
	case strings.HasPrefix(lastMsg, "feat"), strings.HasPrefix(lastMsg, "fix"):
		return PhaseRefactor, nil
	default:
		return PhaseIdle, nil
	}
}

func getStagedFiles(rootDir string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = rootDir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func getLastCommitMsg(rootDir string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = rootDir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func isOnlyTestsStaged(staged string) bool {
	lines := strings.Split(staged, "\n")
	if len(lines) == 0 || lines[0] == "" {
		return false
	}
	for _, l := range lines {
		if l != "" && !strings.HasSuffix(l, "_test.go") {
			return false
		}
	}
	return true
}

// ExecuteGateChecks runs the quality gates appropriate for the current phase.
func ExecuteGateChecks(rootDir string, phase WorkflowPhase) error {
	fmt.Printf("🛡️ Checking Gates for Phase: %s\n", phase)

	switch phase {
	case PhaseRed:
		return verifyRedGate(rootDir)
	case PhaseGreen:
		return verifyGreenGate(rootDir)
	case PhaseRefactor:
		return verifyRefactorGate(rootDir)
	case PhaseIdle:
		fmt.Println("ℹ️  No active feature or idle state. Skipping gates.")
	}

	return nil
}

func verifyRedGate(rootDir string) error {
	fmt.Println("🔍 Verifying RED phase (staged test must fail)...")
	if err := runTests(rootDir); err == nil {
		return fmt.Errorf("RED phase violation: staged tests pass, but they should fail")
	}
	fmt.Println("✅ RED gate passed: Staged test fails as expected.")
	return nil
}

func verifyGreenGate(rootDir string) error {
	fmt.Println("🔍 Verifying GREEN phase (implementation must pass tests)...")
	if err := runTests(rootDir); err != nil {
		return fmt.Errorf("GREEN phase violation: tests failed: %w", err)
	}
	fmt.Println("🔍 Checking API/CLI contract drift...")
	if err := checkContractDrift("."); err != nil {
		return fmt.Errorf("GREEN phase violation: contract drift detected: %w", err)
	}
	fmt.Println("✅ GREEN gate passed: Tests pass and contracts are in sync.")
	return nil
}

func verifyRefactorGate(rootDir string) error {
	fmt.Println("🔍 Verifying REFACTOR phase (full audit)...")
	if err := runTests(rootDir); err != nil {
		return err
	}
	if err := CompareCoverageThreshold(".metrics/coverage"); err != nil {
		if e := CompareCoverageThreshold(".agents/coverage-threshold"); e != nil {
			return fmt.Errorf("REFACTOR phase violation: coverage threshold not met: %w", e)
		}
	}
	fmt.Println("✅ REFACTOR gate passed.")
	return nil
}

func runTests(rootDir string) error {
	// We use -short to avoid long-running benchmarks if any
	cmd := exec.Command("go", "test", "-short", "-count=1", "./...")
	cmd.Dir = rootDir
	return cmd.Run()
}

func checkContractDrift(rootDir string) error {
	// 1. Code Routes
	codeRoutes, err := ExtractRoutes(filepath.Join(rootDir, "cmd"), filepath.Join(rootDir, "internal/handler"), filepath.Join(rootDir, "internal/module"), filepath.Join(rootDir, "internal/infra"))
	if err != nil {
		return err
	}

	// 2. OpenAPI Routes
	bundleCmd := exec.Command("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
	bundleCmd.Dir = rootDir
	openapiData, err := bundleCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to bundle openapi.yaml: %w", err)
	}
	openapiRoutes, err := ExtractOpenAPIRoutes(openapiData)
	if err != nil {
		return err
	}

	drifts := CompareRoutes("Code", "OpenAPI", codeRoutes, openapiRoutes)
	if len(drifts) > 0 {
		return fmt.Errorf("contract drift detected: %s", strings.Join(drifts, "; "))
	}
	return nil
}
