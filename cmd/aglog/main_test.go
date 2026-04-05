package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/history"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Setup: use a temporary directory for history files
	tmpDir := t.TempDir()
	oldDir := history.DefaultStorageDir
	history.DefaultStorageDir = tmpDir
	defer func() { history.DefaultStorageDir = oldDir }()

	t.Run("From Flags", func(t *testing.T) {
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"--feature", "FlagFeature", "--arch", "Arch", "--po", "PO", "--sdet", "SDET", "--be", "BE", "--summary", "Flag Summary"})
		err := cmd.Execute()
		assert.NoError(t, err)

		output := strings.TrimSpace(out.String())
		assert.Contains(t, output, "FlagFeature")
		assert.Contains(t, output, tmpDir)
	})

	t.Run("From STDIN", func(t *testing.T) {
		jsonInput := `{"FeatureName": "JSONFeature", "SystemsArchitect": "Arch", "ProductOwner": "PO", "SDET": "SDET", "BackendEngineer": "BE", "DecisionSummary": "JSON Summary"}`
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetIn(strings.NewReader(jsonInput))
		err := cmd.Execute()
		assert.NoError(t, err)

		output := strings.TrimSpace(out.String())
		assert.Contains(t, output, "JSONFeature")
		assert.Contains(t, output, tmpDir)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetIn(strings.NewReader("invalid"))
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("Empty Input", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetIn(strings.NewReader(""))
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FeatureName is required")
	})

	t.Run("Flag Parse Error", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"--invalid-flag"})
		// Cobra prints error to stderr by default unless SilenceErrors is true
		cmd.SilenceErrors = true
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("Store Error", func(t *testing.T) {
		oldDir := history.DefaultStorageDir
		history.DefaultStorageDir = "/non-existent/dir"
		defer func() { history.DefaultStorageDir = oldDir }()

		cmd := NewRootCmd()
		cmd.SetArgs([]string{"--feature", "ErrFeature"})
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to store decision")
	})

	t.Run("Stdout Write Error", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetOut(&errorWriter{})
		cmd.SetArgs([]string{"--feature", "WriteErr"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("Stdin Read Error", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetIn(&errorReader{})
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read from stdin")
	})

	t.Run("Partial Flags", func(t *testing.T) {
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"--feature", "Partial", "--be", "BE"})
		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, out.String(), "Partial")
	})
}

func TestMain_Success(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir := history.DefaultStorageDir
	history.DefaultStorageDir = tmpDir
	defer func() { history.DefaultStorageDir = oldDir }()

	oldArgs := os.Args
	os.Args = []string{"aglog", "--feature", "MainCall"}
	defer func() { os.Args = oldArgs }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "MainCall")
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}
