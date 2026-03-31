package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func VerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <gate_id> [pattern]",
		Short: "Run the validation gate for the current phase",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			gateID := args[0]
			pattern := ""
			if len(args) > 1 {
				pattern = args[1]
			}
			state, err := getState()
			if err != nil {
				return err
			}

			// Allow drift checks to run without an active feature (useful for CI)
			isDriftCheck := gateID == agent.GateApiSpec || gateID == agent.GateTemplateDrift || gateID == agent.GateSecurityStatic

			if state.Feature == "" && !isDriftCheck {
				return fmt.Errorf("no active feature found in state file")
			}
			if state.WorkflowType == "" {
				state.WorkflowType = agent.WorkflowFeature
			}

			if flagText {
				fmt.Printf("Verifying gate: %s for feature: %s (%s) [%s]\n", gateID, state.Feature, state.Phase, state.WorkflowType)
			}

			// Dependency Checks
			switch gateID {
			case agent.GateImplementation:
				if state.Gates.RedTest != agent.GatePassed || state.Gates.ApiSpec != agent.GatePassed {
					return fmt.Errorf("implementation requires red-test and api-spec to be PASS. 💡 HINT: If this is a UI layer change, use './scripts/agent-exec.sh verify red-test ui-bypass'. 💡 HINT: Note that you MUST still pass the lint gate and verify the UI using the browser_subagent")
				}
			case agent.GateLint, agent.GateCoverage, agent.GateBrowserVerification:
				if state.Gates.Implementation != agent.GatePassed {
					return fmt.Errorf("%s requires implementation to be PASS", gateID)
				}
			}

			success := false
			switch gateID {
			case agent.GateRedTest:
				success = agent.VerifyRedTest(pattern)
			case agent.GateApiSpec:
				success = agent.VerifyApiSpec(state.WorkflowType)
			case agent.GateImplementation:
				success = agent.VerifyImplementation()
			case agent.GateLint:
				success = agent.VerifyLint()
			case agent.GateCoverage:
				success = agent.VerifyCoverage()
			case agent.GateTemplateDrift:
				success = agent.VerifyTemplateDrift()
			case agent.GateSecurityStatic:
				targets := []string{"cmd", "internal"}
				if pattern != "" {
					targets = strings.Fields(pattern)
				}
				success = agent.VerifySecurityStaticGate(targets...)
			case agent.GateBrowserVerification:
				fmt.Println("⚠️  AGENT INSTRUCTION: You must use the browser_subagent tool to verify the UI. Once the subagent finishes and the UI is verified, run: ./scripts/agent-exec.sh workflow gate browser-verification PASS")
				if state.Gates.BrowserVerification == agent.GatePassed {
					fmt.Println("✅ Gate PASS: browser-verification already marked as PASS.")
					success = true
				} else {
					fmt.Println("❌ Gate FAIL: browser-verification must be manually passed or verified via browser subagent.")
					success = false
				}
			default:
				return fmt.Errorf("unknown gate_id '%s'", gateID)
			}

			// Update gate status
			var status agent.GateStatus
			statusStr := "FAIL"
			if success {
				status = agent.GatePassed
				statusStr = "PASS"
			} else {
				status = agent.GateFailed
			}

			// Save in local state
			switch gateID {
			case agent.GateRedTest:
				state.Gates.RedTest = status
				if success {
					fmt.Printf("🚀 Gate PASS: red-test passed. Transitioning to GREEN phase...\n")
				}
			case agent.GateApiSpec:
				state.Gates.ApiSpec = status
			case agent.GateImplementation:
				state.Gates.Implementation = status
				if success {
					fmt.Printf("🚀 Gate PASS: implementation passed. Proceeding to audit...\n")
				}
			case agent.GateLint:
				state.Gates.Lint = status
			case agent.GateCoverage:
				state.Gates.Coverage = status
			case agent.GateTemplateDrift:
				state.Gates.TemplateDrift = status
			case agent.GateSecurityStatic:
				state.Gates.SecurityStatic = status
			case agent.GateBrowserVerification:
				// Already handled above
			}

			if gateID != agent.GateBrowserVerification {
				if err = saveState(state); err != nil {
					return err
				}
			}

			// This matches update_gate in bash script
			// #nosec G204 - Internal harness tool calling itself
			c := exec.Command("scripts/agent-exec.sh", "workflow", "gate", gateID, statusStr)
			_ = c.Run()

			// Auto transition
			state, err = getState() // Reload as exec may have updated it
			if err != nil {
				return err
			}
			switch state.Phase {
			case "RED":
				if state.Gates.RedTest == agent.GatePassed && state.Gates.ApiSpec == agent.GatePassed {
					if flagText {
						fmt.Println("✨ All RED gates passed. Transitioning to GREEN phase.")
					}
					// #nosec G204 - Internal harness tool calling itself
					_ = exec.Command("scripts/agent-exec.sh", "workflow", "set-phase", "GREEN").Run()
				}
			case "GREEN":
				if state.Gates.Implementation == agent.GatePassed {
					if flagText {
						fmt.Println("✨ Implementation passed. Transitioning to REFACTOR phase.")
					}
					// #nosec G204 - Internal harness tool calling itself
					_ = exec.Command("scripts/agent-exec.sh", "workflow", "set-phase", "REFACTOR").Run()
					if uErr := checkAndApplyProgressUpdate(); uErr != nil {
						return fmt.Errorf("failed to update progress: %w", uErr)
					}
				}
			}

			// Print current status
			state, err = getState()
			if err != nil {
				return err
			}
			if flagText {
				fmt.Println("--- Current Workflow Status ---")
				b, _ := json.Marshal(state.Gates)
				fmt.Printf("Feature: %s [%s] (%s)\nGates: %s\n", state.Feature, state.WorkflowType, state.Phase, string(b))
			} else {
				printJSON(success, "verify", map[string]any{
					"message": "Gate verification completed",
					"gate":    gateID,
					"state":   state,
				}, nil)
			}

			if !success {
				return fmt.Errorf("gate verification failed: %s", gateID)
			}

			if err = agent.ArchivePassedCategories(".tester/tasks/progress.md", ".tester/tasks/progress_archive.md", 20); err != nil {
				fmt.Printf("⚠️  Warning: failed to archive progress: %v\n", err)
			}

			return nil
		},
	}
}
