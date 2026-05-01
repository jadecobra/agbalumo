package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

type playwrightReport struct {
	Stats struct {
		Expected   int `json:"expected"`
		Unexpected int `json:"unexpected"`
		Flaky      int `json:"flaky"`
		Skipped    int `json:"skipped"`
	} `json:"stats"`
}

var visualAuditCmd = &cobra.Command{
	Use:   "visual-audit",
	Short: "Run deterministic visual audit (static + Playwright)",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Phase 1: Static Checks
		designViolations, err := maintenance.CheckDesignStandards("ui/templates")
		if err != nil {
			return fmt.Errorf("failed design standards check: %w", err)
		}
		keyViolations, err := maintenance.CheckTemplateKeyGaps("ui/templates")
		if err != nil {
			return fmt.Errorf("failed template key gaps check: %w", err)
		}
		staticViolations := append(designViolations, keyViolations...)

		// Phase 2: Playwright
		// #nosec G204 -- maintenance utility running local tests
		pwCmd := exec.Command("npx", "playwright", "test", "tests/e2e/visual-audit.spec.ts", "--reporter=json")
		pwOut, _ := pwCmd.Output() // playwright returns non-zero on test failure, we want the JSON

		var pwReport playwrightReport
		if len(pwOut) > 0 {
			_ = json.Unmarshal(pwOut, &pwReport)
		}

		if jsonOutput {
			report := map[string]interface{}{
				"static": staticViolations,
				"e2e":    pwReport,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(report)
		}

		// Summary Table
		fmt.Printf("\n%-30s %-10s %-10s\n", "CHECK", "COUNT", "STATUS")
		fmt.Println(strings.Repeat("-", 55))

		printRow := func(name string, count int) {
			status := "✅"
			if count > 0 {
				status = "❌"
			}
			fmt.Printf("%-30s %-10d %-10s\n", name, count, status)
		}

		printRow("Design Standards (Static)", len(designViolations))
		printRow("Template Key Gaps (Static)", len(keyViolations))
		printRow("Playwright E2E Visual", pwReport.Stats.Unexpected)

		total := len(staticViolations) + pwReport.Stats.Unexpected
		fmt.Printf("\nTotal Violations: %d\n", total)

		if total > 0 {
			return fmt.Errorf("visual audit failed with %d issues", total)
		}
		return nil
	},
}

func init() {
	visualAuditCmd.Flags().Bool("json", false, "Output full report as JSON")
}
