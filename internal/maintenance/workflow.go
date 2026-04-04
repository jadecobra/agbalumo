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
	// 1. Get staged files
	cmdDiff := exec.Command("git", "diff", "--cached", "--name-only")
	cmdDiff.Dir = rootDir
	stagedBytes, err := cmdDiff.Output()
	if err != nil {
		return PhaseIdle, fmt.Errorf("failed to get staged files: %w", err)
	}
	staged := strings.TrimSpace(string(stagedBytes))

	// 2. Get last commit message
	cmdLog := exec.Command("git", "log", "-1", "--pretty=%B")
	cmdLog.Dir = rootDir
	logBytes, err := cmdLog.Output()
	if err != nil {
		return PhaseIdle, fmt.Errorf("failed to get last commit: %w", err)
	}
	lastMsg := strings.TrimSpace(string(logBytes))

	// Inference Logic:
	// A. If no staged files...
	if staged == "" {
		return PhaseIdle, nil
	}

	// B. If staged files are ONLY tests...
	isOnlyTests := true
	lines := strings.Split(staged, "\n")
	for _, l := range lines {
		if l == "" {
			continue
		}
		if !strings.HasSuffix(l, "_test.go") {
			isOnlyTests = false
			break
		}
	}
	if isOnlyTests && len(lines) > 0 && lines[0] != "" {
		return PhaseRed, nil
	}

	// C. If last commit was a test...
	if strings.HasPrefix(lastMsg, "test") {
		return PhaseGreen, nil
	}

	// D. If last commit was feat/fix...
	if strings.HasPrefix(lastMsg, "feat") || strings.HasPrefix(lastMsg, "fix") {
		return PhaseRefactor, nil
	}

	return PhaseIdle, nil
}

// ExecuteGateChecks runs the quality gates appropriate for the current phase.
func ExecuteGateChecks(rootDir string, phase WorkflowPhase) error {
	fmt.Printf("🛡️ Checking Gates for Phase: %s\n", phase)

	switch phase {
	case PhaseRed:
		// RED: Test MUST fail
		fmt.Println("🔍 Verifying RED phase (staged test must fail)...")
		if err := runTests(rootDir); err == nil {
			return fmt.Errorf("RED phase violation: staged tests pass, but they should fail")
		}
		fmt.Println("✅ RED gate passed: Staged test fails as expected.")

	case PhaseGreen:
		// GREEN: Implementation MUST pass and satisfy basic contract checks
		fmt.Println("🔍 Verifying GREEN phase (implementation must pass tests)...")
		if err := runTests(rootDir); err != nil {
			return fmt.Errorf("GREEN phase violation: tests failed: %w", err)
		}

		fmt.Println("🔍 Checking API/CLI contract drift...")
		// Use default paths for drift check
		ctx := "."
		if err := checkContractDrift(ctx); err != nil {
			return fmt.Errorf("GREEN phase violation: contract drift detected: %w", err)
		}
		fmt.Println("✅ GREEN gate passed: Tests pass and contracts are in sync.")

	case PhaseRefactor:
		// REFACTOR: Must pass all quality gates including coverage and lint
		fmt.Println("🔍 Verifying REFACTOR phase (full audit)...")
		if err := runTests(rootDir); err != nil {
			return err
		}

		// Coverage check
		if err := CompareCoverageThreshold(".metrics/coverage"); err != nil {
			// fallback to legacy
			if e := CompareCoverageThreshold(".agents/coverage-threshold"); e != nil {
				return fmt.Errorf("REFACTOR phase violation: coverage threshold not met: %w", e)
			}
		}

		fmt.Println("✅ REFACTOR gate passed.")

	case PhaseIdle:
		fmt.Println("ℹ️  No active feature or idle state. Skipping gates.")
	}

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
