package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

const localCIImageTag = "agbalumo:local-ci-check"

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run the full CI pipeline in parallel with dynamic concurrency",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		tasks := []maintenance.CITask{
			{Name: "Verifying Repository Cleanliness", Fn: func() error { return maintenance.VerifyGitClean(".") }},
			{Name: "Verifying GitHub Action SHAs", Fn: func() error { return maintenance.VerifyActionSHAs(".") }},
			{Name: "Verifying CI Toolset", Fn: func() error { return maintenance.VerifyCITools(".") }},
			{Name: "Verifying JS Syntax", Fn: func() error { return maintenance.VerifyJSSyntax(".") }},
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
			{Name: "Checking UI Design standards", Fn: func() error { return designCmd.RunE(cmd, args) }},
			{Name: "Checking Visual Snapshot Parity", Fn: func() error { return snapshotParityCmd.RunE(cmd, args) }},
			{Name: "Checking Playwright Docker Parity", Fn: func() error { return playwrightVersionCmd.RunE(cmd, args) }},
			{Name: "Checking Coverage Threshold", Fn: func() error { return coverageCmd.RunE(cmd, args) }},
			{Name: "Running Performance Audit (Benchmarks)", Fn: func() error { return perfCmd.RunE(cmd, args) }},
			{Name: "Dynamic Server Startup Audit", Fn: func() error { return maintenance.VerifyServerStartup(".") }},
		}

		withDocker, _ := cmd.Flags().GetBool("with-docker")
		if withDocker {
			tasks = append(tasks, maintenance.CITask{
				Name: "Running Playwright E2E Tests (Linux via Docker)",
				Fn:   func() error { return maintenance.RunPlaywrightInDocker(".") },
			})
		}

		// Run group 1: All checks in parallel (scaled by NumCPU)
		if err := maintenance.RunParallelCI(ctx, tasks); err != nil {
			return err
		}

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
		if err := maintenance.VerifyJSSyntax("."); err != nil {
			return err
		}

		// 1b. Skill checks on staging
		stagedOut, stagedErr := runCmdOutput("git", "diff", "--cached", "--name-only")
		if stagedErr == nil && strings.Contains(string(stagedOut), ".agents/skills/") {

			fmt.Println("🔍 Skill files staged. Running skill conformance and resolvability checks...")
			if e := skillConformanceCmd.RunE(cmd, args); e != nil {
				return e
			}
			if e := checkResolvableCmd.RunE(cmd, args); e != nil {
				return e
			}
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
			// Re-stage after formatting to ensure commit is clean
			fmt.Println("📎 Re-staging formatted files...")
			for _, file := range stagedGoFiles {
				if err := runCmd("git", "add", file); err != nil {
					return fmt.Errorf("failed to re-stage %s: %w", file, err)
				}
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

		// 6. Test only staged packages (fast path)
		if len(stagedGoFiles) == 0 {
			fmt.Println("⏩ No Go files staged. Skipping test phase.")
		} else {
			fmt.Println("🧪 Running tests on staged packages only...")
			dirMap := make(map[string]bool)
			for _, file := range stagedGoFiles {
				dir := filepath.Dir(file)
				dirMap[dir] = true
			}
			for dir := range dirMap {
				pkgPath := "./" + dir
				if dir == "." {
					pkgPath = "."
				}
				fmt.Printf("   Testing %s...\n", pkgPath)
				if err := runCmd("go", "test", "-short", pkgPath); err != nil {
					return fmt.Errorf("tests failed in %s: %w", pkgPath, err)
				}
			}
		}

		// 6b. CSS rebuild if templates changed
		stagedHTMLFiles, _ := getStagedFiles(".html")
		stagedCSSFiles, _ := getStagedFiles(".css")
		if len(stagedHTMLFiles) > 0 || len(stagedCSSFiles) > 0 {
			fmt.Println("🎨 Rebuilding CSS (template/CSS files staged)...")
			outputCSS := "ui/static/css/output.css"
			if err := runCmd("npm", "run", "build:css"); err != nil {
				fmt.Printf("⚠️  'npm run build:css' failed: %v\n", err)
				if _, statErr := os.Stat(outputCSS); statErr == nil {
					fmt.Println("ℹ️  Using existing output.css for commit.")
				} else {
					fmt.Printf("⚠️  Trying 'npx tailwindcss' fallback...\n")
					if err2 := runCmd("npx", "tailwindcss", "-i", "./ui/static/css/input.css", "-o", outputCSS, "--minify"); err2 != nil {
						return fmt.Errorf("CSS build failed and no existing output.css found: %w", err2)
					}
				}
			}
			// Re-stage the rebuilt output.css so the commit includes it
			if err := runCmd("git", "add", outputCSS); err != nil {
				return fmt.Errorf("failed to stage rebuilt CSS: %w", err)
			}
		}

		// 6c. Template Contract Verification (if templates staged)
		if len(stagedHTMLFiles) > 0 {
			fmt.Println("📜 Templates staged. Verifying deterministic contracts...")
			if err := runCmd("go", "test", "./internal/ui/", "-v", "-run", "TestGlobalTemplateCoverage"); err != nil {
				return fmt.Errorf("template contract verification failed: %w", err)
			}
			if err := sandboxParityCmd.RunE(cmd, args); err != nil {
				return err
			}
		}

		// 7. Design gate check (ensures no rounding in admin or hardcoded hex)
		if err := designCmd.RunE(cmd, args); err != nil {
			return err
		}

		// 7b. Snapshot Parity Check (if snapshots staged)
		stagedPNGFiles, _ := getStagedFiles(".png")
		hasSnapshots := false
		for _, f := range stagedPNGFiles {
			if strings.Contains(f, "visual.spec.ts-snapshots") {
				hasSnapshots = true
				break
			}
		}
		if hasSnapshots {
			fmt.Println("📸 Visual snapshots staged. Verifying platform parity...")
			if err := snapshotParityCmd.RunE(cmd, args); err != nil {
				return err
			}
		}

		// 8. Documentation Drift Check (if .md files staged)
		stagedMDFiles, _ := getStagedFiles(".md")
		if len(stagedMDFiles) > 0 {
			fmt.Println("📚 Documentation staged. Verifying path/command integrity...")
			if err := docDriftCmd.RunE(cmd, args); err != nil {
				return err
			}
		}

		// 9. Final Cleanliness Check
		fmt.Println("🔍 Final check: Ensuring no unstaged changes remain...")
		if err := maintenance.VerifyGitClean("."); err != nil {
			return fmt.Errorf("pre-commit failed: working tree is dirty. Please stage or stash all changes: %w", err)
		}

		fmt.Println("✅ Pre-commit verification passed!")
		return nil
	},
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

func runVerifyGatedTask(cmd *cobra.Command) error {
	phase, err := maintenance.InferCurrentPhase(".")
	if err != nil {
		return err
	}
	return maintenance.ExecuteGateChecks(".", phase)
}

// getStagedFiles returns a list of staged files matching the given extension.
// Note: .html matching correctly identifies templates in ui/templates/.
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

// runDockerBuild performs a local docker build to catch Dockerfile regressions before push.
func runDockerBuild() error {
	// Build CSS first. We try npm run build:css first, then fallback to npx tailwindcss.
	// We allow this to fail if the output file already exists, as the host node environment
	// can be non-deterministic during parallel CI execution.
	fmt.Println("🎨 Building CSS for Docker...")
	outputCSS := "./ui/static/css/output.css"
	if err := runCmd("npm", "run", "build:css"); err != nil {
		fmt.Printf("⚠️  'npm run build:css' failed: %v\n", err)
		if _, statErr := os.Stat(outputCSS); statErr == nil {
			fmt.Println("ℹ️  Using existing output.css for Docker build.")
		} else {
			fmt.Printf("⚠️  Trying 'npx tailwindcss' fallback...\n")
			if err2 := runCmd("npx", "tailwindcss", "-i", "./ui/static/css/input.css", "-o", outputCSS, "--minify"); err2 != nil {
				return fmt.Errorf("css build failed and no existing output.css found: %w", err2)
			}
		}
	}

	cmd := exec.Command("docker", "build", "--pull", "-t", localCIImageTag, ".")
	cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runTrivyScan runs Trivy against the locally built image.
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
