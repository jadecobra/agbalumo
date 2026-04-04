# Phase 12: Graceful Shutdown Implementation

## Objective
Enforce strict zero-downtime container termination by intercepting OS signals and allowing `Echo` and `SQLite` to drain in-flight web requests and database WAL checkpoints before process death.

## Context
When deploying via container platforms (Docker/Fly.io) or restarting locally, pulling the plug on the Go binary immediately kills running HTTP handlers in the middle of executing. This can drop connections on users and occasionally cause SQLite lock issues if an insert is violently interrupted.

## Steps for Execution
1. Open up the entry file where Echo starts (likely `main.go` or inside `cmd/server/main.go`).
2. Implement signal capturing for `SIGTERM` and `SIGINT` using a `goroutine` to wrap the blocking `e.Start()` call:

```go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/infra/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	e, err := server.Setup(cfg)
	if err != nil {
		slog.Error("Server setup failed", "error", err)
		os.Exit(1)
	}

	// 1. Run the server in a goroutine
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server forced to shutdown", "error", err)
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
    
	// 4. (Optional) Access the repository and explicitly close the *sql.DB here if accessible
    
	slog.Info("Server exited properly")
}
```
3. Run `go run cmd/verify/main.go ci` to ensure tests are not broken by the startup refactor. 
4. Commit natively: `feat(infra): implement strict graceful shutdown context trap`.

## Verification
- Running `go run main.go serve` and hitting `CTRL+C` outputs "Shutdown signal received..." rather than instantly snapping back to the terminal prompt.
