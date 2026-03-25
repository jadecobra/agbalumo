package commands

import (
	"encoding/json"
	"fmt"
	"os"
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

func getState() (*agent.State, error) {
	state, err := agent.LoadState(stateFile)
	if err != nil {
		if agent.IsNotExist(err) {
			return &agent.State{}, nil
		}
		return nil, fmt.Errorf("error loading state: %w", err)
	}
	return state, nil
}

func saveState(state *agent.State) error {
	if err := os.MkdirAll(".agents", 0755); err != nil {
		return fmt.Errorf("error creating .agent directory: %w", err)
	}
	if err := agent.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("error saving state: %w", err)
	}
	return nil
}

func hasPending(steps interface{}) bool {
	if sList, ok := steps.([]interface{}); ok {
		for _, step := range sList {
			if s, ok := step.(string); ok && !strings.Contains(s, "(Completed)") {
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
		Use:           "harness",
		Short:         "agentic-harness is a robust CLI for 10x Engineer workflows",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.PersistentFlags().BoolVar(&flagText, "text", false, "Output in human-readable text format (JSON is default)")
	
	rootCmd.AddCommand(InitCmd())
	rootCmd.AddCommand(SetPhaseCmd())
	rootCmd.AddCommand(GateCmd())
	rootCmd.AddCommand(VerifyCmd())
	rootCmd.AddCommand(StatusCmd())
	rootCmd.AddCommand(UpdateCoverageCmd())
	rootCmd.AddCommand(CostCmd())

	return rootCmd
}
