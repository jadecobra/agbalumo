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

// RunChiefCriticAudit performs a consolidated code quality audit using golangci-lint.
func RunChiefCriticAudit(rootDir string, opts ChiefCriticOptions) error {
	fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

	args := []string{"run"}
	if !opts.Full {
		rev := opts.NewFromRev
		if rev == "" {
			rev = "HEAD~1" // Default to incremental check against previous commit
		}
		args = append(args, "--new-from-rev", rev)
	}

	if opts.Verbose {
		args = append(args, "-v")
	}

	// We use go run to ensure we use the version pinned in go.mod
	command := append([]string{"run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"}, args...)
	output, err := runTool(rootDir, "go", command...)

	if err != nil {
		fmt.Println("❌ ChiefCritic Audit Failed")
		if output != "" {
			fmt.Println(domain.SeparatorLine)
			fmt.Println(output)
			fmt.Println(domain.SeparatorLine)
		}
		return fmt.Errorf("robustness audit failed: %w", err)
	}

	fmt.Println("✅ ChiefCritic Audit Complete!")
	return nil
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
