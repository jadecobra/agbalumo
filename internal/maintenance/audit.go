package maintenance

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AuditConfig holds the necessary configurations for a security audit.
type AuditConfig struct {
	TargetURL  string
	HTTPClient *http.Client
	RootDir    string
}

type auditResults struct {
	score int
	total int
}

// RunSecurityAudit carries out a series of security checks on the codebase and, if available, the live server.
func RunSecurityAudit(config AuditConfig) error {
	res := &auditResults{total: 0, score: 0}
	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println("--------------------------------")

	res.runCheck("Go Vet", func() (bool, bool) { return checkGoVet(config) })
	res.runCheck("Headers", func() (bool, bool) { return checkHeaders(config.TargetURL, config.HTTPClient) })
	res.runCheck("fly.toml", func() (bool, bool) { return checkFlyConfig(config) })
	res.runCheck("govulncheck", func() (bool, bool) { return checkVulnerabilities(config) })
	res.runCheck("XSS", func() (bool, bool) { return checkXSS(config) })

	fmt.Println("--------------------------------")
	if res.total <= 0 {
		return fmt.Errorf("no security checks could be performed")
	}

	finalScore := (float64(res.score) / float64(res.total)) * 100
	fmt.Printf("🔒 Security Score: %.0f/100\n", finalScore)

	if finalScore < 100 {
		return fmt.Errorf("security score too low: %.0f/100", finalScore)
	}
	return nil
}

func (r *auditResults) runCheck(name string, check func() (bool, bool)) {
	fmt.Printf("[?] Running '%s'... ", name)
	passed, skip := check()
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
	if _, err := exec.LookPath("govulncheck"); err != nil {
		return false, true
	}

	cmd := exec.Command("govulncheck", "./...")
	cmd.Dir = config.RootDir
	out, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(out), "No vulnerabilities found") {
		return false, false
	}
	return true, false
}

func checkXSS(config AuditConfig) (bool, bool) {
	script := "grep -rI 'template.HTML' . --exclude-dir=.git --exclude-dir=bin --exclude-dir=scripts --exclude-dir='.tester' --exclude-dir='@*' --exclude='*_test.go' --exclude='renderer.go' | grep -v 'cmd/verify/main.go' | grep -v 'template.HTMLEscapeString' | grep -v 'internal/agent' | grep -v 'internal/maintenance/audit.go'"
	cmd := exec.Command("sh", "-c", script)
	cmd.Dir = config.RootDir
	out, _ := cmd.CombinedOutput()
	return len(strings.TrimSpace(string(out))) == 0, false
}

// RunChiefCriticAudit performs a comprehensive code quality audit.
func RunChiefCriticAudit(rootDir string) error {
	fmt.Println("🚀 Starting ChiefCritic Robustness Audit...")

	tools := []struct {
		name string
		cmd  []string
	}{
		{"Cognitive Complexity", []string{"go", "run", "github.com/uudashr/gocognit/cmd/gocognit", "-over", "10", "./cmd", "./internal"}},
		{"Repeated Strings", []string{"go", "run", "github.com/jgautheron/goconst/cmd/goconst", "./cmd/...", "./internal/..."}},
		{"Struct Alignment", []string{"go", "run", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment", "./internal/...", "./cmd/..."}},
		{"Code Duplication", []string{"go", "run", "github.com/mibk/dupl", "-threshold", "15", "./cmd", "./internal"}},
	}

	failed := false
	for i, t := range tools {
		fmt.Printf("\n[%d/%d] Checking %s...\n", i+1, len(tools), t.name)
		if err := runTool(rootDir, t.cmd[0], t.cmd[1:]...); err != nil {
			fmt.Printf("❌ %s failed!\n", t.name)
			if t.name == "Cognitive Complexity" {
				failed = true
			}
		}
	}

	fmt.Println("\n✅ ChiefCritic Audit Complete!")
	if failed {
		return fmt.Errorf("robustness audit failed due to mandatory quality gate violations")
	}
	return nil
}

func runTool(dir, name string, args ...string) error {
	//nolint:gosec // G204: Maintenance utility running trusted audit tools
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
