package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

const stateFile = ".agent/state.json"

func getState() *agent.State {
	state, err := agent.LoadState(stateFile)
	if err != nil {
		if agent.IsNotExist(err) {
			return &agent.State{}
		}
		fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
		os.Exit(1)
	}
	return state
}

func saveState(state *agent.State) {
	os.MkdirAll(".agent", 0755)
	if err := agent.SaveState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "harness",
		Short: "agentic-harness is a robust CLI for 10x Engineer workflows",
	}

	var initCmd = &cobra.Command{
		Use:   "init <feature_name> [workflow_type]",
		Short: "Initialize the harness tracking structure",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			feature := args[0]
			workflowType := "feature"
			if len(args) > 1 {
				workflowType = args[1]
			}
			if workflowType != "feature" && workflowType != "bugfix" && workflowType != "refactor" {
				fmt.Fprintf(os.Stderr, "Error: invalid workflow type '%s'\n", workflowType)
				os.Exit(1)
			}

			state := &agent.State{
				Feature:      feature,
				WorkflowType: workflowType,
				Phase:        "IDLE",
				Gates: agent.Gates{
					RedTest:             agent.GatePending,
					ApiSpec:             agent.GatePending,
					Implementation:      agent.GatePending,
					Lint:                agent.GatePending,
					Coverage:            agent.GatePending,
					BrowserVerification: agent.GatePending,
				},
			}
			saveState(state)
			fmt.Printf("Workflow initialized for %s: %s\n", workflowType, feature)
		},
	}

	var setPhaseCmd = &cobra.Command{
		Use:   "set-phase <IDLE|RED|GREEN|REFACTOR>",
		Short: "Set the current workflow phase",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			phase := args[0]
			if phase != "IDLE" && phase != "RED" && phase != "GREEN" && phase != "REFACTOR" {
				fmt.Fprintf(os.Stderr, "Error: invalid phase '%s'\n", phase)
				os.Exit(1)
			}
			state := getState()
			state.Phase = phase
			saveState(state)
			fmt.Printf("Phase set to: %s\n", phase)
		},
	}

	var gateCmd = &cobra.Command{
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
			// original bash script logged the input statusStr 
			// echo "Gate '$gate' set to: $status"
			fmt.Printf("Gate '%s' set to: %s\n", gateID, statusStr)
		},
	}

	var verifyCmd = &cobra.Command{
		Use:   "verify <gate_id>",
		Short: "Run the verification script for a gate",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gate := args[0]
			c := exec.Command("bash", "scripts/agent-gate.sh", gate)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
				fmt.Fprintf(os.Stderr, "Error running gate script: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Print the current status of the harness",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			state := getState()
			b, err := json.MarshalIndent(state, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling state: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
		},
	}

	rootCmd.AddCommand(initCmd, setPhaseCmd, gateCmd, verifyCmd, statusCmd)

	return rootCmd
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
