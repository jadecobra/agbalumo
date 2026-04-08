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

func makeSimpleCmd(use, short string, fn func() error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fn()
		},
	}
}

func toSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

func diffSets(items []string, against map[string]bool, format string) []string {
	var drifts []string
	for _, item := range items {
		if !against[item] {
			drifts = append(drifts, fmt.Sprintf(format, item))
		}
	}
	return drifts
}

const errFmt = "❌ %s\n"

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
		openapiData, err := runCmdOutput("npx", "-y", "swagger-cli", "bundle", "docs/openapi.yaml", "-r", "-t", "yaml")
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
				codeMap := toSet(cliCodeCmds)
				mdMap := toSet(cliMDCmds)

				allDrifts = append(allDrifts, diffSets(cliCodeCmds, mdMap, "Missing in CLI Docs: %s (found in Code)")...)
				allDrifts = append(allDrifts, diffSets(cliMDCmds, codeMap, "Missing in Code: %s (found in CLI Docs)")...)
			}
		}

		return reportDrift("Contract", allDrifts, "✅ API and CLI contracts are in sync.")
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
		return reportDrift("Template function", drifts, "✅ All template functions are defined.")
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

var gitleaksCmd = makeSimpleCmd("gitleaks", "Run gitleaks secret scan on staged files", func() error {
	return maintenance.CheckGitleaks(".")
})

var ignoredFilesCmd = makeSimpleCmd("ignored-files", "Check for ignored files staged for commit", func() error {
	return maintenance.CheckIgnoredFiles(".")
})

func reportDrift(label string, drifts []string, successMsg string) error {
	if len(drifts) == 0 {
		if successMsg != "" {
			fmt.Println(successMsg)
		}
		return nil
	}
	sort.Strings(drifts)
	for _, d := range drifts {
		fmt.Printf(errFmt, d)
	}
	return fmt.Errorf("%s drift detected (%d issues)", label, len(drifts))
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...) //nolint:gosec // maintenance utility runs trusted commands
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output() //nolint:gosec // maintenance utility runs trusted commands
}

// runDockerBuild performs a local docker build to catch Dockerfile regressions before push.
// It does not push or tag for any registry — build validation only.
func runDockerBuild() error {
	// Build CSS first (Dockerfile expects ui/styles.css to exist)
	if err := runCmd("npm", "run", "build:css"); err != nil {
		return fmt.Errorf("css build failed (required by Docker): %w", err)
	}
	return runCmd("docker", "build", "--no-cache", "-t", "agbalumo:local-ci-check", ".")
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
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint", "run"); err != nil {
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

		withDocker, _ := cmd.Flags().GetBool("with-docker")
		if withDocker {
			fmt.Println("\n=== 9. Docker Build Validation ===")
			if err := runDockerBuild(); err != nil {
				return fmt.Errorf("docker build failed: %w", err)
			}
		}

		fmt.Println("\n✅ CI Pipeline Passed Successfully!")
		return nil
	},
}

func getStagedFiles(extension string) ([]string, error) {
	out, err := runCmdOutput("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR")
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
		if err := runCmd("git", "diff", "--exit-code", "go.mod", "go.sum"); err != nil {
			return fmt.Errorf("drift detected in go.mod/go.sum. Please commit changes: %w", err)
		}

		// 3. Fmt staged files
		stagedGoFiles, err := getStagedFiles(".go")
		if err != nil {
			return fmt.Errorf("failed to get staged files: %w", err)
		}
		if len(stagedGoFiles) > 0 {
			fmt.Printf("🧹 Formatting %d staged files...\n", len(stagedGoFiles))
			fmtArgs := append([]string{"-w"}, stagedGoFiles...)
			if err := runCmd("gofmt", fmtArgs...); err != nil {
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
		if err := runCmd("go", "run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint", "run", "--new-from-rev=HEAD"); err != nil {
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

var critiqueCmd = makeSimpleCmd("critique", "Run ChiefCritic robustness audit natively", func() error {
	return maintenance.RunChiefCriticAudit(".")
})

var perfCmd = makeSimpleCmd("perf", "Run performance audit natively", func() error {
	return maintenance.RunPerformanceAudit(".")
})

var checkGatesCmd = &cobra.Command{
	Use:   "check-gates",
	Short: "Verify TDD workflow gates based on Git history and staged changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		phase, err := maintenance.InferCurrentPhase(".")
		if err != nil {
			return err
		}
		return maintenance.ExecuteGateChecks(".", phase)
	},
}

var testCmd = &cobra.Command{
	Use:   "test [pkg]",
	Short: "Run tests with race detection and coverage enforcement",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := "./..."
		if len(args) > 0 {
			pkg = args[0]
		}
		race, path := getVerificationOpts(cmd)
		return maintenance.RunTests(pkg, race, path)
	},
}

func setupVerifyFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("race", true, "Enable race detection")
	cmd.Flags().String("threshold-path", "", "Path to coverage threshold file")
}

func getVerificationOpts(cmd *cobra.Command) (bool, string) {
	race, _ := cmd.Flags().GetBool("race")
	path, _ := cmd.Flags().GetString("threshold-path")
	if path == "" {
		// Check .metrics/coverage then .agents/coverage-threshold
		if _, err := os.Stat(".metrics/coverage"); err == nil {
			path = ".metrics/coverage"
		} else if _, err := os.Stat(".agents/coverage-threshold"); err == nil {
			path = ".agents/coverage-threshold"
		} else {
			path = ".metrics/coverage"
		}
	}
	return race, path
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

var gosecRationaleCmd = &cobra.Command{
	Use:   "gosec-rationale",
	Short: "Verify that all #nosec directives include a rationale comment",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Checking for mandatory rationale in #nosec directives...")
		return maintenance.CheckGosecRationale(".")
	},
}

func init() {
	setupVerifyFlags(testCmd)
	setupVerifyFlags(coverageCmd)
	setupVerifyFlags(ciCmd)
	setupVerifyFlags(precommitCmd)
	auditCmd.Flags().String("mode", "", "Audit mode: 'static' (no server required) or 'dynamic' (requires live server). Default runs all checks.")
	ciCmd.Flags().Bool("with-docker", false, "Also run docker build as a final validation step (requires Docker)")

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
		checkGatesCmd,
		testCmd,
		watchCmd,
		gosecRationaleCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
