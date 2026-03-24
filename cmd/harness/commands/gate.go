package commands

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func GateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gate <gate_id> <PENDING|PASS|FAIL>",
		Short: "Run the validation gate for the current phase",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			gateID := args[0]
			statusStr := args[1]
			if statusStr != "PENDING" && statusStr != "PASS" && statusStr != "FAIL" && statusStr != "PASSED" && statusStr != "FAILED" {
				fmt.Fprintf(os.Stderr, "Error: invalid status '%s'\n", statusStr)
				os.Exit(1)
			}
			
			var status agent.GateStatus
			switch statusStr {
			case "PENDING":
				status = agent.GatePending
			case "PASS", "PASSED":
				if gateID == "red-test" || gateID == "coverage" || gateID == "lint" || gateID == "implementation" || gateID == "api-spec" {
					fmt.Fprintf(os.Stderr, "❌ Error: The '%s' gate cannot be manually bypassed.\n", gateID)
					fmt.Fprintln(os.Stderr, "💡 HINT: You must pass this gate through automated verification by fixing the code/adding tests and running: ./scripts/agent-exec.sh verify "+gateID)
					os.Exit(1)
				}
				status = agent.GatePassed
			case "FAIL", "FAILED":
				status = agent.GateFailed
			}

			state := getState()
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
				state.Gates.BrowserVerification = status
			default:
				fmt.Fprintf(os.Stderr, "Error: unknown gate '%s'\n", gateID)
				os.Exit(1)
			}
			saveState(state)
			if flagText {
				fmt.Printf("Gate '%s' set to: %s\n", gateID, statusStr)
			} else {
				printJSON(true, "gate", map[string]any{
					"message": fmt.Sprintf("Gate '%s' set to: %s", gateID, statusStr),
					"gate":    gateID,
					"status":  statusStr,
					"state":   state,
				}, nil)
			}
		},
	}
}
