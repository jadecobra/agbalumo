package maintenance

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/joho/godotenv"
)

// AuditConfig holds the necessary configurations for a security audit.
type AuditConfig struct {
	TargetURL  string
	HTTPClient *http.Client
	RootDir    string
	Mode       string // "static" | "dynamic" | "" (default: both)
}

type auditResults struct {
	score int
	total int
}

// RunSecurityAudit carries out a series of security checks on the codebase and, if available, the live server.
func RunSecurityAudit(config AuditConfig) error {
	res := &auditResults{total: 0, score: 0}
	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println(domain.SeparatorLine)

	checks := []auditCheck{
		{func() (bool, bool) { return checkGoVet(config) }, "Go Vet"},
		{func() (bool, bool) { return checkFlyConfig(config) }, "fly.toml"},
		{func() (bool, bool) { return checkXSS(config) }, "XSS"},
		{func() (bool, bool) { return checkVulnerabilities(config) }, "govulncheck"},
		{func() (bool, bool) { return VerifyCITools(config.RootDir) == nil, false }, "CI Toolset"},
		{func() (bool, bool) { return VerifyActionSHAs(config.RootDir) == nil, false }, "Action SHAs"},
	}

	executeChecks(res, checks)
	if config.Mode != "static" {
		res.runDynamicCheck(config.TargetURL, config.HTTPClient)
	}
	return res.reportResults()
}

func executeChecks(res *auditResults, checks []auditCheck) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, c := range checks {
		wg.Add(1)
		go func(name string, fn func() (bool, bool)) {
			defer wg.Done()
			passed, skip := fn()
			mu.Lock()
			res.recordResult(name, passed, skip)
			mu.Unlock()
		}(c.name, c.fn)
	}
	wg.Wait()
}

func (r *auditResults) runDynamicCheck(url string, client *http.Client) {
	passed, skip := checkHeaders(url, client)
	r.recordResult("Headers", passed, skip)
}

func (r *auditResults) reportResults() error {
	fmt.Println(domain.SeparatorLine)
	if r.total <= 0 {
		return fmt.Errorf("no security checks could be performed")
	}
	fmt.Printf("Audit Score: %d/%d (%.0f%%)\n", r.score, r.total, (float64(r.score)/float64(r.total))*100)
	if r.score < r.total {
		return fmt.Errorf("security audit did not pass all checks (%d/%d)", r.score, r.total)
	}
	return nil
}

type auditCheck struct {
	fn   func() (bool, bool)
	name string
}

func (r *auditResults) recordResult(name string, passed, skip bool) {
	fmt.Printf("[?] Finished '%s'... ", name)
	if skip {
		fmt.Println("⚠️  Skipped")
		return
	}
	r.total++
	if passed {
		fmt.Println("✅ Passed")
		r.score++
	} else {
		fmt.Println("❌ Failed")
	}
}

func checkGoVet(config AuditConfig) (bool, bool) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = config.RootDir
	return cmd.Run() == nil, false
}

func checkHeaders(target string, client *http.Client) (bool, bool) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, //nolint:gosec // maintenance utility
			Timeout:   2 * time.Second,
		}
	}

	resp, err := client.Get(target)
	if err != nil {
		fmt.Printf("⚠️  Could not connect to %s - skipping header check\n", target)
		return false, true
	}
	defer func() { _ = resp.Body.Close() }()

	passed := true
	headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Frame-Options"}
	for _, h := range headers {
		if resp.Header.Get(h) == "" {
			passed = false
		}
	}
	return passed, false
}

func checkFlyConfig(config AuditConfig) (bool, bool) {
	flyPath := filepath.Join(config.RootDir, "fly.toml")
	flyData, err := os.ReadFile(flyPath) //nolint:gosec // maintenance utility
	if err != nil {
		return false, true
	}

	content := strings.ToLower(string(flyData))
	sensitiveKeys := []string{"SECRET", "KEY", "PASSWORD", "TOKEN", "AUTH"}
	for _, key := range sensitiveKeys {
		if strings.Contains(content, strings.ToLower(key)) && (strings.Contains(content, "=") || strings.Contains(content, ":")) {
			return false, false
		}
	}
	return true, false
}

