package main

import (
	"fmt"
	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"sort"
	"strings"
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

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run the full CI pipeline natively in Go",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting Native CI Pipeline...")

		fmt.Println("\n=== 1. Running Lint ===")
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "-c", "scripts/.golangci.yml"); err != nil {
			return fmt.Errorf("lint failed: %w", err)
		}

		fmt.Println("\n=== 2. Running Tests ===")
		if err := runCmd("go", "test", "-race", "-cover", "-count=1", "./..."); err != nil {
			return fmt.Errorf("tests failed: %w", err)
		}

		fmt.Println("\n=== 3. Running Vulncheck ===")
		if err := runCmd("go", "run", "golang.org/x/vuln/cmd/govulncheck", "./..."); err != nil {
			return fmt.Errorf("vulncheck failed: %w", err)
		}

		fmt.Println("\n=== 4. Checking API/CLI Contract Drift ===")
		if err := apiSpecCmd.RunE(cmd, args); err != nil {
			return err
		}

		fmt.Println("\n=== 5. Checking Template Drift ===")
		if err := templateDriftCmd.RunE(cmd, args); err != nil {
			return err
		}

		fmt.Println("\n✅ CI Pipeline Passed Successfully!")
		return nil
	},
}

func getStagedFiles(extension string) ([]string, error) {
	out, err := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR").Output()
	if err != nil {
		return nil, err
	}
	var files []string
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		if strings.HasSuffix(l, extension) {
			files = append(files, l)
		}
	}
	return files, nil
}

var precommitCmd = &cobra.Command{
	Use:   "precommit",
	Short: "Highly optimized, parallelized checks restricted only to staged files",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting Fast Pre-Commit Engine...")

		// 1. Mod Tidy drift check
		fmt.Println("📦 Checking go.mod/go.sum drift...")
		if err := runCmd("go", "mod", "tidy"); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}
		if err := exec.Command("git", "diff", "--exit-code", "go.mod", "go.sum").Run(); err != nil {
			return fmt.Errorf("drift detected in go.mod/go.sum. Please commit changes: %w", err)
		}

		// 2. Fmt staged files
		stagedGoFiles, err := getStagedFiles(".go")
		if err != nil {
			return fmt.Errorf("failed to get staged files: %w", err)
		}
		if len(stagedGoFiles) > 0 {
			fmt.Printf("🧹 Formatting %d staged files...\n", len(stagedGoFiles))
			args := append([]string{"-w"}, stagedGoFiles...)
			if err := runCmd("gofmt", args...); err != nil {
				return fmt.Errorf("gofmt failed: %w", err)
			}
		}

		// 3. Build check
		fmt.Println("🔨 Running fast build syntax check...")
		if err := runCmd("go", "build", "-o", "/dev/null", "./..."); err != nil {
			return fmt.Errorf("build check failed: %w", err)
		}

		// 4. Lint Stage (diff only)
		fmt.Println("🛡️ Running staged-only lint...")
		// Use go run for the pinned tool
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "-c", "scripts/.golangci.yml", "--new-from-rev=HEAD"); err != nil {
			return fmt.Errorf("lint stage failed: %w", err)
		}

		fmt.Println("✅ Pre-commit verification passed!")
		return nil
	},
}

func main() {
	rootCmd.AddCommand(apiSpecCmd, templateDriftCmd, costCmd, coverageCmd, auditCmd, ciCmd, precommitCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
