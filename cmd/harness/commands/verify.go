package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func VerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <gate_id> [pattern]",
		Short: "Run the validation gate for the current phase",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			gateID := args[0]
			pattern := ""
			if len(args) > 1 {
				pattern = args[1]
			}
			state := getState()

			if state.Feature == "" {
				fmt.Fprintln(os.Stderr, "Error: No active feature found in state file")
				os.Exit(1)
			}

			if flagText {
				fmt.Printf("Verifying gate: %s for feature: %s (%s) [%s]\n", gateID, state.Feature, state.Phase, state.WorkflowType)
			}

			// Dependency Checks
			switch gateID {
			case "implementation":
				if state.Gates.RedTest != agent.GatePassed || state.Gates.ApiSpec != agent.GatePassed {
					fmt.Fprintln(os.Stderr, "❌ Error: 'implementation' requires 'red-test' and 'api-spec' to be PASS.")
					fmt.Fprintln(os.Stderr, "💡 HINT: If this is a UI layer change, use './scripts/agent-exec.sh verify red-test ui-bypass'.")
					fmt.Fprintln(os.Stderr, "💡 HINT: Note that you MUST still pass the lint gate and verify the UI using the browser_subagent.")
					os.Exit(1)
				}
			case "lint", "coverage", "browser-verification":
				if state.Gates.Implementation != agent.GatePassed {
					fmt.Fprintf(os.Stderr, "❌ Error: '%s' requires 'implementation' to be PASS.\n", gateID)
					os.Exit(1)
				}
			}

			success := false
			switch gateID {
			case "red-test":
				success = agent.VerifyRedTest(pattern)
			case "api-spec":
				success = agent.VerifyApiSpec(state.WorkflowType)
			case "implementation":
				success = agent.VerifyImplementation()
			case "lint":
				success = agent.VerifyLint()
			case "coverage":
				success = agent.VerifyCoverage()
			case "browser-verification":
				fmt.Println("⚠️  AGENT INSTRUCTION: You must use the browser_subagent tool to verify the UI. Once the subagent finishes and the UI is verified, run: ./scripts/agent-exec.sh workflow gate browser-verification PASS")
				if state.Gates.BrowserVerification == agent.GatePassed {
					fmt.Println("✅ Gate PASS: browser-verification already marked as PASS.")
					success = true
				} else {
					fmt.Println("❌ Gate FAIL: browser-verification must be manually passed or verified via browser subagent.")
					success = false
				}
			default:
				fmt.Fprintf(os.Stderr, "Error: Unknown gate_id '%s'\n", gateID)
				os.Exit(1)
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
			case "red-test":
				state.Gates.RedTest = status
			case "api-spec":
				state.Gates.ApiSpec = status
			case "implementation":
				state.Gates.Implementation = status
			case "lint":
				state.Gates.Lint = status
			case "coverage":
				state.Gates.Coverage = status
			case "browser-verification":
				// Already handled above
			}

			if gateID != "browser-verification" {
				saveState(state)
			}
			
			// This matches update_gate in bash script
			c := exec.Command("scripts/agent-exec.sh", "workflow", "gate", gateID, statusStr)
			_ = c.Run()

			// Auto transition
			state = getState() // Reload as exec may have updated it
			switch state.Phase {
			case "RED":
				if state.Gates.RedTest == agent.GatePassed && state.Gates.ApiSpec == agent.GatePassed {
				if flagText {
					fmt.Println("✨ All RED gates passed. Transitioning to GREEN phase.")
				}
					_ = exec.Command("scripts/agent-exec.sh", "workflow", "set-phase", "GREEN").Run()
				}
			case "GREEN":
				if state.Gates.Implementation == agent.GatePassed {
				if flagText {
					fmt.Println("✨ Implementation passed. Transitioning to REFACTOR phase.")
				}
					_ = exec.Command("scripts/agent-exec.sh", "workflow", "set-phase", "REFACTOR").Run()
					checkAndApplyProgressUpdate()
				}
			}

			// Print current status
			state = getState()
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
				os.Exit(1)
			}
		},
	}
}
