package testutil

import (
	"fmt"
	"strings"
	"testing"
)

func TestCaptureStdout(t *testing.T) {
	t.Run("captures small output", func(t *testing.T) {
		got := CaptureStdout(t, func() error {
			fmt.Print("hello world")
			return nil
		})
		if got != "hello world" {
			t.Errorf("got %q, want 'hello world'", got)
		}
	})

	t.Run("captures large output exceeding 64KB pipe buffer without deadlock", func(t *testing.T) {
		largeChunk := strings.Repeat("A", 128*1024)
		got := CaptureStdout(t, func() error {
			fmt.Print(largeChunk)
			return nil
		})
		if len(got) != 128*1024 {
			t.Errorf("got length %d, want %d", len(got), 128*1024)
		}
	})
}
