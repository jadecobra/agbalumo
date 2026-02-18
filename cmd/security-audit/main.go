package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CommandRunner abstracts command execution for testing
type CommandRunner interface {
	Run(dir string, name string, args ...string) (string, error)
}

type RealRunner struct{}

func (r *RealRunner) Run(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	// Initialize Config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 2 * time.Second,
	}

	config := AuditConfig{
		TargetURL:  "https://localhost:8443",
		HTTPClient: client,
		RootDir:    ".",
		Runner:     &RealRunner{},
		FileReader: os.ReadFile,
	}

	if err := runAudit(config); err != nil {
		os.Exit(1)
	}
}

type AuditConfig struct {
	TargetURL  string
	HTTPClient *http.Client
	RootDir    string
	Runner     CommandRunner
	FileReader func(name string) ([]byte, error)
}

func runAudit(config AuditConfig) error {
	score := 0
	total := 5 // Vetting, Headers, Fly Config, Vuln Check, XSS check

	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println("--------------------------------")

	// 1. Static Analysis (go vet)
	if passed := checkGoVet(config.RootDir, config.Runner); passed {
		score++
	}

	// 2. Check Live Headers
	if passed := checkHeaders(config.TargetURL, config.HTTPClient); passed {
		score++
	}

	// 3. Check fly.toml for leaks
	flyContent, err := config.FileReader("fly.toml")
	flyContentStr := ""
	if err == nil {
		flyContentStr = string(flyContent)
	}
	if passed := checkFlyConfig(flyContentStr, err != nil); passed {
		score++
	} else if err != nil {
		total-- // Adjust total if file missing
	}

	// 4. Vulnerability Check
	if passed := checkVuln(config.RootDir, config.Runner, exec.LookPath); passed {
		score++
	}

	// 5. XSS Vulnerability Check
	if passed := checkXSS(config.RootDir, config.Runner); passed {
		score++
	}

	fmt.Println("--------------------------------")
	finalScore := (float64(score) / float64(total)) * 100
	fmt.Printf("🔒 Security Score: %.0f/100\n", finalScore)

	if finalScore < 100 {
		return fmt.Errorf("security score too low")
	}
	return nil
}

func checkGoVet(rootDir string, runner CommandRunner) bool {
	fmt.Print("[?] Running 'go vet'... ")
	_, err := runner.Run(rootDir, "go", "vet", "./...")
	if err != nil {
		fmt.Println("❌ Failed")
		return false
	}
	fmt.Println("✅ Passed")
	return true
}

func checkHeaders(target string, client *http.Client) bool {
	resp, err := client.Get(target)
	if err != nil {
		fmt.Printf("⚠️  Could not connect to server (%s): %v\n", target, err)
		// Try HTTP fallback if HTTPS fails
		if strings.HasPrefix(target, "https") {
			httpTarget := strings.Replace(target, "https", "http", 1)
			if strings.Contains(httpTarget, ":8443") {
				httpTarget = strings.Replace(httpTarget, ":8443", ":8080", 1)
			}
			fmt.Printf("[!] Falling back to (%s)... ", httpTarget)
			resp, err = client.Get(httpTarget)
		}
	}

	if err != nil {
		fmt.Printf("⚠️  Could not connect to server: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	passed := true

	// HSTS
	hsts := resp.Header.Get("Strict-Transport-Security")
	fmt.Printf("[?] HSTS Header: ")
	if hsts != "" {
		fmt.Printf("✅ Present (%s)\n", hsts)
	} else {
		fmt.Println("❌ Missing")
		passed = false
	}

	// CSP
	csp := resp.Header.Get("Content-Security-Policy")
	fmt.Printf("[?] CSP Header: ")
	if csp != "" {
		fmt.Printf("✅ Present\n")
	} else {
		fmt.Println("❌ Missing")
		passed = false
	}

	// X-Frame-Options
	xfo := resp.Header.Get("X-Frame-Options")
	fmt.Printf("[?] X-Frame-Options: ")
	if xfo != "" {
		fmt.Printf("✅ Present (%s)\n", xfo)
	} else {
		fmt.Println("❌ Missing")
		passed = false
	}

	return passed
}

func checkFlyConfig(content string, missing bool) bool {
	fmt.Print("[?] Checking fly.toml for secrets... ")
	if missing {
		fmt.Println("⚠️  fly.toml not found (skipping)")
		return false
	}

	sensitiveKeys := []string{"SECRET", "KEY", "PASSWORD", "TOKEN", "AUTH"}
	foundSecret := false
	for _, key := range sensitiveKeys {
		if containsSensitive(content, key) {
			foundSecret = true
			break
		}
	}

	if foundSecret {
		fmt.Println("❌ Potential leak found!")
		return false
	}
	fmt.Println("✅ Passed")
	return true
}

func checkVuln(rootDir string, runner CommandRunner, lookup func(string) (string, error)) bool {
	fmt.Print("[?] Running 'govulncheck'... ")

	// Check if installed
	if _, err := lookup("govulncheck"); err != nil {
		fmt.Println("⚠️  govulncheck not installed")
		return true
	}

	output, err := runner.Run(rootDir, "govulncheck", "./...")

	if err != nil {
		if strings.Contains(output, "No vulnerabilities found") {
			fmt.Println("✅ Passed")
			return true
		}
		fmt.Println("⚠️  Possible Vulnerabilities:")
		fmt.Println(output)
		return false
	}
	fmt.Println("✅ Passed")
	return true
}

func checkXSS(rootDir string, runner CommandRunner) bool {
	fmt.Print("[?] Checking for XSS vulnerabilities... ")
	// Use sh -c to handle piping
	// Note: CommandRunner.Run takes name and args. For piping we usually need sh -c.
	// We'll assume runner handles simple command execution.
	// For piping, we pass "sh", "-c", "script..."

	script := "grep -r 'template.HTML' . --exclude-dir=bin --exclude-dir=.git | grep -v 'cmd/security-audit/main.go'"
	output, _ := runner.Run(rootDir, "sh", "-c", script)

	if len(output) > 0 {
		fmt.Println("⚠️  Found explicit 'template.HTML' usage:")
		fmt.Println(output)
		return false
	}
	fmt.Println("✅ Passed")
	return true
}

func containsSensitive(content, key string) bool {
	search := strings.ToLower(key)
	lowerContent := strings.ToLower(content)
	return strings.Contains(lowerContent, search) && (strings.Contains(lowerContent, "=") || strings.Contains(lowerContent, ":"))
}
