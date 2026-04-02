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
