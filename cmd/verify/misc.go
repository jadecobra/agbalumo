package main

import (
	"fmt"
	"os"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var costCmd = makeSimpleCmd("context-cost", "Calculate codebase token density and context window usage", func() error {
	fmt.Println("📊 Calculating Context Cost...")
	report, err := maintenance.CalculateContextCost(".")
	if err != nil {
		return err
	}
	fmt.Printf("Total Files:  %d\n", report.TotalFiles)
	fmt.Printf("Total Lines:  %d\n", report.TotalLines)
	fmt.Printf("Total Tokens: %d\n", report.TotalTokens)
	fmt.Printf("RMS (Lines):  %.2f\n", report.RMS)
	fmt.Printf("Context Usage: %.2f%% of Claude Sonnet window (200k)\n", report.ContextWindowPct)
	return nil
})

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Enforce coverage threshold anti-degradation",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, path := getVerificationOpts(cmd)
		fmt.Printf("🛡️ Checking Coverage Anti-Degradation (%s)...\n", path)
		if err := maintenance.CompareCoverageThreshold(path); err != nil {
			fmt.Printf("❌ %v\n", err)
			return err
		}
		fmt.Println("✅ Coverage threshold check passed.")
		return nil
	},
}

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run comprehensive security and health audit",
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, _ := cmd.Flags().GetString("mode")
		cfg := maintenance.AuditConfig{
			TargetURL: os.Getenv("APP_URL"),
			RootDir:   ".",
			Mode:      mode,
		}
		if cfg.TargetURL == "" {
			cfg.TargetURL = "https://localhost:8443"
		}
		return maintenance.RunSecurityAudit(cfg)
	},
}

var verifyShasCmd = makeSimpleCmd("verify-shas", "Verify all GitHub Action SHAs are pinned", func() error {
	return maintenance.VerifyActionSHAs(".")
})

var ciToolsCmd = makeSimpleCmd("ci-tools", "Verify CI toolset availability and OS friendliness", func() error {
	return maintenance.VerifyCITools(".")
})

var jsSyntaxCmd = makeSimpleCmd("js-syntax", "Verify JavaScript syntax using node -c", func() error {
	return maintenance.VerifyJSSyntax(".")
})

var gitleaksCmd = makeSimpleCmd("gitleaks", "Run gitleaks secret scan on staged files", func() error {
	return maintenance.CheckGitleaks(".")
})

var ignoredFilesCmd = makeSimpleCmd("ignored-files", "Check for ignored files staged for commit", func() error {
	return maintenance.CheckIgnoredFiles(".")
})

var critiqueCmd = &cobra.Command{
	Use:   "critique",
	Short: "Run ChiefCritic robustness audit natively",
	RunE: func(cmd *cobra.Command, args []string) error {
		full, _ := cmd.Flags().GetBool("full")
		rev, _ := cmd.Flags().GetString("baseline")
		verbose, _ := cmd.Flags().GetBool("verbose")
		return maintenance.RunChiefCriticAudit(".", maintenance.ChiefCriticOptions{
			Full:       full,
			NewFromRev: rev,
			Verbose:    verbose,
		})
	},
}

var healCmd = makeSimpleCmd("heal", "Perform automated remediation of quality violations", func() error {
	return maintenance.RunHeal(".")
})

var perfCmd = makeSimpleCmd("perf", "Run performance audit natively", func() error {
	return maintenance.RunPerformanceAudit(".")
})

var checkGatesCmd = &cobra.Command{
	Use:   "check-gates",
	Short: "Verify TDD workflow gates based on Git history and staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVerifyGatedTask(cmd)
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch [command] [args...]",
	Short: "Watch files and restart a command (e.g., serve or test)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdName := "go"
		cmdArgs := []string{"run", "main.go", "serve"}
		if len(args) > 0 {
			cmdName = args[0]
			cmdArgs = args[1:]
		}
		return maintenance.Watch(cmd.Context(), cmdName, cmdArgs)
	},
}

var gosecRationaleCmd = makeSimpleCmd("gosec-rationale", "Verify that all #nosec directives include a rationale comment", func() error {
	fmt.Println("🔍 Checking for mandatory rationale in #nosec directives...")
	return maintenance.CheckGosecRationale(".")
})
