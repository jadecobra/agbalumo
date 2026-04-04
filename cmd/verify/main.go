package main

import (
	"fmt"
	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"sort"
)

var rootCmd = &cobra.Command{
	Use:   "verify",
	Short: "Agbalumo Maintenance and Verification Utility",
}

var apiSpecCmd = &cobra.Command{
	Use:   "api-spec",
	Short: "Detect drift between Code, OpenAPI, and Markdown docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking API and CLI Drift...")

		// 1. Code Routes
		codeRoutes, err := maintenance.ExtractRoutes("cmd", "internal/handler", "internal/module", "internal/infra")
		if err != nil {
			return fmt.Errorf("failed to extract code routes: %w", err)
		}

		// 2. OpenAPI Routes
		bundleCmd := exec.Command("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
		openapiData, err := bundleCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to bundle openapi.yaml (ensure npx/swagger-cli is available): %w", err)
		}
		openapiRoutes, err := maintenance.ExtractOpenAPIRoutes(openapiData)
		if err != nil {
			return fmt.Errorf("failed to extract openapi routes: %w", err)
		}

		// 3. Markdown Routes
		mdData, err := os.ReadFile("docs/api.md")
		if err != nil {
			return fmt.Errorf("failed to read docs/api.md: %w", err)
		}
		mdRoutes, err := maintenance.ExtractMarkdownRoutes(mdData)
		if err != nil {
			return fmt.Errorf("failed to extract markdown routes: %w", err)
		}

		// 4. Compare
		var allDrifts []string
		allDrifts = append(allDrifts, maintenance.CompareRoutes("Code", "OpenAPI", codeRoutes, openapiRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("OpenAPI", "Code", openapiRoutes, codeRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("Code", "API Docs", codeRoutes, mdRoutes)...)
		allDrifts = append(allDrifts, maintenance.CompareRoutes("API Docs", "Code", mdRoutes, codeRoutes)...)

		// 5. CLI Drift
		cliCodeCmds, err := maintenance.ExtractCLICodeCommands("cmd")
		if err == nil {
			cliMDCmds, err := maintenance.ExtractCLIMarkdownCommands("docs/cli.md", "docs/cli")
			if err == nil {
				codeMap := make(map[string]bool)
				for _, c := range cliCodeCmds {
					codeMap[c] = true
				}
				mdMap := make(map[string]bool)
				for _, c := range cliMDCmds {
					mdMap[c] = true
				}
				for _, c := range cliCodeCmds {
					if !mdMap[c] {
						allDrifts = append(allDrifts, fmt.Sprintf("Missing in CLI Docs: %s (found in Code)", c))
					}
				}
				for _, c := range cliMDCmds {
					if !codeMap[c] {
						allDrifts = append(allDrifts, fmt.Sprintf("Missing in Code: %s (found in CLI Docs)", c))
					}
				}
			}
		}

		if len(allDrifts) > 0 {
			sort.Strings(allDrifts)
			for _, d := range allDrifts {
				fmt.Printf("❌ %s\n", d)
			}
			return fmt.Errorf("contract drift detected (%d issues)", len(allDrifts))
		}

		fmt.Println("✅ API and CLI contracts are in sync.")
		return nil
	},
}

var templateDriftCmd = &cobra.Command{
	Use:   "template-drift",
	Short: "Detect undefined template functions in HTML templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking Template Function Drift...")

		defined, err := maintenance.ExtractRendererFunctions("internal/ui/renderer.go")
		if err != nil {
			return err
		}
		used, err := maintenance.ExtractTemplateFunctionCalls("ui/templates")
		if err != nil {
			return err
		}

		drifts := maintenance.CheckTemplateDrift(defined, used)
		if len(drifts) > 0 {
			for _, d := range drifts {
				fmt.Printf("❌ %s\n", d)
			}
			return fmt.Errorf("template function drift detected (%d issues)", len(drifts))
		}

		fmt.Println("✅ All template functions are defined.")
		return nil
	},
}

var costCmd = &cobra.Command{
	Use:   "context-cost",
	Short: "Calculate codebase token density and context window usage",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Enforce coverage threshold anti-degradation",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ".metrics/coverage"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Check legacy path if new one doesn't exist
			if _, e := os.Stat(".agents/coverage-threshold"); e == nil {
				path = ".agents/coverage-threshold"
			}
		}
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
		cfg := maintenance.AuditConfig{
			TargetURL: os.Getenv("APP_URL"),
			RootDir:   ".",
		}
		if cfg.TargetURL == "" {
			cfg.TargetURL = "https://localhost:8443"
		}
		return maintenance.RunSecurityAudit(cfg)
	},
}

func main() {
	rootCmd.AddCommand(apiSpecCmd, templateDriftCmd, costCmd, coverageCmd, auditCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
