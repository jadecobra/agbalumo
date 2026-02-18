package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agbalumo",
	Short: "Agbalumo is a directory and request platform",
	Long:  `Agbalumo is a high-performance directory and request platform for the West African diaspora.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be defined here
}
