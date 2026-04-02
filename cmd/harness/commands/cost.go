package commands

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func CostCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cost [dir]",
		Short: "Audit the codebase: LOC RMS gate (active) + token-based context window % vs Claude Sonnet 200k (observation)",
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
				fmt.Printf("Total Files:          %d\n", report.TotalFiles)
				fmt.Printf("Total Lines:          %d\n", report.TotalLines)
				fmt.Printf("Total Tokens:         %d\n", report.TotalTokens)
				fmt.Printf("\n")
				fmt.Printf("LOC  RMS:             %.2f  (gate: ≤110 — active)\n", report.RMS)
				fmt.Printf("Token RMS:            %.2f  (observation only — no gate yet)\n", report.TokenRMS)
				fmt.Printf("Context Window Used:  %.2f%%  of Claude Sonnet 200k (worst case)\n", report.ContextWindowPct)
				fmt.Printf("                      Gemini Flash 1M satisfied automatically if Sonnet passes\n")
				fmt.Printf("\n")
				if len(report.TopFiles) > 0 {
					fmt.Printf("Top %d Most Expensive Files (by token count):\n", len(report.TopFiles))
					for i, fc := range report.TopFiles {
						fmt.Printf("%2d. [%6d tokens | %4d lines] %s\n", i+1, fc.Tokens, fc.Lines, fc.FilePath)
					}
				}
			} else {
				printJSON(true, "cost", map[string]any{
					"totalFiles":       report.TotalFiles,
					"totalLines":       report.TotalLines,
					"totalTokens":      report.TotalTokens,
					"rms":              report.RMS,
					"tokenRms":         report.TokenRMS,
					"contextWindowPct": report.ContextWindowPct,
					"topFiles":         report.TopFiles,
				}, nil)
			}

			return nil
		},
	}
}
