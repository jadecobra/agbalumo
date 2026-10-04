package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var handlerConcurrencyCmd = &cobra.Command{
	Use:   "handler-concurrency",
	Short: "Audit HTTP handlers in internal/module/ for unbounded or excessive database concurrency",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		opts := maintenance.HandlerConcurrencyOptions{
			MaxGoroutinesPerFunc: limit,
		}

		fmt.Println("🔍 Checking HTTP request handlers for unbounded or excessive concurrency...")
		violations, err := maintenance.CheckHandlerConcurrency(".", opts)
		if err != nil {
			return err
		}

		if len(violations) > 0 {
			fmt.Printf("❌ Found %d handler concurrency violations:\n", len(violations))
			for _, v := range violations {
				fmt.Printf("  - %s:%d [%s]: %s\n", v.File, v.Line, v.Func, v.Message)
			}
			return fmt.Errorf("handler concurrency check failed")
		}

		fmt.Println("✅ All HTTP handlers satisfy concurrency bounds.")
		return nil
	},
}

func init() {
	handlerConcurrencyCmd.Flags().Int("limit", 4, "Maximum allowed concurrent goroutines per handler function")
}
