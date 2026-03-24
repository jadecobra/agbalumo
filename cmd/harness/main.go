package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

const stateFile = ".agents/state.json"

var flagText bool

type CommandOutput struct {
	Success  bool           `json:"success"`
	Command  string         `json:"command"`
	Output   any            `json:"output"`
	Warnings []string       `json:"warnings"`
}

func printJSON(success bool, command string, output any, warnings []string) {
	if warnings == nil {
		warnings = []string{}
	}
	out := CommandOutput{
		Success:  success,
		Command:  command,
		Output:   output,
		Warnings: warnings,
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

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
	if err := os.MkdirAll(".agents", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .agent directory: %v\n", err)
		os.Exit(1)
	}
	if err := agent.SaveState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		os.Exit(1)
	}
}

func hasPending(steps interface{}) bool {
	if sList, ok := steps.([]interface{}); ok {
		for _, step := range sList {
			if s, ok := step.(string); ok && strings.Contains(s, "(Pending)") {
				return true
			}
		}
	}
	return false
}

func summarizeProgress() {
	data, err := os.ReadFile(".tester/tasks/progress.json")
	if err != nil {
		return
	}
	var tracker struct {
		Features []struct {
			Category string `json:"category"`
			Passes   bool   `json:"passes"`
		} `json:"features"`
	}
	if err := json.Unmarshal(data, &tracker); err != nil {
		return
	}

	passed, pending := 0, 0
	var pendingCategories []string

	for _, f := range tracker.Features {
		if f.Passes {
			passed++
		} else {
			pending++
			pendingCategories = append(pendingCategories, f.Category)
		}
	}

	if flagText {
		fmt.Printf("\n--- Project Progress Summary ---\n")
		fmt.Printf("Total Tracked Goals: %d\n", len(tracker.Features))
		fmt.Printf("✅ Passed / Completed: %d\n", passed)
		fmt.Printf("⏳ Pending / In-Progress: %d\n", pending)
		if pending > 0 {
			fmt.Printf("Pending Categories:\n")
			for _, cat := range pendingCategories {
				fmt.Printf("  - %s\n", cat)
			}
		}
		fmt.Printf("--------------------------------\n\n")
	}
}

func checkAndApplyProgressUpdate() {
	updateFile := ".tester/tasks/pending_update.json"
	targetFile := ".tester/tasks/progress.json"

	if _, err := os.Stat(updateFile); os.IsNotExist(err) {
		return // No update file provided
	}

	fmt.Println("📦 Found pending_update.json. Triggering automatic progress tracker update...")
	updateData, err := os.ReadFile(updateFile)
	if err != nil {
		fmt.Println("⚠️ Failed to read pending update:", err)
		return
	}

	var newFeature map[string]interface{}
	err = json.Unmarshal(updateData, &newFeature)
	if err != nil {
		fmt.Println("⚠️ Failed to parse pending update JSON:", err)
		return
	}
	newFeature["passes"] = !hasPending(newFeature["steps"])

	targetData, err := os.ReadFile(targetFile)
	if err != nil {
		fmt.Println("⚠️ Failed to read progress.json:", err)
		return
	}

	var tracker map[string]interface{}
	err = json.Unmarshal(targetData, &tracker)
	if err != nil {
		fmt.Println("⚠️ Failed to parse progress.json:", err)
		return
	}

	if features, ok := tracker["features"].([]interface{}); ok {
		merged := false
		newCategory, _ := newFeature["category"].(string)

		for i, f := range features {
			if featMap, ok := f.(map[string]interface{}); ok {
				if cat, _ := featMap["category"].(string); cat == newCategory && newCategory != "" {
					// Merge steps
					if existingSteps, ok := featMap["steps"].([]interface{}); ok {
						if newSteps, ok := newFeature["steps"].([]interface{}); ok {
							existingSteps = append(existingSteps, newSteps...)
							featMap["steps"] = existingSteps
						}
					}

					featMap["passes"] = !hasPending(featMap["steps"])

					features[i] = featMap
					merged = true
					break
				}
			}
		}

		if !merged {
			tracker["features"] = append(features, newFeature)
		} else {
			tracker["features"] = features
		}
	} else {
		fmt.Println("⚠️ progress.json missing features array")
		return
	}

	outData, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		fmt.Println("⚠️ Failed to encode updated progress.json:", err)
		return
	}

	if err := os.WriteFile(targetFile, outData, 0644); err != nil {
		fmt.Println("⚠️ Failed to save updated progress.json:", err)
		return
	}

	fmt.Println("✅ Successfully updated progress.json with new feature implementation!")
	_ = os.Remove(updateFile)
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
			if flagText {
				fmt.Printf("Workflow initialized for %s: %s\n", workflowType, feature)
				summarizeProgress()
			} else {
				printJSON(true, "init", map[string]any{
					"message":      fmt.Sprintf("Workflow initialized for %s: %s", workflowType, feature),
					"feature":      feature,
					"workflowType": workflowType,
					"state":        state,
				}, nil)
			}
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
			if flagText {
				fmt.Printf("Phase set to: %s\n", phase)
			} else {
				printJSON(true, "set-phase", map[string]any{
					"message": fmt.Sprintf("Phase set to: %s", phase),
					"phase":   phase,
					"state":   state,
				}, nil)
			}
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
			// original bash script logged the input statusStr 
			// echo "Gate '$gate' set to: $status"
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

	var verifyCmd = &cobra.Command{
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

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Print the current status of the harness",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			state := getState()
			if flagText {
				b, err := json.MarshalIndent(state, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling state: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(b))
			} else {
				printJSON(true, "status", map[string]any{"state": state}, nil)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&flagText, "text", false, "Output in human-readable text format (JSON is default)")
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
