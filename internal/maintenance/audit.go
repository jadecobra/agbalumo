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

// RunChiefCriticAudit performs a comprehensive code quality audit.
func RunChiefCriticAudit(rootDir string) error {
	fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")
	failed := executeAuditTools(rootDir)
	fmt.Println("\n✅ ChiefCritic Audit Complete!")
	if failed {
		return fmt.Errorf("robustness audit failed due to mandatory quality gate violations")
	}
	return nil
}

func executeAuditTools(rootDir string) bool {
	tools := []struct {
		name      string
		cmd       []string
		mandatory bool
	}{
		{"Cognitive Complexity", []string{"go", "run", "github.com/uudashr/gocognit/cmd/gocognit", "-over", "10", "./cmd", "./internal"}, true},
		{"Repeated Strings", []string{"go", "run", "github.com/jgautheron/goconst/cmd/goconst", "./cmd/...", "./internal/..."}, false},
		{"Struct Alignment", []string{"go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment", "./internal/...", "./cmd/..."}, false},
		{"Code Duplication", []string{"go", "run", "github.com/mibk/dupl", "-threshold", "15", "./cmd", "./internal"}, false},
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	failed := false
	verbose := os.Getenv("VERBOSE") == "true"

	for i, t := range tools {
		wg.Add(1)
		go func(idx int, name string, command []string, isMandatory bool) {
			defer wg.Done()
			f := runAuditWorker(rootDir, idx, len(tools), name, command, isMandatory, verbose, &mu)
			if f {
				mu.Lock()
				failed = true
				mu.Unlock()
			}
		}(i, t.name, t.cmd, t.mandatory)
	}
	wg.Wait()
	return failed
}

func runAuditWorker(rootDir string, idx, total int, name string, command []string, isMandatory, verbose bool, mu *sync.Mutex) bool {
	output, err := runTool(rootDir, command[0], command[1:]...)

	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("[%d/%d] %s: ", idx+1, total, name)

	failed := false
	if err != nil {
		fmt.Print("❌ ")
		if isMandatory {
			failed = true
		}
	} else {
		fmt.Print("✅ ")
	}

	summary := parseSummary(name, output)
	if summary != "" {
		fmt.Printf("(%s)", summary)
	}
	fmt.Println()

	if (err != nil && isMandatory) || verbose {
		if output != "" {
			fmt.Println(domain.SeparatorLine)
			fmt.Println(output)
			fmt.Println(domain.SeparatorLine)
		}
	}
	return failed
}

func parseSummary(name, output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || output == "" {
		return ""
	}

	switch name {
	case "Code Duplication":
		// dupl output ends with "Found total X clone groups."
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], "Found total") {
				return strings.TrimSpace(lines[i])
			}
		}
	case "Repeated Strings":
		return fmt.Sprintf("found %d violations", len(lines))
	case "Cognitive Complexity":
		return fmt.Sprintf("found %d complexity violations", len(lines))
	}

	if len(lines) > 5 {
		return fmt.Sprintf("%d lines of output", len(lines))
	}
	return ""
}

func runTool(dir, name string, args ...string) (string, error) {
	//nolint:gosec // G204: Maintenance utility running trusted audit tools
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
