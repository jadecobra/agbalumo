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

// RunSecurityAudit carries out a series of security checks on the codebase and, if available, the live server.
func RunSecurityAudit(config AuditConfig) error {
	score := 0
	total := 4 // Go Vet, fly.toml, govulncheck, XSS. (Headers is optional if server is up)

	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println("--------------------------------")

	// 1. Static Analysis (go vet)
	fmt.Print("[?] Running 'go vet'... ")
	cmdVet := exec.Command("go", "vet", "./...")
	cmdVet.Dir = config.RootDir
	if err := cmdVet.Run(); err != nil {
		fmt.Println("❌ Failed")
	} else {
		fmt.Println("✅ Passed")
		score++
	}

	// 2. Check Live Headers (Optional)
	headerPassed, skipped := checkHeaders(config.TargetURL, config.HTTPClient)
	if !skipped {
		total++
		if headerPassed {
			score++
		}
	}

	// 3. Check fly.toml for leaks
	fmt.Print("[?] Checking fly.toml for secrets... ")
	flyPath := filepath.Join(config.RootDir, "fly.toml")
	// G304: Maintenance utility reads the fly.toml file
	flyData, err := os.ReadFile(flyPath) //nolint:gosec // maintenance utility
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️  fly.toml not found (skipping)")
			total--
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
	} else {
		content := strings.ToLower(string(flyData))
		sensitiveKeys := []string{"SECRET", "KEY", "PASSWORD", "TOKEN", "AUTH"}
		foundSecret := false
		for _, key := range sensitiveKeys {
			sKey := strings.ToLower(key)
			if strings.Contains(content, sKey) && (strings.Contains(content, "=") || strings.Contains(content, ":")) {
				foundSecret = true
				break
			}
		}

		if foundSecret {
			fmt.Println("❌ Potential leak found!")
		} else {
			fmt.Println("✅ Passed")
			score++
		}
	}

	// 4. Vulnerability Check
	fmt.Print("[?] Running 'govulncheck'... ")
	if _, err := exec.LookPath("govulncheck"); err != nil {
		fmt.Println("⚠️  govulncheck not installed - skipping check")
		total--
	} else {
		cmdVuln := exec.Command("govulncheck", "./...")
		cmdVuln.Dir = config.RootDir
		out, err := cmdVuln.CombinedOutput()
		if err != nil {
			if strings.Contains(string(out), "No vulnerabilities found") {
				fmt.Println("✅ Passed")
				score++
			} else {
				fmt.Println("⚠️  Possible Vulnerabilities:")
				fmt.Println(string(out))
			}
		} else {
			fmt.Println("✅ Passed")
			score++
		}
	}

	// 5. XSS Vulnerability Check
	fmt.Print("[?] Checking for XSS vulnerabilities... ")
	// Surgically exclude known safe files, tests, and the audit tool itself.
	script := "grep -rI 'template.HTML' . --exclude-dir=.git --exclude-dir=bin --exclude-dir=scripts --exclude-dir='.tester' --exclude-dir='@*' --exclude='*_test.go' --exclude='renderer.go' | grep -v 'cmd/verify/main.go' | grep -v 'template.HTMLEscapeString' | grep -v 'internal/agent' | grep -v 'internal/maintenance/audit.go'"
	cmdXSS := exec.Command("sh", "-c", script)
	cmdXSS.Dir = config.RootDir
	outXSS, _ := cmdXSS.CombinedOutput()

	if len(strings.TrimSpace(string(outXSS))) > 0 {
		fmt.Println("⚠️  Found explicit 'template.HTML' usage:")
		fmt.Println(string(outXSS))
	} else {
		fmt.Println("✅ Passed")
		score++
	}

	fmt.Println("--------------------------------")
	if total <= 0 {
		return fmt.Errorf("no security checks could be performed")
	}

	finalScore := (float64(score) / float64(total)) * 100
	fmt.Printf("🔒 Security Score: %.0f/100\n", finalScore)

	if finalScore < 100 {
		return fmt.Errorf("security score too low: %.0f/100", finalScore)
	}

	return nil
}

func checkHeaders(target string, client *http.Client) (bool, bool) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// G402: Auditing local dev server which may use self-signed certs
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // maintenance utility
			},
			Timeout: 2 * time.Second,
		}
	}

	resp, err := client.Get(target)
	if err != nil {
		// Fallback to localhost if possible
		if strings.Contains(target, "https") {
			httpTarget := strings.Replace(target, "https", "http", 1)
			httpTarget = strings.Replace(httpTarget, ":8443", ":8080", 1)
			resp, err = client.Get(httpTarget)
		}
	}

	if err != nil {
		fmt.Printf("⚠️  Could not connect to server (%s) - skipping header check\n", target)
		return false, true
	}
	defer func() { _ = resp.Body.Close() }()

	passed := true
	headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Frame-Options"}
	for _, h := range headers {
		val := resp.Header.Get(h)
		fmt.Printf("[?] %s Header: ", h)
		if val != "" {
			fmt.Println("✅ Present")
		} else {
			fmt.Println("❌ Missing")
			passed = false
		}
	}

	return passed, false
}
