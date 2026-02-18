package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Agbalumo web server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file
		// Try loading from local .env or the scripts location
		godotenv.Load(".env")
		if err := godotenv.Load("../scripts/agbalumo/.env"); err != nil {
			log.Printf("Error loading ../scripts/agbalumo/.env: %v", err)
		}

		// Environment Configuration
		env := os.Getenv("AGBALUMO_ENV")

		// Setup Server
		e, err := SetupServer()
		if err != nil {
			log.Fatalf("Failed to setup server: %v", err)
		}

		// Server Config
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		// In production (Fly.io), TLS is handled by the proxy. We just listen on PORT.
		// In dev, we might want TLS if certificates exist, OR just HTTP.
		if env == "production" {
			log.Printf("Starting Server in PRODUCTION mode on :%s", port)
			if os.Getenv("AGBALUMO_DRY_RUN") == "true" {
				return
			}
			if err := e.Start(":" + port); err != nil {
				e.Logger.Fatal(err)
			}
		} else {
			// Development Mode
			certFile := "certs/cert.pem"
			keyFile := "certs/key.pem"

			if _, err := os.Stat(certFile); err == nil {
				if _, err := os.Stat(keyFile); err == nil {
					log.Println("Starting Secure Server on :8443 (HTTPS)")
					if os.Getenv("AGBALUMO_DRY_RUN") == "true" {
						return
					}
					if err := e.StartTLS(":8443", certFile, keyFile); err != nil {
						e.Logger.Fatal(err)
					}
					return
				}
			}

			log.Printf("Starting Server in DEV mode on :%s (HTTP)", port)
			if os.Getenv("AGBALUMO_DRY_RUN") == "true" {
				return
			}
			if err := e.Start(":" + port); err != nil {
				e.Logger.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
