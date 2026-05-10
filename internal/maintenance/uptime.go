package maintenance

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// CheckServerUptime verifies that the server is responding on the given targetURL.
// It performs an HTTP GET and expects a 200 OK response.
func CheckServerUptime(targetURL string) error {
	if targetURL == "" {
		return fmt.Errorf("target URL is empty; check APP_URL environment variable")
	}

	fmt.Printf("🔍 Checking server uptime at %s...\n", targetURL)

	client := &http.Client{
		Transport: &http.Transport{
			// #nosec G402 // local development often uses self-signed certs
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return fmt.Errorf("server is unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	fmt.Println("✅ Server is UP and responding.")
	return nil
}
