package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func GateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gate <gate_id> <PENDING|PASS|FAIL>",
		Short: "Run the validation gate for the current phase",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			gateID := args[0]
			statusStr := args[1]
			if statusStr != "PENDING" && statusStr != "PASS" && statusStr != "FAIL" && statusStr != "PASSED" && statusStr != "FAILED" {
				return fmt.Errorf("invalid status '%s'", statusStr)
			}

			var status agent.GateStatus
			switch statusStr {
			case "PENDING":
				status = agent.GatePending
			case "PASS", "PASSED":
				if gateID == agent.GateRedTest || gateID == agent.GateCoverage || gateID == agent.GateLint || gateID == agent.GateImplementation || gateID == agent.GateApiSpec {
					fmt.Fprintf(os.Stderr, "❌ Error: The '%s' gate cannot be manually bypassed.\n", gateID)
					if f, err := os.OpenFile(".tester/tasks/bypass_audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err == nil {
						_, _ = fmt.Fprintf(f, "%s | BLOCKED | gate=%s | attempted=%s | user=agent\n", time.Now().UTC().Format(time.RFC3339), gateID, statusStr)
						_ = f.Close()
					}
					return fmt.Errorf("manual bypass not allowed for gate '%s'", gateID)
				}
				status = agent.GatePassed
			case "FAIL", "FAILED":
				status = agent.GateFailed
			}

			state, err := getState()
			if err != nil {
				return err
			}

			switch gateID {
			case agent.GateRedTest:
				state.Gates.RedTest = status
			case agent.GateApiSpec:
				state.Gates.ApiSpec = status
			case agent.GateImplementation:
				state.Gates.Implementation = status
			case agent.GateLint:
				state.Gates.Lint = status
			case agent.GateCoverage:
				state.Gates.Coverage = status
			case agent.GateBrowserVerification:
				state.Gates.BrowserVerification = status
			default:
				return fmt.Errorf("unknown gate '%s'", gateID)
			}

			if err := saveState(state); err != nil {
				return err
			}

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
			return nil
		},
	}
}
