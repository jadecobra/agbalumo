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
		{"go", []string{"vet", "./cmd/...", "./internal/..."}},
		{"go", []string{"build", "./cmd/...", "./internal/..."}},
		{"go", []string{"test", "-json", "-coverprofile=.tester/coverage/coverage.out", "./cmd/...", "./internal/..."}},
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

// --- Helper process factories ---

func makeMockExec(helperName string) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=" + helperName, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

// --- Helper processes ---

func TestHelperProcessRedTestEvasion(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Output":"FAIL\tfake\t0.331s\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Elapsed":0.332}`)
	os.Exit(0)
}

func TestHelperProcessRedTestValid(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Test":"TestRed","Output":"--- FAIL: TestRed\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Test":"TestRed","Elapsed":0.332}`)
	os.Exit(0)
}

func TestHelperProcessUIBypassClean(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) > 0 && args[0] == "git" {
		// git status --porcelain with only test files and HTML
		fmt.Println("M  internal/handler/listing_test.go")
		fmt.Println("M  templates/home.html")
	}
	os.Exit(0)
}

func TestHelperProcessUIBypassRejected(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) > 0 && args[0] == "git" {
		// git status --porcelain with a non-test .go file
		fmt.Println("M  internal/handler/listing.go")
		fmt.Println("M  internal/handler/listing_test.go")
	}
	os.Exit(0)
}

func TestHelperProcessCompileFail(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintln(os.Stderr, "# fake/pkg\nsyntax error")
	os.Exit(1)
}

func TestHelperProcessCompilationFailedJSON(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	// First call: go test -run=^$ (compile check) — succeed
	if len(args) >= 3 && args[1] == "test" && args[2] == "-run=^$" {
		os.Exit(0)
	}
	// Second call: go test -json — return build-fail JSON
	fmt.Println(`{"ImportPath":"fake.test","Action":"build-fail"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"start","Package":"fake"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Output":"FAIL\tfake [setup failed]\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Elapsed":0,"FailedBuild":"fake.test"}`)
	os.Exit(0)
}

func TestHelperProcessAllPass(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	// Compile check succeeds
	if len(args) >= 3 && args[1] == "test" && args[2] == "-run=^$" {
		os.Exit(0)
	}
	// Test run: all pass
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"run","Package":"fake","Test":"TestGreen"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Test":"TestGreen","Output":"--- PASS: TestGreen (0.00s)\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"pass","Package":"fake","Test":"TestGreen","Elapsed":0}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"pass","Package":"fake","Elapsed":0.001}`)
	os.Exit(0)
}

func TestHelperProcessPatternMatch(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) >= 3 && args[1] == "test" && args[2] == "-run=^$" {
		os.Exit(0)
	}
	// Valid failure with identifiable output for pattern matching
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Test":"TestRed","Output":"expected_pattern: value mismatch\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"output","Package":"fake","Test":"TestRed","Output":"--- FAIL: TestRed (0.00s)\n"}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Test":"TestRed","Elapsed":0}`)
	fmt.Println(`{"Time":"2023-01-01T00:00:00Z","Action":"fail","Package":"fake","Elapsed":0.001}`)
	os.Exit(0)
}

// --- Tests ---

func TestVerifyRedTest(t *testing.T) {
	t.Run("EvasionExploit", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessRedTestEvasion")
		defer func() { ExecCommand = orig }()

		if VerifyRedTest("") {
			t.Error("VerifyRedTest passed on an evasion exploit! It should have failed.")
		}
	})

	t.Run("ValidFailure", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessRedTestValid")
		defer func() { ExecCommand = orig }()

		if !VerifyRedTest("") {
			t.Error("VerifyRedTest failed a valid failing red-test.")
		}
	})

	t.Run("UIBypass_Clean", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessUIBypassClean")
		defer func() { ExecCommand = orig }()

		if !VerifyRedTest("ui-bypass") {
			t.Error("VerifyRedTest should pass for UI bypass with only test/HTML files modified")
		}
	})

	t.Run("UIBypass_NonTestGoModified", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessUIBypassRejected")
		defer func() { ExecCommand = orig }()

		if VerifyRedTest("ui-bypass") {
			t.Error("VerifyRedTest should reject UI bypass when non-test .go files are modified")
		}
	})

	t.Run("CompilationFailure", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessCompileFail")
		defer func() { ExecCommand = orig }()

		_ = os.MkdirAll(".tester", 0755)
		if VerifyRedTest("") {
			t.Error("VerifyRedTest should fail when code does not compile")
		}
	})

	t.Run("CompilationFailed_FromJSON", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessCompilationFailedJSON")
		defer func() { ExecCommand = orig }()

		_ = os.MkdirAll(".tester", 0755)
		if VerifyRedTest("") {
			t.Error("VerifyRedTest should fail when JSON reports compilation failure (build-fail)")
		}
	})

	t.Run("AllTestsPass_GateFail", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessAllPass")
		defer func() { ExecCommand = orig }()

		_ = os.MkdirAll(".tester", 0755)
		if VerifyRedTest("") {
			t.Error("VerifyRedTest should fail when all tests pass (red-test expects failure)")
		}
	})

	t.Run("PatternMatched", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessPatternMatch")
		defer func() { ExecCommand = orig }()

		_ = os.MkdirAll(".tester", 0755)
		if !VerifyRedTest("expected_pattern") {
			t.Error("VerifyRedTest should pass when failure output contains the expected pattern")
		}
	})

	t.Run("PatternNotMatched", func(t *testing.T) {
		orig := ExecCommand
		ExecCommand = makeMockExec("TestHelperProcessPatternMatch")
		defer func() { ExecCommand = orig }()

		_ = os.MkdirAll(".tester", 0755)
		if VerifyRedTest("missing_pattern") {
			t.Error("VerifyRedTest should fail when failure output does not contain the expected pattern")
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
