package agent

import (
	"os"
	"testing"
)

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
