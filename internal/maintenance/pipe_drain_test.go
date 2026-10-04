package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

type pipeDrainTestCase struct {
	name          string
	content       string
	wantViolation bool
}

func getPipeDrainCases() []pipeDrainTestCase {
	return []pipeDrainTestCase{
		{
			name: "synchronous drain after close blocks on large output",
			content: `package sample_test

import (
	"io"
	"os"
	"testing"
)

func TestSyncDrain(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_ = doSomething()

	_ = w.Close()
	os.Stdout = old
	_, _ = io.ReadAll(r)
}
`,
			wantViolation: true,
		},
		{
			name: "concurrent drain in goroutine is safe",
			content: `package sample_test

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestConcurrentDrain(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	_ = doSomething()

	_ = w.Close()
	os.Stdout = old
	_ = <-outChan
}
`,
			wantViolation: false,
		},
		{
			name: "concurrent drain via anonymous goroutine reading pipe reader",
			content: `package sample_test

import (
	"io"
	"os"
	"testing"
)

func TestGoroutineRead(t *testing.T) {
	r, w, _ := os.Pipe()
	go func() {
		_, _ = io.ReadAll(r)
	}()
	_ = w.Close()
}
`,
			wantViolation: false,
		},
		{
			name: "test file without os.Pipe has no violations",
			content: `package sample_test

import "testing"

func TestNormal(t *testing.T) {
	t.Log("plain test")
}
`,
			wantViolation: false,
		},
		{
			name: "os.Pipe inside helper function called synchronously without goroutine",
			content: `package sample_test

import (
	"io"
	"os"
	"testing"
)

func capture(fn func()) string {
	r, w, _ := os.Pipe()
	fn()
	_ = w.Close()
	data, _ := io.ReadAll(r)
	return string(data)
}
`,
			wantViolation: true,
		},
	}
}

func runPipeDrainTestCase(t *testing.T, tc pipeDrainTestCase) {
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "sample_test.go")
	if err := os.WriteFile(testFilePath, []byte(tc.content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	violations, err := CheckPipeDrain(tmpDir)
	if err != nil {
		t.Fatalf("CheckPipeDrain failed unexpectedly: %v", err)
	}

	hasViolation := len(violations) > 0
	if hasViolation != tc.wantViolation {
		t.Errorf("got violation = %v (%d violations), wantViolation = %v", hasViolation, len(violations), tc.wantViolation)
		for _, v := range violations {
			t.Logf("violation: %s:%d: %s", v.File, v.Line, v.Message)
		}
	}
}

func TestCheckPipeDrain(t *testing.T) {
	for _, tc := range getPipeDrainCases() {
		t.Run(tc.name, func(t *testing.T) {
			runPipeDrainTestCase(t, tc)
		})
	}
}
