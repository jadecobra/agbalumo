package main

import (
	"fmt"

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
