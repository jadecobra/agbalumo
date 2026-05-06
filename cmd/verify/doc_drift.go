package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var docDriftCmd = &cobra.Command{
	Use:   "doc-drift",
	Short: "Detect stale file path references in documentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		violations, err := maintenance.CheckDocDrift(".")
		if err != nil {
			return fmt.Errorf("failed to check doc drift: %w", err)
		}

		if len(violations) == 0 {
			fmt.Println("✅ No documentation drift detected.")
			return nil
		}

		fmt.Printf("❌ Detected %d stale documentation references:\n", len(violations))
		for _, v := range violations {
			fmt.Printf("  %s:%d -> %s (Not Found)\n", v.DocFile, v.Line, v.ReferencedPath)
		}

		os.Exit(1)
		return nil
	},
}
