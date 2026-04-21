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

// VerifyServerStartup starts the server and verifies it boots successfully without panics.
func VerifyServerStartup(rootDir string) error {
	fmt.Println("🚀 Verifying Server Startup (Dynamic Audit)...")
	
	testerDir := filepath.Join(rootDir, ".tester", "servers")
	if err := os.MkdirAll(testerDir, 0750); err != nil {
		return fmt.Errorf("failed to create tester directory: %w", err)
	}

	logPath := filepath.Join(testerDir, "verify_startup.log")
	logFile, err := os.Create(filepath.Clean(logPath))
	if err != nil {
		return fmt.Errorf("failed to create startup log: %w", err)
	}
	defer func() { _ = logFile.Close() }()

	binPath, err := buildTestBinary(rootDir, testerDir)
	if err != nil {
		return err
	}

	cmd, done, err := startTestServer(binPath, rootDir, logFile)
	if err != nil {
		return err
	}

	return waitForReadiness(cmd, done, logPath)
}

func buildTestBinary(rootDir, testerDir string) (string, error) {
	fmt.Println("🔨 Building temporary test binary...")
	binPath := filepath.Join(testerDir, "server_bin")
	// #nosec G204
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Dir = rootDir
	if err := buildCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build server for dynamic audit: %w", err)
	}
	return binPath, nil
}

func startTestServer(binPath, rootDir string, logFile *os.File) (*exec.Cmd, chan error, error) {
	// #nosec G204
	cmd := exec.Command(binPath, "serve")
	cmd.Dir = rootDir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Env = append(os.Environ(), "AGBALUMO_ENV=test", "PORT=8444")

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start server: %w", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	return cmd, done, nil
}

func waitForReadiness(cmd *exec.Cmd, done chan error, logPath string) error {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, //nolint:gosec // maintenance utility uses self-signed local certs
		Timeout:   1 * time.Second,
	}
	
	fmt.Print("⏳ Waiting for server readiness...")
	for i := 0; i < 15; i++ {
		if ready, err := attemptReadinessCheck(client, done, logPath); err != nil {
			return err
		} else if ready {
			_ = cmd.Process.Kill()
			return nil
		}
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}

	_ = cmd.Process.Kill()
	return formatStartupError(logPath, fmt.Errorf("readiness timeout"))
}

func attemptReadinessCheck(client *http.Client, done chan error, logPath string) (bool, error) {
	select {
	case err := <-done:
		fmt.Println("\n❌ Server process exited prematurely!")
		return false, formatStartupError(logPath, err)
	default:
		resp, err := client.Get("https://localhost:8444/healthz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				fmt.Println(" ✅ Ready!")
				return true, nil
			}
		}
		return false, nil
	}
}

func formatStartupError(logPath string, err error) error {
	content, _ := os.ReadFile(filepath.Clean(logPath))
	if strings.Contains(string(content), "panic:") {
		return fmt.Errorf("server panicked during startup:\n%s", string(content))
	}
	return fmt.Errorf("server failed to start (exit error: %v). Logs:\n%s", err, string(content))
}
