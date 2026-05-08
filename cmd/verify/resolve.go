package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve [intent]",
	Short: "Resolve an intent to a skill or command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		intent := args[0]
		rootDir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Find root dir by looking for .git
		for {
			if _, errStat := os.Stat(filepath.Join(rootDir, ".git")); errStat == nil {
				break
			}
			parent := filepath.Dir(rootDir)
			if parent == rootDir {
				break
			}
			rootDir = parent
		}

		matches, err := maintenance.ResolveIntent(rootDir, intent)
		if err != nil {
			return err
		}

		if len(matches) == 0 {
			fmt.Printf("No matches found for intent: %q\n", intent)
			return nil
		}

		for _, m := range matches {
			if m.Type == "skill" {
				fmt.Printf("→ 📖 Skill: %s → %s\n", m.Name, m.Path)
			} else {
				desc := m.Description
				if desc == "" {
					fmt.Printf("→ 🛠️  Command: %s\n", m.Name)
				} else {
					fmt.Printf("→ 🛠️  Command: %s (%s)\n", m.Name, desc)
				}
			}
		}

		return nil
	},
}
