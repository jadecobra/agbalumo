package main

import (
	"fmt"
	"sort"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var designCmd = &cobra.Command{
	Use:   "design",
	Short: "Detect UI Dialect violations (rounding in admin, hardcoded hex codes)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking UI Design Standards...")

		violations, err := maintenance.CheckDesignStandards("ui/templates")
		if err != nil {
			return fmt.Errorf("failed to check design standards: %w", err)
		}

		if len(violations) == 0 {
			fmt.Println("✅ All templates follow UI Dialect standards.")
			return nil
		}

		// Group violations by file for better reporting
		fileViolations := make(map[string][]maintenance.DesignViolation)
		var files []string
		for _, v := range violations {
			if _, ok := fileViolations[v.File]; !ok {
				files = append(files, v.File)
			}
			fileViolations[v.File] = append(fileViolations[v.File], v)
		}
		sort.Strings(files)

		for _, file := range files {
			fmt.Printf("\n📄 %s\n", file)
			for _, v := range fileViolations[file] {
				fmt.Printf("  L%d: %s\n      %s\n", v.Line, v.Reason, v.Content)
			}
		}

		return fmt.Errorf("design violations detected (%d issues)", len(violations))
	},
}
