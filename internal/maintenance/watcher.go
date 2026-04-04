package maintenance

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch monitors the codebase for changes and executes a command (e.g., serve or test).
// This is an agent-friendly replacement for legacy shell loops.
func Watch(ctx context.Context, cmdName string, cmdArgs []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Initial run
	fmt.Printf("👀 Watcher started. Initial execution: %s %s\n", cmdName, strings.Join(cmdArgs, " "))
	var currentCmd *exec.Cmd

	run := func() {
		if currentCmd != nil && currentCmd.Process != nil {
			fmt.Println("🔄 Restarting...")
			_ = currentCmd.Process.Kill()
		}

		// G204: Maintenance utility executes command provided by user/config
		currentCmd = exec.Command(cmdName, cmdArgs...) //nolint:gosec // maintenance utility
		currentCmd.Stdout = os.Stdout
		currentCmd.Stderr = os.Stderr
		if startErr := currentCmd.Start(); startErr != nil {
			fmt.Printf("❌ Failed to start command: %v\n", startErr)
		}
	}

	run()

	// Add directories to watch
	dirs := []string{".", "cmd", "internal", "ui/templates"}
	for _, d := range dirs {
		err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if strings.HasPrefix(path, ".") && path != "." {
					return filepath.SkipDir
				}
				if path == "node_modules" || path == "vendor" || path == ".git" {
					return filepath.SkipDir
				}
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			log.Printf("Warning: failed to watch directory %s: %v", d, err)
		}
	}

	debounceTimer := time.NewTimer(500 * time.Millisecond)
	if !debounceTimer.Stop() {
		<-debounceTimer.C
	}

	for {
		select {
		case <-ctx.Done():
			if currentCmd != nil && currentCmd.Process != nil {
				_ = currentCmd.Process.Kill()
			}
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// Only trigger on write, create, or remove of Go/HTML files
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				ext := filepath.Ext(event.Name)
				if ext == ".go" || ext == ".html" || ext == ".css" {
					debounceTimer.Reset(300 * time.Millisecond)
				}
			}
		case <-debounceTimer.C:
			run()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
