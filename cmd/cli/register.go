package cli

import (
	"github.com/spf13/cobra"
)

// RegisterCommands adds all CLI commands to the root command.
func RegisterCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(AdminCmd)
	rootCmd.AddCommand(CategoryCmd)
	rootCmd.AddCommand(ListingCmd)
	rootCmd.AddCommand(BenchmarkCmd)
	rootCmd.AddCommand(StressCmd)
	rootCmd.AddCommand(SeedCmd)
}
