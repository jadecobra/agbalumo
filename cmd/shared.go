package cmd

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/spf13/cobra"
)

func initRepo() *sqlite.SQLiteRepository {
	dbPath := getDatabaseURL()
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	return repo
}

func getDatabaseURL() string {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}
	return ".tester/data/agbalumo.db"
}

func exitOnErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		os.Exit(1)
	}
}

func printListResponse(cmd *cobra.Command, items any, count int, emptyMsg string) bool {
	if count == 0 {
		if !flagText {
			cmd.Println("[]")
		} else {
			cmd.Println(emptyMsg)
		}
		return true
	}

	if !flagText {
		data, _ := json.MarshalIndent(items, "", "  ")
		cmd.Println(string(data))
		return true
	}

	return false
}
