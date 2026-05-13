package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const defaultTargetURL = "https://localhost:8443"

var hitboxesCmd = &cobra.Command{
	Use:   "hitboxes",
	Short: "Audit touch target hitboxes and interaction layers",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = godotenv.Load(".env")
		targetURL := os.Getenv("BASE_URL")
		if targetURL == "" {
			targetURL = defaultTargetURL
		}

		fmt.Printf("🔍 Auditing hitboxes at %s...\n", targetURL)
		violations, err := maintenance.CheckHitboxes(targetURL)
		if err != nil {
			return fmt.Errorf("hitbox audit failed: %w", err)
		}

		if len(violations) == 0 {
			fmt.Println("✅ No hitbox or interaction layer violations found.")
			return nil
		}

		fmt.Printf("❌ Found %d hitbox violation(s):\n", len(violations))
		for _, v := range violations {
			msg := fmt.Sprintf("- [%s] '%s': %s", v.Tag, v.Text, v.Reason)
			if v.Width > 0 || v.Height > 0 {
				msg += fmt.Sprintf(" (Actual: %.1fx%.1f)", v.Width, v.Height)
			}
			if v.BlockedBy != "" {
				msg += fmt.Sprintf(" (Blocked by: %s)", v.BlockedBy)
			}
			fmt.Println(msg)
		}

		return fmt.Errorf("hitbox audit failed with %d violation(s)", len(violations))
	},
}
