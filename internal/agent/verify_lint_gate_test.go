package agent

import (
	"fmt"
	"os"
	"testing"
)

func TestVerifyLint(t *testing.T) {
	origExec := ExecCommand
	origLook := LookPath
	origStat := OSStat
	ExecCommand = mockExecCommand
	LookPath = func(file string) (string, error) {
		if file == "golangci-lint" {
			return "golangci-lint", nil
		}
		return "", fmt.Errorf("not found")
	}
	OSStat = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}
	defer func() {
		ExecCommand = origExec
		LookPath = origLook
		OSStat = origStat
	}()

	recordedCommands = nil
	if !VerifyLint() {
		t.Fatal("VerifyLint failed in mock")
	}

	expected := []struct {
		Name string
		Args []string
	}{
		{"golangci-lint", []string{"run", "-c", "scripts/.golangci.yml"}},
	}

	if len(recordedCommands) != len(expected) {
		t.Fatalf("expected %d commands, got %d", len(expected), len(recordedCommands))
	}

	for i, exp := range expected {
		if recordedCommands[i].Name != exp.Name {
			t.Errorf("cmd %d: expected name %s, got %s", i, exp.Name, recordedCommands[i].Name)
		}
		if len(recordedCommands[i].Args) != len(exp.Args) {
			t.Errorf("cmd %d: expected %d args, got %d", i, len(exp.Args), len(recordedCommands[i].Args))
			continue
		}
		for j, arg := range exp.Args {
			if recordedCommands[i].Args[j] != arg {
				t.Errorf("cmd %d, arg %d: expected %s, got %s", i, j, arg, recordedCommands[i].Args[j])
			}
		}
	}
}

func TestVerifyLint_NotFound(t *testing.T) {
	origLook := LookPath
	origStat := OSStat
	LookPath = func(file string) (string, error) { return "", fmt.Errorf("not found") }
	OSStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	defer func() { LookPath = origLook; OSStat = origStat }()

	if !VerifyLint() {
		t.Error("VerifyLint should return true (pass skipped) when golangci-lint not found")
	}
}
