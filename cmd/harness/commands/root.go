package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	progressPath := ".tester/tasks/progress.md"
	if _, err := util.SafeStat(progressPath); util.SafeIsNotExist(err) {
		// Fallback to JSON during migration
		progressPath = ".tester/tasks/progress.json"
	}

	data, err := util.SafeReadFile(progressPath)
	if err != nil {
		return err
	}

	var tracker agent.ProgressTracker
	if strings.HasSuffix(progressPath, ".md") {
		tracker, err = agent.ParseMarkdownTracker(string(data))
		if err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &tracker); err != nil {
			return err
		}
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
	updateFileMD := ".tester/tasks/pending_update.md"
	updateFileJSON := ".tester/tasks/pending_update.json"
	targetFileMD := ".tester/tasks/progress.md"
	targetFileJSON := ".tester/tasks/progress.json"

	var updateFile string
	if _, err := util.SafeStat(updateFileMD); err == nil {
		updateFile = updateFileMD
	} else if _, err := util.SafeStat(updateFileJSON); err == nil {
		updateFile = updateFileJSON
	} else {
		return nil // No update file provided
	}

	fmt.Printf("📦 Found %s. Triggering automatic progress tracker update...\n", updateFile)
	updateData, err := util.SafeReadFile(updateFile)
	if err != nil {
		return fmt.Errorf("failed to read pending update: %w", err)
	}

	var newFeature agent.Feature
	if strings.HasSuffix(updateFile, ".md") {
		tempTracker, tErr := agent.ParseMarkdownTracker(string(updateData))
		if tErr != nil || len(tempTracker.Features) == 0 {
			return fmt.Errorf("failed to parse pending update Markdown: %w", tErr)
		}
		newFeature = tempTracker.Features[0]
	} else {
		err = json.Unmarshal(updateData, &newFeature)
		if err != nil {
			return fmt.Errorf("failed to parse pending update JSON: %w", err)
		}
	}
	newFeature.Passes = !agent.HasPending(newFeature.Steps)

	// Determine target file
	targetFile := targetFileMD
	if _, sErr := util.SafeStat(targetFileMD); util.SafeIsNotExist(sErr) {
		if _, sErr := util.SafeStat(targetFileJSON); sErr == nil {
			targetFile = targetFileJSON
		}
	}

	targetData, err := util.SafeReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read target progress file: %w", err)
	}

	var tracker agent.ProgressTracker
	if strings.HasSuffix(targetFile, ".md") {
		tracker, err = agent.ParseMarkdownTracker(string(targetData))
		if err != nil {
			return fmt.Errorf("failed to parse progress Markdown: %w", err)
		}
	} else {
		err = json.Unmarshal(targetData, &tracker)
		if err != nil {
			return fmt.Errorf("failed to parse progress JSON: %w", err)
		}
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

	var outData []byte
	if strings.HasSuffix(targetFile, ".md") {
		outData = []byte(agent.ToMarkdown(tracker))
	} else {
		outData, err = json.MarshalIndent(tracker, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode updated progress file: %w", err)
		}
	}

	if err := util.SafeWriteFile(targetFile, outData); err != nil {
		return fmt.Errorf("failed to save updated progress file: %w", err)
	}

	fmt.Printf("✅ Successfully updated %s with new implementation!\n", targetFile)
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

	rootCmd.PersistentFlags().BoolVar(&flagText, "text", os.Getenv("HARNESS_TEXT") == "true", "Output in human-readable text format (JSON is default)")

	rootCmd.AddCommand(InitCmd())
	rootCmd.AddCommand(SetPhaseCmd())
	rootCmd.AddCommand(GateCmd())
	rootCmd.AddCommand(VerifyCmd())
	rootCmd.AddCommand(StatusCmd())
	rootCmd.AddCommand(UpdateCoverageCmd())
	rootCmd.AddCommand(CostCmd())

	return rootCmd
}
