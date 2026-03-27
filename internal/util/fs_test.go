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

func TestSafeReadFile(t *testing.T) {
	filename := "testreadfile"
	data := []byte("safe read content")
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	defer os.Remove(filename)

	content, err := SafeReadFile(filename)
	if err != nil {
		t.Fatalf("SafeReadFile failed: %v", err)
	}

	if string(content) != string(data) {
		t.Errorf("Expected content %s, got %s", string(data), string(content))
	}
}

func TestSafeRemove(t *testing.T) {
	filename := "testremove"
	data := []byte("remove content")
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}

	err = SafeRemove(filename)
	if err != nil {
		t.Fatalf("SafeRemove failed: %v", err)
	}

	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Errorf("Expected file to be removed, got error: %v", err)
	}
}

func TestSafeStat(t *testing.T) {
	filename := "teststat"
	data := []byte("stat content")
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	defer os.Remove(filename)

	info, err := SafeStat(filename)
	if err != nil {
		t.Fatalf("SafeStat failed: %v", err)
	}

	if info.Name() != filename {
		t.Errorf("Expected filename %s, got %s", filename, info.Name())
	}
}
