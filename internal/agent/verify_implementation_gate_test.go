package agent

import (
	"os"
	"os/exec"
	"testing"
)

func TestVerifyImplementation(t *testing.T) {
	origExec := ExecCommand
	origLook := LookPath
	origStat := OSStat
	ExecCommand = mockExecCommand
	LookPath = func(file string) (string, error) {
		return "golangci-lint", nil
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
	if !VerifyImplementation() {
		t.Fatal("VerifyImplementation failed in mock")
	}

	expected := []struct {
		Name string
		Args []string
	}{
		{"golangci-lint", []string{"run", "-c", "scripts/.golangci.yml"}},
		{"go", []string{"build", "-buildvcs=false", "./cmd/...", "./internal/..."}},
		{"go", []string{"test", "-buildvcs=false", "-json", "-coverprofile=.tester/coverage/coverage.out", "./..."}},
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

func TestVerifyImplementation_BuildError(t *testing.T) {
	orig := ExecCommand
	ExecCommand = makeMockExec("TestHelperProcessCompileFail")
	defer func() { ExecCommand = orig }()

	if VerifyImplementation() {
		t.Error("VerifyImplementation should return false on build error")
	}
}

func TestVerifyImplementation_Success(t *testing.T) {
	orig := ExecCommand
	// Mock all commands to succeed
	ExecCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	defer func() { ExecCommand = orig }()

	// We need to pass the lint and api-spec check inside VerifyImplementation
	if !VerifyImplementation() {
		t.Error("VerifyImplementation should return true on success")
	}
}
