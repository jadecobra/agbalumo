package commands

import (
	"github.com/jadecobra/agbalumo/internal/agent"
	"github.com/jadecobra/agbalumo/internal/util"
	"os/exec"
	"testing"
)

func TestVerifyCmd_Errors(t *testing.T) {
	cmd := VerifyCmd()

	// 1. Test without active feature
	stateFile := ".agents/state.json"
	backupFile := stateFile + ".bak"
	if _, err := util.SafeStat(stateFile); err == nil {
		_ = util.SafeRename(stateFile, backupFile)
		defer func() { _ = util.SafeRename(backupFile, stateFile) }()
	}

	err := cmd.RunE(cmd, []string{"red-test"})
	if err == nil {
		t.Errorf("expected error when no active feature exists, got nil")
	}

	// 2. Test with invalid gate
	state := &agent.State{
		Feature: "test",
		Phase:   "RED",
	}
	_ = agent.SaveState(stateFile, state)
	defer func() { _ = util.SafeRemove(stateFile) }()

	err = cmd.RunE(cmd, []string{"invalid-gate"})
	if err == nil {
		t.Errorf("expected error for invalid gate_id, got nil")
	}

	// 3. Test dependency failure (Implementation requires red-test PASS)
	state.Gates.RedTest = agent.GatePending
	_ = agent.SaveState(stateFile, state)
	err = cmd.RunE(cmd, []string{"implementation"})
	if err == nil || err.Error() == "" {
		t.Errorf("expected dependency error for implementation, got %v", err)
	}

	// 4. Test all gates with mocks
	oldExec := agent.ExecCommand
	defer func() { agent.ExecCommand = oldExec }()

	agent.ExecCommand = func(name string, args ...string) *exec.Cmd {
		// Mock success for all commands
		return exec.Command("echo", "success")
	}

	state.Gates.RedTest = agent.GatePassed
	state.Gates.ApiSpec = agent.GatePassed
	state.Gates.Implementation = agent.GatePassed
	state.Phase = "REFACTOR"
	_ = agent.SaveState(stateFile, state)

	gates := []string{"red-test", "api-spec", "implementation", "lint", "coverage", "template-drift", "security-static", "browser-verification"}

	for _, g := range gates {
		t.Run("Gate_"+g, func(t *testing.T) {
			_ = cmd.RunE(cmd, []string{g})
		})
	}

	// 5. Test security-static with pattern
	t.Run("Gate_security-static_pattern", func(t *testing.T) {
		_ = cmd.RunE(cmd, []string{"security-static", "cmd internal"})
	})

	// Test browser-verification instructions
	state.Gates.BrowserVerification = agent.GatePending
	_ = agent.SaveState(stateFile, state)
	_ = cmd.RunE(cmd, []string{"browser-verification"})
}
