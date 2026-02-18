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
		dbPath := ResolveSeedConfig(args, os.Getenv)

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

// ResolveSeedConfig determines the database path from arguments or environment variables.
func ResolveSeedConfig(args []string, getEnv func(string) string) string {
	dbPath := "agbalumo.db"
	if len(args) > 0 {
		dbPath = args[0]
	}

	// Also check env var if not argument provided (or default is used)
	if dbPath == "agbalumo.db" {
		if envPath := getEnv("DATABASE_URL"); envPath != "" {
			dbPath = envPath
		}
	}
	return dbPath
}
