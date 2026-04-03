package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CheckVerifiedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-verified",
		Short: "Checks if all core automated gates have passed (exit 0 if verified, 1 otherwise)",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := getState()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: Could not load state: %v\n", err)
				os.Exit(1)
			}

			if state.IsVerified() {
				if flagText {
					fmt.Println("✅ State is VERIFIED. All core automated gates have passed.")
				}
				os.Exit(0)
			} else {
				if flagText {
					fmt.Println("❌ State is NOT VERIFIED. Some core automated gates are still PENDING or FAILED.")
					fmt.Printf("Status: %+v\n", state.Gates)
				}
				os.Exit(1)
			}
		},
	}
}