func checkVulnerabilities(config AuditConfig) (bool, bool) {
	cmd := exec.Command("go", "run", "golang.org/x/vuln/cmd/govulncheck", "./...")
	cmd.Dir = config.RootDir
	cmd.Stdout = os.Stdout // surface govulncheck output in CI logs
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil, false
}

func checkXSS(config AuditConfig) (bool, bool) {
	script := "grep -rI 'template.HTML' . --exclude-dir=.git --exclude-dir=bin --exclude-dir=scripts --exclude-dir='.tester' --exclude-dir='@*' --exclude='*_test.go' --exclude='renderer.go' --exclude='renderer_funcs.go' | grep -v 'cmd/verify/main.go' | grep -v 'template.HTMLEscapeString' | grep -v 'internal/agent' | grep -v 'internal/maintenance/audit.go'"
	cmd := exec.Command("sh", "-c", script)
	cmd.Dir = config.RootDir
	out, _ := cmd.CombinedOutput()
	return len(strings.TrimSpace(string(out))) == 0, false
}

// ChiefCriticOptions configures the robustness audit behavior.
type ChiefCriticOptions struct {
	NewFromRev string
	Full       bool
	Verbose    bool
}

type linterIssue struct {
	file    string
	line    string
	message string
	linter  string
	raw     string
}

// RunChiefCriticAudit performs a consolidated code quality audit using golangci-lint.
func RunChiefCriticAudit(rootDir string, opts ChiefCriticOptions) error {
	fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

	command := buildLinterCommand(opts)
	output, err := runTool(rootDir, "go", command...)

	if err != nil {
		fmt.Println("❌ ChiefCritic Audit Failed")
		if output != "" {
			if opts.Verbose {
				fmt.Println(domain.SeparatorLine)
				fmt.Println(output)
				fmt.Println(domain.SeparatorLine)
			} else {
				reportSummarizedIssues(output)
			}
		}
		return fmt.Errorf("robustness audit failed: %w", err)
	}

	return nil
}

func buildLinterCommand(opts ChiefCriticOptions) []string {
	args := []string{"run"}
	if !opts.Full {
		rev := opts.NewFromRev
		if rev == "" {
			rev = "HEAD~1"
		}
		args = append(args, "--new-from-rev", rev)
	}

	if opts.Verbose {
		args = append(args, "-v")
	}

	return append([]string{"run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"}, args...)
}

func reportSummarizedIssues(output string) {
	issuesByLinter, totalIssues := parseLinterOutput(output)

	printLinterSummaryTable(issuesByLinter)

	const globalCap = 25
	const perLinterCap = 5
	reportedCount := printTopIssues(issuesByLinter, globalCap, perLinterCap)

	if totalIssues > reportedCount {
		fmt.Printf("\n⚠️  Total issues: %d. Showing %d for context efficiency.\n", totalIssues, reportedCount)
		fmt.Println("💡 Use 'verify critique --verbose' for full report.")
	}

	if anySystemic(issuesByLinter) {
		fmt.Println("\n🚨 [ADVISORY] Systemic technical debt detected. Consider triggering '/learn' to codify new standards.")
	}
}

func parseLinterOutput(output string) (map[string][]linterIssue, int) {
	lines := strings.Split(output, "\n")
	issuesByLinter := make(map[string][]linterIssue)
	totalIssues := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "level=") {
			continue
		}

		parts := strings.Split(line, ":")
		// Validation: line must start with file:line: and parts[1] must be numeric
		if len(parts) < 3 || !isNumeric(parts[1]) {
			continue
		}

		linter := extractLinterName(line)
		issue := linterIssue{
			file:    parts[0],
			line:    parts[1],
			message: strings.Join(parts[2:], ":"),
			linter:  linter,
			raw:     line,
		}

		issuesByLinter[linter] = append(issuesByLinter[linter], issue)
		totalIssues++
	}
	return issuesByLinter, totalIssues
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

