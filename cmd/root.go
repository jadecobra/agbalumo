package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agbalumo",
	Short: "Agbalumo is a directory and request platform",
	Long:  `Agbalumo is a high-performance directory and request platform for the West African diaspora.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be defined here
}
