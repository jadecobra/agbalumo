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
		RunE: func(cmd *cobra.Command, args []string) error {
			pkgPath := args[0]
			thresholdStr := args[1]

			threshold, err := strconv.ParseFloat(thresholdStr, 64)
			if err != nil {
				return fmt.Errorf("invalid threshold '%s'", thresholdStr)
			}

			const thresholdPath = ".agents/coverage-thresholds.json"
			// #nosec G304 - File Inclusion: Reading internal configuration from the fixed location.
			data, err := os.ReadFile(thresholdPath)
			if err != nil {
				return fmt.Errorf("reading coverage thresholds: %w", err)
			}

			thresholds, err := agent.ParseThresholds(data)
			if err != nil {
				return fmt.Errorf("parsing coverage thresholds: %w", err)
			}

			// ONLY INCREASE, NEVER LOWER
			if existing, ok := thresholds[pkgPath]; ok && threshold < existing {
				return fmt.Errorf("CRITICAL RULE VIOLATION: NEVER lower the test coverage below its current value (%.1f%%). If test coverage drops due to your changes, you MUST write new tests to cover the new or modified code", existing)
			}

			thresholds[pkgPath] = threshold

			if err := agent.SaveThresholds(thresholdPath, thresholds); err != nil {
				return fmt.Errorf("saving coverage thresholds: %w", err)
			}

			if flagText {
				fmt.Printf("✅ Successfully updated coverage threshold for %s to %.1f%%\n", pkgPath, threshold)
			} else {
				printJSON(true, "update-coverage", map[string]any{
					"package":   pkgPath,
					"threshold": threshold,
				}, nil)
			}

			return nil
		},
	}
}
