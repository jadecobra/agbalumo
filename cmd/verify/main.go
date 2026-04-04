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

var verifyShasCmd = &cobra.Command{
	Use:   "verify-shas",
	Short: "Verify all GitHub Action SHAs are pinned",
	RunE: func(cmd *cobra.Command, args []string) error {
		return maintenance.VerifyActionSHAs(".")
	},
}

var ciToolsCmd = &cobra.Command{
	Use:   "ci-tools",
	Short: "Verify CI toolset availability and OS friendliness",
	RunE: func(cmd *cobra.Command, args []string) error {
		return maintenance.VerifyCITools(".")
	},
}

var gitleaksCmd = &cobra.Command{
	Use:   "gitleaks",
	Short: "Run gitleaks secret scan on staged files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return maintenance.CheckGitleaks(".")
	},
}

var ignoredFilesCmd = &cobra.Command{
	Use:   "ignored-files",
	Short: "Check for ignored files staged for commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		return maintenance.CheckIgnoredFiles(".")
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

		fmt.Println("\n=== 1. Verifying GitHub Action SHAs ===")
		if err := maintenance.VerifyActionSHAs("."); err != nil {
			return err
		}

		fmt.Println("\n=== 2. Verifying CI Toolset ===")
		if err := maintenance.VerifyCITools("."); err != nil {
			return err
		}

		fmt.Println("\n=== 3. Running Lint ===")
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "-c", "scripts/.golangci.yml"); err != nil {
			return fmt.Errorf("lint failed: %w", err)
		}

		fmt.Println("\n=== 4. Running Tests ===")
		if err := runCmd("go", "test", "-race", "-cover", "-count=1", "./..."); err != nil {
			return fmt.Errorf("tests failed: %w", err)
		}

		fmt.Println("\n=== 5. Running Vulncheck ===")
		if err := runCmd("go", "run", "golang.org/x/vuln/cmd/govulncheck", "./..."); err != nil {
			return fmt.Errorf("vulncheck failed: %w", err)
		}

		fmt.Println("\n=== 6. Checking API/CLI Contract Drift ===")
		if err := apiSpecCmd.RunE(cmd, args); err != nil {
			return err
		}

		fmt.Println("\n=== 7. Checking Template Drift ===")
		if err := templateDriftCmd.RunE(cmd, args); err != nil {
			return err
		}

		fmt.Println("\n=== 8. Checking Coverage Threshold ===")
		if err := coverageCmd.RunE(cmd, args); err != nil {
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

		// 1. VCS Checks
		fmt.Println("🔍 Running Stage Isolation Checks...")
		if err := maintenance.CheckIgnoredFiles("."); err != nil {
			return err
		}
		if err := maintenance.CheckGitleaks("."); err != nil {
			return err
		}

		// 2. Mod Tidy drift check
		fmt.Println("📦 Checking go.mod/go.sum drift...")
		if err := runCmd("go", "mod", "tidy"); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}
		if err := exec.Command("git", "diff", "--exit-code", "go.mod", "go.sum").Run(); err != nil {
			return fmt.Errorf("drift detected in go.mod/go.sum. Please commit changes: %w", err)
		}

		// 3. Fmt staged files
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

		// 4. Build check
		fmt.Println("🔨 Running fast build syntax check...")
		if err := runCmd("go", "build", "-o", "/dev/null", "./..."); err != nil {
			return fmt.Errorf("build check failed: %w", err)
		}

		// 5. Lint Stage (diff only)
		fmt.Println("🛡️ Running staged-only lint...")
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint", "run", "-c", "scripts/.golangci.yml", "--new-from-rev=HEAD"); err != nil {
			return fmt.Errorf("lint stage failed: %w", err)
		}

		// 6. Coverage gate check (optional for pre-commit but good for anti-degradation)
		if err := coverageCmd.RunE(cmd, args); err != nil {
			return err
		}

		fmt.Println("✅ Pre-commit verification passed!")
		return nil
	},
}

var critiqueCmd = &cobra.Command{
	Use:   "critique",
	Short: "Run ChiefCritic robustness audit natively",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

		fmt.Println("\n[1/4] Checking Cognitive Complexity (gocognit)...")
		gocognitErr := runCmd("go", "run", "github.com/uudashr/gocognit/cmd/gocognit", "-over", "10", "./cmd", "./internal")
		if gocognitErr != nil {
			fmt.Println("❌ Complexity threshold exceeded!")
		}

		fmt.Println("\n[2/4] Checking Repeated Strings (goconst)...")
		_ = runCmd("go", "run", "github.com/jgautheron/goconst/cmd/goconst", "./cmd/...", "./internal/...")

		fmt.Println("\n[3/4] Checking Struct Alignment (fieldalignment)...")
		_ = runCmd("go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment", "./internal/...", "./cmd/...")

		fmt.Println("\n[4/4] Checking Code Duplication (dupl)...")
		_ = runCmd("go", "run", "github.com/mibk/dupl", "-threshold", "15", "-t", "./cmd", "./internal")

		fmt.Println("\n✅ ChiefCritic Audit Complete!")
		if gocognitErr != nil {
			return fmt.Errorf("robustness audit failed due to cognitive complexity")
		}
		return nil
	},
}

var perfCmd = &cobra.Command{
	Use:   "perf",
	Short: "Run performance audit natively",
	RunE: func(cmd *cobra.Command, args []string) error {
		return maintenance.RunPerformanceAudit(".")
	},
}

func main() {
	rootCmd.AddCommand(
		apiSpecCmd,
		templateDriftCmd,
		costCmd,
		coverageCmd,
		auditCmd,
		ciCmd,
		precommitCmd,
		verifyShasCmd,
		ciToolsCmd,
		gitleaksCmd,
		ignoredFilesCmd,
		critiqueCmd,
		perfCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
