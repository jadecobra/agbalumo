package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
	"github.com/uudashr/gocognit"
)

var complexityCmd = &cobra.Command{
	Use:   "complexity",
	Short: "Early profiling of cognitive complexity across Go functions",
	RunE: func(cmd *cobra.Command, args []string) error {
		threshold, _ := cmd.Flags().GetInt("threshold")
		top, _ := cmd.Flags().GetInt("top")
		diagnostics, _ := cmd.Flags().GetBool("diagnostics")

		opts := maintenance.ComplexityProfileOptions{
			Threshold:   threshold,
			Top:         top,
			Diagnostics: diagnostics,
		}

		stats, err := maintenance.ProfileComplexity(".", opts)
		if err != nil {
			return err
		}

		printComplexityReport(stats, threshold, diagnostics)
		return nil
	},
}

func init() {
	complexityCmd.Flags().IntP("threshold", "t", 10, "Minimum cognitive complexity score to report")
	complexityCmd.Flags().Int("top", 20, "Number of highest complexity functions to display (0 for all)")
	complexityCmd.Flags().Bool("diagnostics", false, "Show detailed breakdown of why complexity increased")
}

func printComplexityReport(stats []gocognit.Stat, threshold int, diagnostics bool) {
	if len(stats) == 0 {
		fmt.Printf("✅ No functions found with cognitive complexity >= %d.\n", threshold)
		return
	}

	fmt.Printf("📊 Cognitive Complexity Report (Threshold: %d, Found: %d):\n", threshold, len(stats))
	for _, stat := range stats {
		fmt.Printf("📍 [%d] %s (%s:%d)\n", stat.Complexity, stat.FuncName, stat.Pos.Filename, stat.Pos.Line)
		if diagnostics {
			printStatDiagnostics(stat.Diagnostics)
		}
	}
}

func printStatDiagnostics(diagnostics []gocognit.Diagnostic) {
	for _, diag := range diagnostics {
		fmt.Printf("    +%-2d (nesting: %d) line %d: %s\n", diag.Inc, diag.Nesting, diag.Pos.Line, diag.Text)
	}
}
