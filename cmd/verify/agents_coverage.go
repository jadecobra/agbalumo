package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var agentsCoverageCmd = &cobra.Command{
	Use:   "agents-coverage",
	Short: "Check for missing AGENTS.md in Go packages",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		missing, total, err := maintenance.CheckAgentsCoverage(cwd)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		if len(missing) > 0 {
			fmt.Println("❌ Missing AGENTS.md:")
			for _, p := range missing {
				fmt.Printf("→   %s/\n", p)
			}
			covered := total - len(missing)
			percent := float64(covered) / float64(total) * 100
			fmt.Printf("📊 Coverage: %d/%d packages (%.1f%%)\n", covered, total, percent)
			os.Exit(1)
		}

		fmt.Printf("✅ 100%% AGENTS.md coverage (%d/%d packages).\n", total, total)
	},
}
