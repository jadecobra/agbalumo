package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Observe request lifecycle with aggressive logging",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		method, _ := cmd.Flags().GetString("method")
		method = strings.ToUpper(method)

		rootDir := "."
		testerDir := filepath.Join(rootDir, ".tester", "servers")
		if err := os.MkdirAll(testerDir, 0750); err != nil {
			return fmt.Errorf("failed to create tester directory: %w", err)
		}

		logPath := filepath.Join(testerDir, "trace_server.log")
		logFile, err := os.Create(filepath.Clean(logPath))
		if err != nil {
			return fmt.Errorf("failed to create trace log: %w", err)
		}
		defer func() { _ = logFile.Close() }()

		binPath, err := maintenance.BuildTestBinary(rootDir, testerDir)
		if err != nil {
			return err
		}

		// Start server with AGBALUMO_ENV=trace
		fmt.Println("🚀 Starting server in TRACE mode...")
		serverCmd, done, err := maintenance.StartTestServer(binPath, rootDir, logFile, "AGBALUMO_ENV=trace")
		if err != nil {
			return err
		}
		defer func() { _ = serverCmd.Process.Kill() }()

		// Wait for readiness
		client := &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, //nolint:gosec // maintenance utility
			Timeout:   5 * time.Second,
		}

		fmt.Print("⏳ Waiting for server readiness...")
		ready := false
		for i := 0; i < 20; i++ {
			resp, rErr := client.Get("https://localhost:8444/healthz")
			if rErr == nil {
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					fmt.Println(" ✅ Ready!")
					ready = true
					break
				}
			}
			fmt.Print(".")
			select {
			case sErr := <-done:
				return fmt.Errorf("server exited prematurely: %w", sErr)
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}

		if !ready {
			return fmt.Errorf("server failed to become ready")
		}

		// Execute request
		req, err := http.NewRequest(method, "https://localhost:8444"+path, nil)
		if err != nil {
			return err
		}

		fmt.Printf("📡 Executing %s %s...\n", method, path)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()

		// Give a tiny bit of time for logs to flush
		time.Sleep(200 * time.Millisecond)

		// Print server logs
		fmt.Println("\n--- AGGRESSIVE TRACE LOGS ---")
		logContent, lErr := os.ReadFile(filepath.Clean(logPath))
		if lErr != nil {
			return lErr
		}
		fmt.Println(string(logContent))
		fmt.Println("-----------------------------")

		return nil
	},
}

func init() {
	traceCmd.Flags().String("path", "/healthz", "Path to request")
	traceCmd.Flags().String("method", "GET", "HTTP method")
}
