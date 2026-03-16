package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/spf13/cobra"
)

var stressCount int

var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Generate stress test data for listings",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := ResolveSeedConfig(args, os.Getenv)

		repo, err := sqlite.NewSQLiteRepository(dbPath)
		if err != nil {
			slog.Error("Failed to open DB", "path", dbPath, "error", err)
			os.Exit(1)
		}

		ctx := context.Background()

		slog.Info("Starting stress generation...", "count", stressCount)
		seeder.GenerateStressData(ctx, repo, stressCount)
		slog.Info("Stress generation complete!")
	},
}

func init() {
	stressCmd.Flags().IntVarP(&stressCount, "count", "c", 100000, "Number of listings to generate")
	rootCmd.AddCommand(stressCmd)
}
