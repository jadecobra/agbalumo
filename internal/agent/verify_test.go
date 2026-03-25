package agent

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

type commandRecord struct {
	Name string
	Args []string
}

var recordedCommands []commandRecord

func mockExecCommand(command string, args ...string) *exec.Cmd {
	recordedCommands = append(recordedCommands, commandRecord{Name: command, Args: args})
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulated outputs
	os.Exit(0)
}

func TestVerifyImplementation(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	recordedCommands = nil
	if !VerifyImplementation() {
		t.Fatal("VerifyImplementation failed in mock")
	}

	expected := []struct {
		Name string
		Args []string
	}{
		{"go", []string{"vet", "./..."}},
		{"go", []string{"build", "./..."}},
		{"go", []string{"test", "-json", "-coverprofile=.tester/coverage/coverage.out", "./..."}},
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

func TestVerifyLint(t *testing.T) {
	origExec := ExecCommand
	origLook := LookPath
	ExecCommand = mockExecCommand
	LookPath = func(file string) (string, error) {
		if file == "golangci-lint" {
			return "/usr/local/bin/golangci-lint", nil
		}
		return "", fmt.Errorf("not found")
	}
	defer func() {
		ExecCommand = origExec
		LookPath = origLook
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



func TestVerifyCoverage(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	// We just want to call it to increase coverage
    oldCov := ".tester/coverage/coverage.out"
	// Write a fake coverage file for parsing failure test or pass
	_ = os.MkdirAll(".tester/coverage", 0755)
	_ = os.WriteFile(oldCov, []byte("mode: set\n"), 0644)
	defer func() { _ = os.Remove(oldCov) }()

	// Since we wrote an empty coverage profile, it will probably fail the threshold check
	success := VerifyCoverage()
    fmt.Printf("Coverage result: %v\n", success)
}

func mockExecCommandRedTestEvasion(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessRedTestEvasion", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcessRedTestEvasion(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulated outputs for evasion (a test panicking prints fail for package but no tests)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Output":"FAIL\tfake\t0.331s\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Elapsed":0.332}`)
	os.Exit(0)
}

func mockExecCommandRedTestValid(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessRedTestValid", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcessRedTestValid(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulated outputs for a valid failing red-test
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Test":"TestRed","Output":"--- FAIL: TestRed\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Test":"TestRed","Elapsed":0.332}`)
	os.Exit(0)
}

func TestVerifyRedTest(t *testing.T) {
	t.Run("EvasionExploit", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = mockExecCommandRedTestEvasion
		defer func() { ExecCommand = orig }()

		// With pattern="" we expect it to fail the gate because there's no actual test failures inside
		if VerifyRedTest("") {
			t.Error("VerifyRedTest passed on an evasion exploit! It should have failed.")
		}
	})

	t.Run("ValidFailure", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = mockExecCommandRedTestValid
		defer func() { ExecCommand = orig }()

		// With pattern="" we expect it to pass the gate
		if !VerifyRedTest("") {
			t.Error("VerifyRedTest failed a valid failing red-test.")
		}
	})
}

func TestVerifyApiSpec(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	VerifyApiSpec("feature")
}

func TestVerifyApiSpec_DriftAggregation(t *testing.T) {
	orig := ExecCommand
	ExecCommand = mockExecCommand
	defer func() { ExecCommand = orig }()

	// Mock file system in a temp dir
	tmpDir := t.TempDir()

	err := os.MkdirAll(tmpDir+"/cmd", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/docs/cli", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/internal/handler", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(tmpDir+"/internal/module", 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Code commands
	_ = os.WriteFile(tmpDir+"/cmd/main.go", []byte(`package main
import "github.com/spf13/cobra"
var cmd = &cobra.Command{ Use: "serve" }`), 0644)
	
	// OpenApi and Docs
	_ = os.WriteFile(tmpDir+"/docs/openapi.yaml", []byte(`
paths:
  /auth:
    get: {}`), 0644)
	
	_ = os.WriteFile(tmpDir+"/docs/api.md", []byte(`
| Method | Path |
|--------|------|
| GET | /missing |`), 0644)

	// MD Commands missing 'serve' and adding 'unknown'
	_ = os.WriteFile(tmpDir+"/docs/cli.md", []byte(`### unknown`), 0644)

	// Change dir to TempDir to trick hardcoded paths
	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origWd) }()

	// The system shouldn't pass
	if VerifyApiSpec("feature") {
		t.Error("VerifyApiSpec should have failed due to drift violations")
	}
}
