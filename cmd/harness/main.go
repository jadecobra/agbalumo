package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command for the harness CLI.
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "harness",
		Short: "agentic-harness is a robust CLI for 10x Engineer workflows",
		Long:  `agentic-harness provides utilities for state machine transitions, gating, and environment verification.`,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the harness tracking structure",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("init called")
		},
	}

	var setPhaseCmd = &cobra.Command{
		Use:   "set-phase",
		Short: "Set the current workflow phase",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("set-phase called")
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Print the current status of the harness",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("status called")
		},
	}

	var gateCmd = &cobra.Command{
		Use:   "gate",
		Short: "Run the validation gate for the current phase",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gate called")
		},
	}

	rootCmd.AddCommand(initCmd, setPhaseCmd, statusCmd, gateCmd)

	return rootCmd
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
