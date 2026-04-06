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

type watcherState struct {
	cmdName    string
	cmdArgs    []string
	currentCmd *exec.Cmd
}

// Watch monitors the codebase for changes and executes a command (e.g., serve or test).
func Watch(ctx context.Context, cmdName string, cmdArgs []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer func() { _ = watcher.Close() }()

	state := &watcherState{cmdName: cmdName, cmdArgs: cmdArgs}
	fmt.Printf("👀 Watcher started. Initial execution: %s %s\n", cmdName, strings.Join(cmdArgs, " "))

	state.run()
	setupDirs(watcher)

	return watchLoop(ctx, watcher, state)
}

func (s *watcherState) run() {
	if s.currentCmd != nil && s.currentCmd.Process != nil {
		fmt.Println("🔄 Restarting...")
		_ = s.currentCmd.Process.Kill()
	}

	// G204: Maintenance utility executes command provided by user/config
	s.currentCmd = exec.Command(s.cmdName, s.cmdArgs...) //nolint:gosec // maintenance utility
	s.currentCmd.Stdout = os.Stdout
	s.currentCmd.Stderr = os.Stderr
	if startErr := s.currentCmd.Start(); startErr != nil {
		fmt.Printf("❌ Failed to start command: %v\n", startErr)
	}
}

func setupDirs(watcher *fsnotify.Watcher) {
	dirs := []string{".", "cmd", "internal", "ui/templates"}
	for _, d := range dirs {
		_ = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() && shouldWatch(path) {
				_ = watcher.Add(path)
			}
			return nil
		})
	}
}

func shouldWatch(path string) bool {
	if path == "node_modules" || path == "vendor" || path == ".git" {
		return false
	}
	if strings.HasPrefix(path, ".") && path != "." {
		return false
	}
	return true
}

func watchLoop(ctx context.Context, watcher *fsnotify.Watcher, state *watcherState) error {
	debounce := time.NewTimer(300 * time.Millisecond)
	if !debounce.Stop() {
		<-debounce.C
	}

	for {
		select {
		case <-ctx.Done():
			state.cleanup()
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			state.handleEvent(event, debounce)
		case <-debounce.C:
			state.run()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (s *watcherState) cleanup() {
	if s.currentCmd != nil && s.currentCmd.Process != nil {
		_ = s.currentCmd.Process.Kill()
	}
}

func (s *watcherState) handleEvent(event fsnotify.Event, debounce *time.Timer) {
	if isInterestingFile(event) {
		debounce.Reset(300 * time.Millisecond)
	}
}

func isInterestingFile(event fsnotify.Event) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
		return false
	}
	ext := filepath.Ext(event.Name)
	return ext == ".go" || ext == ".html" || ext == ".css"
}
