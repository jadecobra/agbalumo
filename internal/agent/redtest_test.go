package agent_test

import (
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/agent"
)

func TestParseTestJSON_Success(t *testing.T) {
	jsonLog := `{"Time":"2023-10-26T10:00:00Z","Action":"run","Package":"dummy","Test":"TestPass"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Test":"TestPass","Output":"=== RUN   TestPass\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Test":"TestPass","Output":"--- PASS: TestPass (0.00s)\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"pass","Package":"dummy","Test":"TestPass","Elapsed":0}
{"Time":"2023-10-26T10:00:00Z","Action":"pass","Package":"dummy","Elapsed":0.001}`

	res, err := agent.ParseTestJSON(strings.NewReader(jsonLog))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !res.Success {
		t.Errorf("expected Success to be true")
	}
	if res.CompilationFailed {
		t.Errorf("expected CompilationFailed to be false")
	}
	if len(res.Failures) != 0 {
		t.Errorf("expected 0 failures, got %d", len(res.Failures))
	}
}

func TestParseTestJSON_AssertionFailure(t *testing.T) {
	jsonLog := `{"Time":"2023-10-26T10:00:00Z","Action":"run","Package":"dummy","Test":"TestFail"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Test":"TestFail","Output":"=== RUN   TestFail\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Test":"TestFail","Output":"    dummy_test.go:4: fail\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Test":"TestFail","Output":"--- FAIL: TestFail (0.00s)\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"fail","Package":"dummy","Test":"TestFail","Elapsed":0}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Output":"FAIL\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"fail","Package":"dummy","Elapsed":0.001}`

	res, err := agent.ParseTestJSON(strings.NewReader(jsonLog))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Success {
		t.Errorf("expected Success to be false")
	}
	if res.CompilationFailed {
		t.Errorf("expected CompilationFailed to be false")
	}
	if len(res.Failures) != 1 {
		t.Fatalf("expected 1 failure, got %d", len(res.Failures))
	}

	failure := res.Failures[0]
	if failure.TestName != "TestFail" {
		t.Errorf("expected test name 'TestFail', got %q", failure.TestName)
	}
	expectedOutput := "    dummy_test.go:4: fail\n--- FAIL: TestFail (0.00s)\n"
	if !strings.Contains(failure.Output, expectedOutput) {
		t.Errorf("expected output to contain %q, got %q", expectedOutput, failure.Output)
	}
}

func TestParseTestJSON_CompilationFailure(t *testing.T) {
	jsonLog := `{"ImportPath":"dummy.test","Action":"build-output","Output":"# dummy\n"}
{"ImportPath":"dummy.test","Action":"build-output","Output":"dummy2_test.go:4:10: expected ';', found code\n"}
{"ImportPath":"dummy.test","Action":"build-fail"}
{"Time":"2023-10-26T10:00:00Z","Action":"start","Package":"dummy"}
{"Time":"2023-10-26T10:00:00Z","Action":"output","Package":"dummy","Output":"FAIL\tdummy [setup failed]\n"}
{"Time":"2023-10-26T10:00:00Z","Action":"fail","Package":"dummy","Elapsed":0,"FailedBuild":"dummy.test"}`

	res, err := agent.ParseTestJSON(strings.NewReader(jsonLog))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Success {
		t.Errorf("expected Success to be false")
	}
	if !res.CompilationFailed {
		t.Errorf("expected CompilationFailed to be true")
	}
	// Note: We might also want to capture the build output, but our TestFailure struct
	// only currently has TestName and Output. Let's see if we can shove build errors into it.
	if len(res.Failures) != 1 {
		t.Fatalf("expected 1 failure, got %d", len(res.Failures))
	}
	failure := res.Failures[0]
	if failure.TestName != "build" {
		t.Errorf("expected test name 'build', got %q", failure.TestName)
	}
	expectedOutput := "# dummy\ndummy2_test.go:4:10: expected ';', found code\n"
	if failure.Output != expectedOutput {
		t.Errorf("expected build output %q, got %q", expectedOutput, failure.Output)
	}
}

func TestParseTestJSON_SetupFailed(t *testing.T) {
	// Another variant of compilation failure, sometimes just outputting setup failed
	// without explicit build-fail if it's a test compilation error in an older Go version
	// or specific edge case.
	jsonLog := `{"Time":"2026-03-16T00:52:04.921813-05:00","Action":"start","Package":"dummy"}
{"Time":"2026-03-16T00:52:04.921835-05:00","Action":"output","Package":"dummy","Output":"FAIL\tdummy [setup failed]\n"}
{"Time":"2026-03-16T00:52:04.921838-05:00","Action":"fail","Package":"dummy","Elapsed":0,"FailedBuild":"dummy.test"}`

	res, err := agent.ParseTestJSON(strings.NewReader(jsonLog))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Success {
		t.Errorf("expected Success to be false")
	}
	if !res.CompilationFailed {
		t.Errorf("expected CompilationFailed to be true")
	}
}
