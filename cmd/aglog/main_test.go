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
		args := []string{
			"aglog", // binary name usually first arg
			"-feature", "FlagFeature",
			"-arch", "Arch",
			"-po", "PO",
			"-sdet", "SDET",
			"-be", "BE",
			"-summary", "Flag Summary",
		}
		var out bytes.Buffer
		err := run(args, nil, &out)
		assert.NoError(t, err)
		
		output := strings.TrimSpace(out.String())
		assert.Contains(t, output, "FlagFeature")
		assert.Contains(t, output, tmpDir)
	})

	t.Run("From STDIN", func(t *testing.T) {
		jsonInput := `{"FeatureName": "JSONFeature", "SystemsArchitect": "Arch", "ProductOwner": "PO", "SDET": "SDET", "BackendEngineer": "BE", "DecisionSummary": "JSON Summary"}`
		var out bytes.Buffer
		err := run([]string{"aglog"}, strings.NewReader(jsonInput), &out)
		assert.NoError(t, err)
		
		output := strings.TrimSpace(out.String())
		assert.Contains(t, output, "JSONFeature")
		assert.Contains(t, output, tmpDir)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{"aglog"}, strings.NewReader("invalid"), &out)
		assert.Error(t, err)
	})

	t.Run("Empty Input", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{"aglog"}, strings.NewReader(""), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FeatureName is required")
	})

	t.Run("Flag Parse Error", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{"aglog", "-invalid-flag"}, nil, &out)
		assert.Error(t, err)
	})

	t.Run("Store Error", func(t *testing.T) {
		oldDir := history.DefaultStorageDir
		history.DefaultStorageDir = "/non-existent/dir"
		defer func() { history.DefaultStorageDir = oldDir }()

		args := []string{"aglog", "-feature", "ErrFeature"}
		var out bytes.Buffer
		err := run(args, nil, &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to store decision")
	})

	t.Run("Stdout Write Error", func(t *testing.T) {
		args := []string{"aglog", "-feature", "WriteErr"}
		err := run(args, nil, &errorWriter{})
		assert.Error(t, err)
	})

	t.Run("Stdin Read Error", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{"aglog"}, &errorReader{}, &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read from stdin")
	})

	t.Run("Partial Flags", func(t *testing.T) {
		args := []string{"aglog", "-feature", "Partial", "-be", "BE"}
		var out bytes.Buffer
		err := run(args, nil, &out)
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
	os.Args = []string{"aglog", "-feature", "MainCall"}
	defer func() { os.Args = oldArgs }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
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