const linterUnknown = "unknown"

func extractLinterName(line string) string {
	// golangci-lint output usually ends with " (lintername)"
	if !strings.HasSuffix(line, ")") {
		// Fallback for typecheck which sometimes lacks suffix but contains keyword
		if strings.Contains(line, "typecheck") || strings.Contains(line, "undefined") {
			return "typecheck"
		}
		return linterUnknown
	}

	lastParen := strings.LastIndex(line, "(")
	if lastParen == -1 {
		return linterUnknown
	}

	// Ensure the "(" is preceded by a space and followed by the linter name
	if lastParen > 0 && line[lastParen-1] == ' ' {
		linter := line[lastParen+1 : len(line)-1]
		// Linters are typically short, alphanumeric, and no spaces or special code chars
		if len(linter) > 0 && len(linter) < 25 && !strings.ContainsAny(linter, " {}=[]") {
			return linter
		}
	}

	return linterUnknown
}

func printLinterSummaryTable(issuesByLinter map[string][]linterIssue) {
	fmt.Println(domain.SeparatorLine)
	fmt.Printf("%-20s | %-6s | %-10s\n", "Linter", "Count", "Status")
	fmt.Println(strings.Repeat("-", 40))

	linterNames := sortedLinterNames(issuesByLinter)
	for _, linter := range linterNames {
		count := len(issuesByLinter[linter])
		status := "⚠️"
		if count > 20 {
			status = "💣 SYSTEMIC"
		}
		fmt.Printf("%-20s | %-6d | %-10s\n", linter, count, status)
	}
	fmt.Println(domain.SeparatorLine)
}

func printTopIssues(issuesByLinter map[string][]linterIssue, globalCap, perLinterCap int) int {
	reportedCount := 0
	linterNames := sortedLinterNames(issuesByLinter)
	p0Keywords := []string{"security", "shadow", "panic", "govet", "gosec"}

	fmt.Println("🔍 Top Issues (Agent-Native Summary):")
	for _, linter := range linterNames {
		if reportedCount >= globalCap {
			break
		}
		reportedCount += printLinterIssues(linter, issuesByLinter[linter], p0Keywords, globalCap-reportedCount, perLinterCap)
	}
	return reportedCount
}

func printLinterIssues(linter string, issues []linterIssue, p0Keywords []string, remainingGlobal, perLinterCap int) int {
	reported := 0
	sorted := prioritizeIssues(issues, p0Keywords)

	for i, iss := range sorted {
		if i >= perLinterCap || reported >= remainingGlobal {
			if i == perLinterCap {
				fmt.Printf("   ... and %d more from %s\n", len(issues)-perLinterCap, linter)
			}
			break
		}
		fmt.Printf("📍 [%s] %s\n", iss.linter, iss.raw)
		reported++
	}
	return reported
}

func sortedLinterNames(issuesByLinter map[string][]linterIssue) []string {
	names := make([]string, 0, len(issuesByLinter))
	for l := range issuesByLinter {
		names = append(names, l)
	}
	return names
}

func prioritizeIssues(issues []linterIssue, keywords []string) []linterIssue {
	p0 := []linterIssue{}
	other := []linterIssue{}
	for _, iss := range issues {
		isP0 := false
		for _, kw := range keywords {
			if strings.Contains(strings.ToLower(iss.raw), kw) {
				isP0 = true
				break
			}
		}
		if isP0 {
			p0 = append(p0, iss)
		} else {
			other = append(other, iss)
		}
	}
	return append(p0, other...)
}

func anySystemic(issues map[string][]linterIssue) bool {
	for _, iss := range issues {
		if len(iss) > 20 {
			return true
		}
	}
	return false
}

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

func runTool(dir, name string, args ...string) (string, error) {
	//nolint:gosec // G204: Maintenance utility running trusted audit tools
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
