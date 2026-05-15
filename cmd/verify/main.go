package main

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "verify",
	Short:         "Agbalumo Maintenance and Verification Utility",
	SilenceUsage:  true,
	SilenceErrors: true, // We handle errors in main
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

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...) //nolint:gosec // maintenance utility runs trusted commands
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output() //nolint:gosec // maintenance utility runs trusted commands
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
		if _, err := os.Stat(".agents/coverage.json"); err == nil {
			path = ".agents/coverage.json"
		} else if _, err := os.Stat(".metrics/coverage"); err == nil {
			path = ".metrics/coverage"
		} else {
			path = ".metrics/coverage"
		}
	}
	return race, path
}

func init() {
	setupTestFlags(testCmd)
	setupVerifyFlags(coverageCmd)
	setupVerifyFlags(ciCmd)
	setupVerifyFlags(precommitCmd)

	auditCmd.Flags().String("mode", "", "Audit mode: 'static' (no server required) or 'dynamic' (requires live server). Default runs all checks.")
	ciCmd.Flags().Bool("with-docker", false, "Run docker build + trivy image scan (mirrors production CI). Requires Docker and trivy.")
	ciCmd.Flags().Bool("verbose", false, "Restore full linter logs in summary steps")
	critiqueCmd.Flags().Bool("full", false, "Run full audit instead of incremental")
	critiqueCmd.Flags().String("baseline", "", "Git revision to compare against (default: HEAD~1)")
	critiqueCmd.Flags().Bool("verbose", false, "Restore full linter logs (disables summarization)")

	rootCmd.AddCommand(
		// CI Domain
		ciCmd,
		precommitCmd,
		testCmd,
		browserCmd,

		// Drift Domain
		docDriftCmd,
		apiSpecCmd,
		templateDriftCmd,
		templateContractCmd,
		deprecatedCmd,

		// Jobs Domain
		locationBackfillCmd,
		enrichCmd,

		// Misc Domain
		costCmd,
		coverageCmd,
		auditCmd,
		verifyShasCmd,
		ciToolsCmd,
		jsSyntaxCmd,
		gitleaksCmd,
		ignoredFilesCmd,
		critiqueCmd,
		healCmd,
		perfCmd,
		checkGatesCmd,
		watchCmd,
		gosecRationaleCmd,
		preflightCmd,
		sessionContextCmd,
		janitorCmd,
		dumpInvariantsCmd,
		designCmd,
		visualAuditCmd,
		skillConformanceCmd,
		checkResolvableCmd,
		mapCmd,
		schemaCmd,
		traceCmd,
		rootHygieneCmd,
		resolveCmd,
		agentsCoverageCmd,
		playwrightConfigCmd,
		testIsolationCmd,
		snapshotParityCmd,
		uptimeCmd,
		playwrightVersionCmd,
		gitCleanCmd,
		minifyContextCmd,
		hitboxesCmd,
		sandboxParityCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if err.Error() != "" {
			_, _ = os.Stderr.WriteString("Error: " + err.Error() + "\n")
		}
		os.Exit(1)
	}
}
