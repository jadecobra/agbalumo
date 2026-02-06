package main

import (
	"context"
	"log"
	"os"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
)

func main() {
	dbPath := "agbalumo.db" // Assumes running from root
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB at %s: %v", dbPath, err)
	}

	ctx := context.Background()

	log.Println("Starting full seed...")
	seeder.SeedAll(ctx, repo)
	log.Println("Full seed complete!")
}
