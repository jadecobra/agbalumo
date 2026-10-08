package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Dumps the active SQLite schema deterministically",
	RunE: func(cmd *cobra.Command, args []string) error {
		schema, err := maintenance.DumpSQLiteSchema("listings.db")
		if err != nil {
			return err
		}
		fmt.Println(schema)
		return nil
	},
}

var dbPathCmd = &cobra.Command{
	Use:   "db-path",
	Short: "Prints the canonical path to the populated SQLite snapshot/database",
	RunE: func(cmd *cobra.Command, args []string) error {
		candidates := []string{
			".tester/data/prod_snapshot.db",
			"listings.db",
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				fmt.Println(p)
				return nil
			}
		}
		return fmt.Errorf("no sqlite database found among: %v", candidates)
	},
}
