package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"github.com/spf13/cobra"
)

const stateFile = ".agents/state.json"

var flagText bool

type CommandOutput struct {
	Success  bool     `json:"success"`
	Command  string   `json:"command"`
	Output   any      `json:"output"`
	Warnings []string `json:"warnings"`
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
	if err := util.SafeMkdir(".agents"); err != nil {
		return fmt.Errorf("error creating .agent directory: %w", err)
	}
	if err := agent.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("error saving state: %w", err)
	}
	return nil
}


func summarizeProgress() error {
	data, err := util.SafeReadFile(".tester/tasks/progress.json")
	if err != nil {
		return err
	}
	var tracker agent.ProgressTracker
	if err := json.Unmarshal(data, &tracker); err != nil {
		return err
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
	return nil
}

func checkAndApplyProgressUpdate() error {
	updateFile := ".tester/tasks/pending_update.json"
	targetFile := ".tester/tasks/progress.json"

	if _, err := util.SafeStat(updateFile); util.SafeIsNotExist(err) {
		return nil // No update file provided
	}

	fmt.Println("📦 Found pending_update.json. Triggering automatic progress tracker update...")
	updateData, err := util.SafeReadFile(updateFile)
	if err != nil {
		return fmt.Errorf("failed to read pending update: %w", err)
	}

	var newFeature agent.Feature
	err = json.Unmarshal(updateData, &newFeature)
	if err != nil {
		return fmt.Errorf("failed to parse pending update JSON: %w", err)
	}
	newFeature.Passes = !agent.HasPending(newFeature.Steps)

	targetData, err := util.SafeReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read progress.json: %w", err)
	}

	var tracker agent.ProgressTracker
	err = json.Unmarshal(targetData, &tracker)
	if err != nil {
		return fmt.Errorf("failed to parse progress.json: %w", err)
	}

	merged := false
	for i, f := range tracker.Features {
		if f.Category == newFeature.Category && newFeature.Category != "" {
			// Merge steps
			f.Steps = append(f.Steps, newFeature.Steps...)
			f.Passes = !agent.HasPending(f.Steps)
			tracker.Features[i] = f
			merged = true
			break
		}
	}

	if !merged {
		tracker.Features = append(tracker.Features, newFeature)
	}

	outData, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode updated progress.json: %w", err)
	}
	newFeature.Passes = !agent.HasPending(newFeature.Steps)

	if err := util.SafeWriteFile(targetFile, outData); err != nil {
		return fmt.Errorf("failed to save updated progress.json: %w", err)
	}

	fmt.Println("✅ Successfully updated progress.json with new feature implementation!")
	_ = util.SafeRemove(updateFile)
	return nil
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
