package util

import (
	"os"
	"testing"
)

func TestSafeMkdir(t *testing.T) {
	path := "testdir"
	defer os.RemoveAll(path)

	err := SafeMkdir(path)
	if err != nil {
		t.Fatalf("SafeMkdir failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat failed: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Expected directory, got file")
	}

	mode := info.Mode().Perm()
	if mode != 0750 {
		t.Errorf("Expected mode 0750, got %v", mode)
	}
}

func TestSafeWriteFile(t *testing.T) {
	filename := "testfile"
	data := []byte("hello world")
	defer os.Remove(filename)

	err := SafeWriteFile(filename, data)
	if err != nil {
		t.Fatalf("SafeWriteFile failed: %v", err)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("os.Stat failed: %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Expected mode 0600, got %v", mode)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile failed: %v", err)
	}

	if string(content) != string(data) {
		t.Errorf("Expected content %s, got %s", string(data), string(content))
	}
}
