package commands

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func CostCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cost [dir]",
		Short: "Audit the codebase to measure the agent context cost (RMS of LOC)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			report, err := agent.CalculateContextCost(dir)
			if err != nil {
				return fmt.Errorf("calculating context cost: %w", err)
			}
			if flagText {
				fmt.Printf("--- Context Cost Report ---\n")
				fmt.Printf("Total Files: %d\n", report.TotalFiles)
				fmt.Printf("Total Lines: %d\n", report.TotalLines)
				fmt.Printf("RMS (Root Mean Square) LOC: %.2f\n\n", report.RMS)
				if len(report.TopFiles) > 0 {
					fmt.Printf("Top %d Most Expensive Files:\n", len(report.TopFiles))
					for i, fc := range report.TopFiles {
						fmt.Printf("%2d. [%5d lines] %s\n", i+1, fc.Lines, fc.FilePath)
					}
				}
			} else {
				printJSON(true, "cost", map[string]any{
					"totalFiles": report.TotalFiles,
					"totalLines": report.TotalLines,
					"rms":        report.RMS,
					"topFiles":   report.TopFiles,
				}, nil)
			}

			return nil
		},
	}
}
