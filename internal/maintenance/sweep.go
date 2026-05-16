package maintenance

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// SweepResult represents the outcome of a single maintenance gate.
type SweepResult struct {
	Gate    string `json:"gate"`
	Status  string `json:"status"` // "PASS", "FAIL", "WARN"
	Details string `json:"details"`
}

// RunSweep invokes all structural and meta gates and collects their results.
func RunSweep(rootDir string) ([]SweepResult, error) {
	var results []SweepResult

	// 1. Build check
	results = append(results, runBuildCheck(rootDir))

	// 2. Doc Drift
	results = append(results, runDocDriftCheck(rootDir))

	// 3. Deprecated Patterns
	results = append(results, runDeprecatedCheck(rootDir))

	// 4. Skill Conformance
	results = append(results, runSkillConformanceCheck(rootDir))

	// 5. Skill Resolvability
	results = append(results, runResolvableCheck(rootDir))

	// 6. Template Contract
	results = append(results, runTemplateContractCheck(rootDir))

	// 7. Agents Coverage
	results = append(results, runAgentsCoverageCheck(rootDir))

	// 8. Context Cost
	results = append(results, runContextCostCheck(rootDir))

	return results, nil
}

func runBuildCheck(rootDir string) SweepResult {
	// #nosec G204 -- maintenance utility runs build
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = rootDir
	if err := cmd.Run(); err != nil {
		return SweepResult{Gate: "build", Status: "FAIL", Details: "Build failed"}
	}
	return SweepResult{Gate: "build", Status: "PASS", Details: ""}
}

func runDocDriftCheck(rootDir string) SweepResult {
	violations, err := CheckDocDrift(rootDir)
	if err != nil {
		return SweepResult{Gate: "doc-drift", Status: "FAIL", Details: err.Error()}
	}
	if len(violations) > 0 {
		return SweepResult{Gate: "doc-drift", Status: "FAIL", Details: fmt.Sprintf("%d violations", len(violations))}
	}
	return SweepResult{Gate: "doc-drift", Status: "PASS", Details: ""}
}

func runDeprecatedCheck(rootDir string) SweepResult {
	violations, err := CheckDeprecatedPatterns(rootDir)
	if err != nil {
		return SweepResult{Gate: "deprecated", Status: "FAIL", Details: err.Error()}
	}
	if len(violations) > 0 {
		return SweepResult{Gate: "deprecated", Status: "WARN", Details: fmt.Sprintf("%d remaining", len(violations))}
	}
	return SweepResult{Gate: "deprecated", Status: "PASS", Details: ""}
}

func runSkillConformanceCheck(rootDir string) SweepResult {
	skillsDir := filepath.Join(rootDir, ".agents/skills")
	violations := SkillConformance(skillsDir)
	if len(violations) > 0 {
		return SweepResult{Gate: "skill-conformance", Status: "FAIL", Details: fmt.Sprintf("%d violations", len(violations))}
	}
	return SweepResult{Gate: "skill-conformance", Status: "PASS", Details: ""}
}

func runResolvableCheck(rootDir string) SweepResult {
	skillsDir := filepath.Join(rootDir, ".agents/skills")
	resolverPath := filepath.Join(rootDir, ".agents/skills/RESOLVER.md")
	manifestPath := filepath.Join(rootDir, ".agents/verify-manifest.yaml")
	violations := CheckResolvable(skillsDir, resolverPath, manifestPath)
	if len(violations) > 0 {
		return SweepResult{Gate: "check-resolvable", Status: "FAIL", Details: fmt.Sprintf("%d violations", len(violations))}
	}
	return SweepResult{Gate: "check-resolvable", Status: "PASS", Details: ""}
}

func runTemplateContractCheck(rootDir string) SweepResult {
	violations, err := CheckTemplateKeyGaps(rootDir)
	if err != nil {
		return SweepResult{Gate: "template-contract", Status: "FAIL", Details: err.Error()}
	}
	if len(violations) > 0 {
		return SweepResult{Gate: "template-contract", Status: "FAIL", Details: fmt.Sprintf("%d violations", len(violations))}
	}
	return SweepResult{Gate: "template-contract", Status: "PASS", Details: ""}
}

func runAgentsCoverageCheck(rootDir string) SweepResult {
	missing, _, err := CheckAgentsCoverage(rootDir)
	if err != nil {
		return SweepResult{Gate: "agents-coverage", Status: "FAIL", Details: err.Error()}
	}
	if len(missing) > 0 {
		return SweepResult{Gate: "agents-coverage", Status: "FAIL", Details: fmt.Sprintf("%d missing AGENTS.md", len(missing))}
	}
	return SweepResult{Gate: "agents-coverage", Status: "PASS", Details: ""}
}

func runContextCostCheck(rootDir string) SweepResult {
	report, err := CalculateContextCost(rootDir)
	if err != nil {
		return SweepResult{Gate: "context-cost", Status: "FAIL", Details: err.Error()}
	}
	return SweepResult{Gate: "context-cost", Status: "PASS", Details: fmt.Sprintf("RMS: %.0f", report.RMS)}
}
