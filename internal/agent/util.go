package agent

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
)

// NormalizePath consolidates path normalization logic used across the agent.
func NormalizePath(p string) string {
	// 1. replace :id or :UserId with {id} or {UserId}
	p = regexp.MustCompile(`:([a-zA-Z0-9_]+)`).ReplaceAllString(p, "{$1}")

	// 2. Remove trailing slashes (except root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}

	// 3. Deduplicate slashes
	p = regexp.MustCompile(`//+`).ReplaceAllString(p, "/")

	if p == "" {
		p = "/"
	}
	return p
}

// SpawnAgent runs the antigravity chat command asynchronously and detached.
// It executes "antigravity chat -m agent -a task.md <prompt>".
func SpawnAgent(prompt string) error {
	cmd := ExecCommand("antigravity", "chat", "-m", "agent", "-a", "task.md", prompt)

	// Detach the process from the current session/terminal
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Redirect output to avoid blocking the harness CLI
	// #nosec G108 - Redirecting to null is intentional for detached background task
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err == nil {
		cmd.Stdout = null
		cmd.Stderr = null
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Error: Failed to spawn background agent: %v\n", err)
		return err
	}
	return nil
}
