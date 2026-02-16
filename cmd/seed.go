package cmd

import (
	"context"
	"log"
	"os"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := "agbalumo.db" // Assumes running from root
		if len(args) > 0 {
			dbPath = args[0]
		}
		
		// Also check env var if not argument provided
		if dbPath == "agbalumo.db" && os.Getenv("DATABASE_URL") != "" {
			dbPath = os.Getenv("DATABASE_URL")
		}

		repo, err := sqlite.NewSQLiteRepository(dbPath)
		if err != nil {
			log.Fatalf("Failed to open DB at %s: %v", dbPath, err)
		}

		ctx := context.Background()

		log.Println("Starting full seed...")
		seeder.SeedAll(ctx, repo)
		log.Println("Full seed complete!")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
