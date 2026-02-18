package cmd

import (
	"log"
	"os"

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
	if port == "" {
		port = "8080"
	}

	// In production (Fly.io), TLS is handled by the proxy. We just listen on PORT.
	if env == "production" {
		return ServerConfig{Addr: ":" + port, TLS: false}
	}

	// Development Mode
	certFile := "certs/cert.pem"
	keyFile := "certs/key.pem"

	if fileExists(certFile) && fileExists(keyFile) {
		return ServerConfig{
			Addr:     ":8443",
			TLS:      true,
			CertFile: certFile,
			KeyFile:  keyFile,
		}
	}

	return ServerConfig{Addr: ":" + port, TLS: false}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Agbalumo web server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file
		godotenv.Load(".env")
		if err := godotenv.Load("../scripts/agbalumo/.env"); err != nil {
			log.Printf("Error loading ../scripts/agbalumo/.env: %v", err)
		}

		// Environment Configuration
		env := os.Getenv("AGBALUMO_ENV")
		port := os.Getenv("PORT")

		// Setup Server
		e, err := SetupServer()
		if err != nil {
			log.Fatalf("Failed to setup server: %v", err)
		}

		// Resolve Configuration
		config := ResolveServerConfig(env, port, func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		})

		// Dry Run Check
		if os.Getenv("AGBALUMO_DRY_RUN") == "true" {
			log.Printf("Dry Run: Config resolved to %+v", config)
			return
		}

		if config.TLS {
			log.Printf("Starting Secure Server on %s (HTTPS)", config.Addr)
			if err := e.StartTLS(config.Addr, config.CertFile, config.KeyFile); err != nil {
				e.Logger.Fatal(err)
			}
		} else {
			mode := "DEV"
			if env == "production" {
				mode = "PRODUCTION"
			}
			log.Printf("Starting Server in %s mode on %s (HTTP)", mode, config.Addr)
			if err := e.Start(config.Addr); err != nil {
				e.Logger.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
