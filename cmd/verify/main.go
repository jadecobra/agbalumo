package main

import (
	"fmt"
	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/joho/godotenv"
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

var templateDriftCmd = makeSimpleCmd("template-drift", "Detect undefined template functions in HTML templates", func() error {
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
})

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

var gitleaksCmd = makeSimpleCmd("gitleaks", "Run gitleaks secret scan on staged files", func() error {
	return maintenance.CheckGitleaks(".")
})

var ignoredFilesCmd = makeSimpleCmd("ignored-files", "Check for ignored files staged for commit", func() error {
	return maintenance.CheckIgnoredFiles(".")
})

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Run the scraper job manually to enrich listings",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load .env if present
		_ = godotenv.Load(".env")
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = ".tester/data/agbalumo.db"
		}
		repo, err := sqlite.NewSQLiteRepository(dbURL)
		if err != nil {
			return err
		}
		defer func() { _ = repo.Close() }()

		scraper := service.NewWebsiteScraper()
		job := service.NewScraperJob(repo, scraper)

		fmt.Println("🚀 Starting Manual Enrichment Job...")
		count, err := job.EnrichListings(cmd.Context(), 10)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Success! Enriched %d listings with sensory signals.\n", count)
		return nil
	},
}

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

const localCIImageTag = "agbalumo:local-ci-check"

// runDockerBuild performs a local docker build to catch Dockerfile regressions before push.
// It does not push or tag for any registry — build validation only.
func runDockerBuild() error {
	// Build CSS first (Dockerfile expects ui/styles.css to exist)
	if err := runCmd("npm", "run", "build:css"); err != nil {
		return fmt.Errorf("css build failed (required by Docker): %w", err)
	}

	// Enable BuildKit for cache mounts and use --pull to ensure latest base images
	cmd := exec.Command("docker", "build", "--pull", "-t", localCIImageTag, ".")
	cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runTrivyScan runs Trivy against the locally built image with flags that mirror
// production CI exactly (ci.yml lines 197–200). Fails hard if trivy is not installed.
func runTrivyScan() error {
	if _, err := exec.LookPath("trivy"); err != nil {
		return fmt.Errorf(
			"trivy is not installed — required for --with-docker.\n" +
				"Install: brew install trivy  (mac) or https://aquasecurity.github.io/trivy/latest/getting-started/installation/",
		)
	}
	return runCmd("trivy", "image",
		"--exit-code", "1",
		"--ignore-unfixed",
		"--vuln-type", "os,library",
		"--severity", "CRITICAL,HIGH",
		localCIImageTag,
	)
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run the full CI pipeline in parallel with dynamic concurrency",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tasks := []maintenance.CITask{
			{Name: "Verifying GitHub Action SHAs", Fn: func() error { return maintenance.VerifyActionSHAs(".") }},
			{Name: "Verifying CI Toolset", Fn: func() error { return maintenance.VerifyCITools(".") }},
			{Name: "Running Lint", Fn: func() error {
				return runCmd("go", "run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint", "run")
			}},
			{Name: "Enforcing Struct Alignment", Fn: func() error {
				return runCmd("go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest", "./...")
			}},
			{Name: "Running Vulncheck", Fn: func() error {
				return runCmd("go", "run", "golang.org/x/vuln/cmd/govulncheck", "./...")
			}},
			{Name: "Running Heavy Tests (with -race)", Fn: func() error {
				return runCmd("go", "test", "-race", "-cover", "-count=1", "./...")
			}},
			{Name: "Checking ChiefCritic Robustness", Fn: func() error {
				verbose, _ := cmd.Flags().GetBool("verbose")
				return maintenance.RunChiefCriticAudit(".", maintenance.ChiefCriticOptions{
					Full:    true,
					Verbose: verbose,
				})
			}},
			{Name: "Checking API/CLI Contract Drift", Fn: func() error { return apiSpecCmd.RunE(cmd, args) }},
			{Name: "Checking Template Drift", Fn: func() error { return templateDriftCmd.RunE(cmd, args) }},
			{Name: "Checking Coverage Threshold", Fn: func() error { return coverageCmd.RunE(cmd, args) }},
			{Name: "Running Performance Audit (Benchmarks)", Fn: func() error { return perfCmd.RunE(cmd, args) }},
		}

		// Run group 1: All checks in parallel (scaled by NumCPU)
		if err := maintenance.RunParallelCI(ctx, tasks); err != nil {
			return err
		}

		withDocker, _ := cmd.Flags().GetBool("with-docker")
		if withDocker {
			fmt.Println("\n=== Docker Build & Security Scan ===")
			if err := runDockerBuild(); err != nil {
				return fmt.Errorf("docker build failed: %w", err)
			}
			if err := runTrivyScan(); err != nil {
				return fmt.Errorf("trivy scan failed: %w", err)
			}
		}

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
		opts := maintenance.ChiefCriticOptions{Full: false, NewFromRev: "HEAD", Verbose: false}
		if err := maintenance.RunChiefCriticAudit(".", opts); err != nil {
			return fmt.Errorf("robustness audit failed: %w", err)
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

func runVerifyGatedTask(cmd *cobra.Command) error {
	phase, err := maintenance.InferCurrentPhase(".")
	if err != nil {
		return err
	}
	return maintenance.ExecuteGateChecks(".", phase)
}

var testCmd = &cobra.Command{
	Use:   "test [pkg]",
	Short: "Run tests with race detection and coverage enforcement",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runVerifyGatedTask(cmd); err != nil {
			return err
		}
		pkg := "./..."
		if len(args) > 0 {
			pkg = args[0]
		}
		race, path := getVerificationOpts(cmd)
		short, _ := cmd.Flags().GetBool("short")
		parallel, _ := cmd.Flags().GetInt("parallel")
		return maintenance.RunTests(pkg, race, path, short, parallel)
	},
}

func setupVerifyFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("race", true, "Enable race detection")
	cmd.Flags().String("threshold-path", "", "Path to coverage threshold file")
}

func setupTestFlags(cmd *cobra.Command) {
	setupVerifyFlags(cmd)
	cmd.Flags().Bool("short", false, "Skip slow integration tests (e.g. govulncheck)")
	cmd.Flags().Int("parallel", 0, "Max parallel tests per package (0 = Go default)")
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

var gosecRationaleCmd = makeSimpleCmd("gosec-rationale", "Verify that all #nosec directives include a rationale comment", func() error {
	fmt.Println("🔍 Checking for mandatory rationale in #nosec directives...")
	return maintenance.CheckGosecRationale(".")
})

func init() {
	setupTestFlags(testCmd)
	setupVerifyFlags(coverageCmd)
	setupVerifyFlags(ciCmd)
	setupVerifyFlags(precommitCmd)
	auditCmd.Flags().String("mode", "", "Audit mode: 'static' (no server required) or 'dynamic' (requires live server). Default runs all checks.")
	ciCmd.Flags().Bool("with-docker", false, "Run docker build + trivy image scan (mirrors production CI). Requires Docker and trivy.")
	critiqueCmd.Flags().Bool("full", false, "Run full audit instead of incremental")
	critiqueCmd.Flags().String("baseline", "", "Git revision to compare against (default: HEAD~1)")
	critiqueCmd.Flags().Bool("verbose", false, "Restore full linter logs (disables summarization)")
	ciCmd.Flags().Bool("verbose", false, "Restore full linter logs in summary steps")

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
		healCmd,
		perfCmd,
		checkGatesCmd,
		testCmd,
		watchCmd,
		enrichCmd,
		gosecRationaleCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
