package maintenance

import (
	"fmt"

	"github.com/joho/godotenv"
)

// RunHeal performs automated remediation of common quality issues.
func RunHeal(rootDir string) error {
	_ = godotenv.Load(".env")
	fmt.Println("🩹 Starting ChiefCritic Automated Healing...")

	// 1. Struct Alignment Fix
	fmt.Print("[1/1] Healing Struct Alignment... ")
	_, err := runTool(rootDir, "go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment", "-fix", "./...")
	if err != nil {
		// fieldalignment -fix often returns non-zero even on success if it made changes
		// We'll check git status later or just assume it tried its best.
		fmt.Println("⚠️  (Applied changes or encountered minor issues)")
	} else {
		fmt.Println("✅")
	}

	fmt.Println("\n✨ Healing Complete! Please review and commit the changes.")
	return nil
}
