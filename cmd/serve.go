package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// ServerConfig holds the configuration for starting the server
type ServerConfig struct {
	Addr     string
	TLS      bool
	CertFile string
	KeyFile  string
}

// ResolveServerConfig determines the server configuration based on environment and file existence
func ResolveServerConfig(env, port string, fileExists func(string) bool) ServerConfig {
	const certFile = "certs/cert.pem"
	const keyFile = "certs/key.pem"
	hasCerts := fileExists(certFile) && fileExists(keyFile)
	port = resolvePort(port, env, hasCerts)

	// In production (Fly.io), TLS is handled by the proxy. We just listen on PORT.
	if env == "production" {
		return ServerConfig{Addr: ":" + port, TLS: false}
	}

	// Development Mode
	if hasCerts {
		return ServerConfig{
			Addr:     ":" + port,
			TLS:      true,
			CertFile: certFile,
			KeyFile:  keyFile,
		}
	}

	return ServerConfig{Addr: ":" + port, TLS: false}
}

// resolvePort returns the effective port to listen on.
// It checks: explicit arg → APP_URL env → cert-aware default.
func resolvePort(port, env string, hasCerts bool) string {
	if port != "" {
		return port
	}

	if appURL := os.Getenv("APP_URL"); appURL != "" {
		if strings.Contains(appURL, ":") {
			parts := strings.Split(appURL, ":")
			port = strings.TrimSuffix(parts[len(parts)-1], "/")
			if port != "" {
				return port
			}
		}
	}

	if hasCerts && env != "production" {
		return "8443"
	}

	return "8080"
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the agbalumo web server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file
		_ = godotenv.Load(".env")

		// Environment Configuration
		env := os.Getenv("AGBALUMO_ENV")
		port := os.Getenv("PORT")

		// Setup Server
		e, cleanup, err := SetupServer()
		if err != nil {
			slog.Error("Failed to setup server", "error", err)
			os.Exit(1)
		}

		// Resolve Configuration
		config := ResolveServerConfig(env, port, func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		})

		// Dry Run Check
		if os.Getenv("AGBALUMO_DRY_RUN") == "true" {
			slog.Info("Dry Run configuration", "config", config)
			return
		}

		// 1. Run the server in a goroutine
		go func() {
			if config.TLS {
				slog.Info("Starting Secure Server (HTTPS)", "addr", config.Addr)
				if err := e.StartTLS(config.Addr, config.CertFile, config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("Server forced to shutdown", "error", err)
				}
			} else {
				mode := "DEV"
				if env == "production" {
					mode = "PRODUCTION"
				}
				slog.Info("Starting Server (HTTP)", "mode", mode, "addr", config.Addr)
				if err := e.Start(config.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("Server forced to shutdown", "error", err)
				}
			}
		}()

		// 2. Setup channel to listen for OS interrupt signals
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		// Block until a signal is received
		<-quit
		slog.Info("Shutdown signal received, initiating graceful teardown...")

		// 3. Create context with 15 second timeout to allow in-flight requests to drain
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			slog.Error("Echo shutdown failed gracefully", "error", err)
		}

		// 4. Run cleanup (close DB, stop background services)
		cleanup()

		slog.Info("Server exited properly")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
