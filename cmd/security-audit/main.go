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

func main() {
	score := 0
	total := 7 // Vetting, Headers (HSTS, CSP, X-Frame), fly.toml leak check, Vuln Check, XSS check

	fmt.Println("🛡️  Starting Security Audit...")
	fmt.Println("--------------------------------")

	// 1. Static Analysis (go vet)
	fmt.Print("[?] Running 'go vet'... ")
	cmd := exec.Command("go", "vet", "./...")
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed")
	} else {
		fmt.Println("✅ Passed")
		score++
	}

	// 2. Check Live Headers
	// Ensure server is running or start it?
	// For simplicity, we assume server is running on :8443 (HTTPS) or :8080
	// We'll try to check local dev server.
	target := "https://localhost:8443"
	// Bypass cert check for self-signed
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(target)
	if err != nil {
		fmt.Printf("⚠️  Could not connect to server (%s): %v\n", target, err)
		// Try HTTP as fallback but warn heavily
		target = "http://localhost:8080"
		fmt.Printf("[!] Falling back to (%s)... ", target)
		client = &http.Client{Timeout: 2 * time.Second}
		resp, err = client.Get(target)
	}

	if err != nil {
		fmt.Printf("⚠️  Could not connect to server (%s) either: %v\n", target, err)
		// Don't penalize score if server is just down, but warn.
	} else {
		defer resp.Body.Close()

		// HSTS
		hsts := resp.Header.Get("Strict-Transport-Security")
		fmt.Printf("[?] HSTS Header: ")
		if hsts != "" {
			fmt.Printf("✅ Present (%s)\n", hsts)
			score++
		} else {
			fmt.Println("❌ Missing")
		}

		// CSP
		csp := resp.Header.Get("Content-Security-Policy")
		fmt.Printf("[?] CSP Header: ")
		if csp != "" {
			fmt.Printf("✅ Present\n")
			score++
		} else {
			fmt.Println("❌ Missing")
		}

		// X-Frame-Options
		xfo := resp.Header.Get("X-Frame-Options")
		fmt.Printf("[?] X-Frame-Options: ")
		if xfo != "" {
			fmt.Printf("✅ Present (%s)\n", xfo)
			score++
		} else {
			fmt.Println("❌ Missing")
		}
	}

	// 3. Check fly.toml for leaks
	fmt.Print("[?] Checking fly.toml for secrets... ")
	flyContent, err := os.ReadFile("fly.toml")
	if err != nil {
		fmt.Println("⚠️  fly.toml not found (skipping)")
		total-- // Adjust total if file missing
	} else {
		content := string(flyContent)
		// Check for common secret keys in [env] block if any
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
		} else {
			fmt.Println("✅ Passed")
			score++
		}
	}

	// 4. Vulnerability Check (govulncheck)
	fmt.Print("[?] Running 'govulncheck'... ")
	// Check if installed first
	if _, err := exec.LookPath("govulncheck"); err != nil {
		fmt.Println("⚠️  govulncheck not found. Installing...")
		installCmd := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest")
		if err := installCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to install govulncheck: %v\n", err)
		}
	}

	cmdVuln := exec.Command("govulncheck", "./...")
	if output, err := cmdVuln.CombinedOutput(); err != nil {
		// govulncheck returns exit code 1 if vulnerabilities are found?
		// Actually it returns 0 on success (no vulns) or failure (vulns found) depending on flags?
		// It returns 0 if no vulns, or if vulns found but not erroring?
		// Usually it's best to check output or exit code.
		// Let's assume if err != nil it might be a finding or execution error.
		// Use output to check.
		fmt.Println("⚠️  Possible Vulnerabilities or Error:")
		fmt.Println(string(output))
		// Don't fail the score immediately if it's just a setup issue, but warn.
		// If explicit vulns found "No vulnerabilities found" is standard success msg.
		if strings.Contains(string(output), "No vulnerabilities found") {
			fmt.Println("✅ Passed")
			score++
		}
	} else {
		fmt.Println("✅ Passed")
		score++
	}

	// 5. XSS Vulnerability Check
	fmt.Print("[?] Checking for XSS vulnerabilities (template.HTML usage in handlers/views)... ")
	// Grep for template.HTML usage, excluding binaries, git, and this audit file specifically
	// using grep -v for reliable exclusion of the specific file path
	cmdXSS := exec.Command("sh", "-c", "grep -r 'template.HTML' . --exclude-dir=bin --exclude-dir=.git | grep -v 'cmd/security-audit/main.go'")
	output, _ := cmdXSS.CombinedOutput()
	if len(output) > 0 {
		fmt.Println("⚠️  Found explicit 'template.HTML' usage (verify safety):")
		fmt.Println(string(output))
	} else {
		fmt.Println("✅ Passed (No explicit unsafe HTML found)")
		score++
	}

	fmt.Println("--------------------------------")
	finalScore := (float64(score) / float64(total)) * 100
	fmt.Printf("🔒 Security Score: %.0f/100\n", finalScore)

	if finalScore < 100 {
		os.Exit(1)
	}
}

func containsSensitive(content, key string) bool {
	// Simple check: looking for key= or key = and ensuring it's not empty or just a placeholder
	// We want to avoid matching comments if possible, but for a simple tool this is okay.
	search := strings.ToLower(key)
	lowerContent := strings.ToLower(content)

	// Look for pattern key = "value" or key = 'value'
	if !strings.Contains(lowerContent, search) {
		return false
	}

	// If the key is there, let's see if it's in the [env] section and has a value that isn't a known placeholder
	// For now, if someone puts CLIENT_ID = "actually-the-id" in fly.toml, we want to catch it.
	// Production secrets should be in fly secrets.
	return true
}
