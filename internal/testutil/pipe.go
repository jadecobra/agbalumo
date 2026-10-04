package testutil

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// CaptureStdout intercepts os.Stdout while executing fn, draining output
// concurrently in a background goroutine to avoid deadlocks when output
// exceeds the 64KB OS pipe buffer.
func CaptureStdout(t *testing.T, fn func() error) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	outChan := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		_ = r.Close()
		outChan <- buf.String()
	}()

	errFn := fn()

	_ = w.Close()
	os.Stdout = old

	output := <-outChan

	if errFn != nil {
		t.Fatalf("captured function failed: %v", errFn)
	}

	return output
}
