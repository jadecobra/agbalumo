package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

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
		start := time.Now()

		listings := seeder.GenerateStressListings(stressCount)
		
		slog.Info("Inserting listings...", "count", len(listings))
		if err := repo.BulkInsertListings(ctx, listings); err != nil {
			slog.Error("Failed to bulk insert listings", "error", err)
			os.Exit(1)
		}

		duration := time.Since(start)
		fmt.Printf("Stress generation and insertion of %d listings complete in %v\n", stressCount, duration)
	},
}

func init() {
	stressCmd.Flags().IntVarP(&stressCount, "count", "c", 10000, "Number of listings to generate")
	rootCmd.AddCommand(stressCmd)
}
