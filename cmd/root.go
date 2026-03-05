package cmd

import (
	"github.com/spf13/cobra"
)

var flagJSON bool

var rootCmd = &cobra.Command{
	Use:   "agbalumo",
	Short: "agbalumo is a directory and request platform",
	Long:  `agbalumo is a high-performance directory and request platform for the West African diaspora.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format")
}
