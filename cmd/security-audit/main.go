package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	score := 0
	total := 4 // Vetting, Headers (HSTS, CSP, X-Frame)

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

	fmt.Println("--------------------------------")
	finalScore := (float64(score) / float64(total)) * 100
	fmt.Printf("🔒 Security Score: %.0f/100\n", finalScore)

	if finalScore < 100 {
		os.Exit(1)
	}
}
