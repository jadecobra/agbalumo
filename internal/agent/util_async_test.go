package agent

import (
	"os/exec"
	"testing"
)

func TestSpawnAgent_ArgValidation(t *testing.T) {
	// 1. Setup mock
	oldExec := ExecCommand
	defer func() { ExecCommand = oldExec }()

	var capturedName string
	var capturedArgs []string
	var capturedCmd *exec.Cmd

	ExecCommand = func(name string, args ...string) *exec.Cmd {
		capturedName = name
		capturedArgs = args
		cmd := exec.Command(name, args...)
		capturedCmd = cmd
		return cmd
	}

	// 2. Call SpawnAgent with mock behavior
	prompt := "Test Prompt"
	err := SpawnAgent(prompt)
	_ = err // We don't care if it errors due to missing binary in some envs

	if capturedName != "antigravity" {
		t.Errorf("Expected command antigravity, got %s", capturedName)
	}

	expectedArgs := []string{"chat", "-m", "agent", "-a", "task.md", prompt}
	if len(capturedArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(capturedArgs))
	} else {
		for i, arg := range capturedArgs {
			if arg != expectedArgs[i] {
				t.Errorf("Arg %d: expected %s, got %s", i, expectedArgs[i], arg)
			}
		}
	}

	if capturedCmd == nil || capturedCmd.SysProcAttr == nil || !capturedCmd.SysProcAttr.Setsid {
		t.Errorf("Expected SysProcAttr with Setsid: true")
	}
}
