package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var lessonsConformanceCmd = &cobra.Command{
	Use:   "lessons-conformance",
	Short: "Audit active prose strict lessons for TRIGGERS and the 20-lesson ceiling",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking Strict Lessons Conformance...")

		violations, err := maintenance.CheckLessonsConformance(".")
		if err != nil {
			return fmt.Errorf("lessons conformance check failed: %w", err)
		}

		if len(violations) == 0 {
			fmt.Println("✅ All strict lessons follow trigger and ceiling standards.")
			return nil
		}

		fmt.Printf("❌ Found %d strict lessons conformance violation(s):\n", len(violations))
		for _, v := range violations {
			fmt.Printf("  %s:%d -> %s\n", v.File, v.Line, v.Message)
		}

		os.Exit(1)
		return nil
	},
}
