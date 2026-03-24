package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/spf13/cobra"
)

func UpdateCoverageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-coverage <package_path> <threshold>",
		Short: "Update the test coverage threshold for a given package",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			pkgPath := args[0]
			thresholdStr := args[1]
			
			threshold, err := strconv.ParseFloat(thresholdStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid threshold '%s'\n", thresholdStr)
				os.Exit(1)
			}

			path := ".agents/coverage-thresholds.json"
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading coverage thresholds: %v\n", err)
				os.Exit(1)
			}

			thresholds, err := agent.ParseThresholds(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing coverage thresholds: %v\n", err)
				os.Exit(1)
			}

			// ONLY INCREASE, NEVER LOWER
			if existing, ok := thresholds[pkgPath]; ok && threshold < existing {
				fmt.Fprintf(os.Stderr, "❌ CRITICAL RULE VIOLATION: NEVER lower the test coverage below its current value (%.1f%%). If test coverage drops due to your changes, you MUST write new tests to cover the new or modified code.\n", existing)
				os.Exit(1)
			}

			thresholds[pkgPath] = threshold

			if err := agent.SaveThresholds(path, thresholds); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving coverage thresholds: %v\n", err)
				os.Exit(1)
			}

			if flagText {
				fmt.Printf("✅ Successfully updated coverage threshold for %s to %.1f%%\n", pkgPath, threshold)
			} else {
				printJSON(true, "update-coverage", map[string]any{
					"package":   pkgPath,
					"threshold": threshold,
				}, nil)
			}
		},
	}
}
