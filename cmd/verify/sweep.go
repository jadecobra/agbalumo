package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var sweepCmd = &cobra.Command{
	Use:   "sweep",
	Short: "Run all structural and meta gates in one cold start",
	Run: func(cmd *cobra.Command, args []string) {
		results, err := maintenance.RunSweep(".")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		useJSON, _ := cmd.Flags().GetBool("json")
		if useJSON {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println("Gate                 Status   Details")
		fmt.Println("─────────────────────────────────────")
		failed := false
		for _, r := range results {
			fmt.Printf("%-20s %-8s %s\n", r.Gate, formatStatus(r.Status), r.Details)
			if r.Status == "FAIL" {
				failed = true
			}
		}

		if failed {
			os.Exit(1)
		}
	},
}

func formatStatus(status string) string {
	switch status {
	case "PASS":
		return "✅ PASS"
	case "FAIL":
		return "❌ FAIL"
	case "WARN":
		return "⚠️ WARN"
	default:
		return status
	}
}

func init() {
	sweepCmd.Flags().Bool("json", false, "Output results in JSON format")
}
