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
