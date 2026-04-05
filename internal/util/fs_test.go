package util

import (
	"os"
	"testing"
)

func TestSafeMkdir(t *testing.T) {
	path := "testdir"
	_ = os.RemoveAll(path)

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
	if mode != 0755 {
		t.Errorf("Expected mode 0755, got %v", mode)
	}
}

func TestSafeWriteFile(t *testing.T) {
	filename := "testfile"
	data := []byte("hello world")
	_ = os.Remove(filename)

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

	content, err := os.ReadFile(/*nolint:gosec*/ filename)
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
	err := os.WriteFile(/*nolint:gosec*/ filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	_ = os.Remove(filename)

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
	err := os.WriteFile(/*nolint:gosec*/ filename, data, 0600)
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
	err := os.WriteFile(/*nolint:gosec*/ filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	_ = os.Remove(filename)

	info, err := SafeStat(filename)
	if err != nil {
		t.Fatalf("SafeStat failed: %v", err)
	}

	if info.Name() != filename {
		t.Errorf("Expected filename %s, got %s", filename, info.Name())
	}
}

func TestSafeIsNotExist(t *testing.T) {
	t.Run("ExistingFile", func(t *testing.T) {
		filename := "existing_test"
		err := os.WriteFile(/*nolint:gosec*/ filename, []byte("data"), 0600)
		if err != nil {
			t.Fatalf("os.WriteFile failed: %v", err)
		}
		_ = os.Remove(filename)

		_, err = os.Stat(filename)
		if SafeIsNotExist(err) {
			t.Error("Expected SafeIsNotExist to be false for existing file")
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		filename := "non_existent_test"
		_, err := os.Stat(filename)
		if !SafeIsNotExist(err) {
			t.Error("Expected SafeIsNotExist to be true for non-existent file")
		}
	})

	t.Run("NilError", func(t *testing.T) {
		if SafeIsNotExist(nil) {
			t.Error("Expected SafeIsNotExist to be false for nil error")
		}
	})
}

func TestSafeRename(t *testing.T) {
	oldFile := "rename_old"
	newFile := "rename_new"
	data := []byte("rename test")

	err := os.WriteFile(/*nolint:gosec*/ oldFile, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	defer func() {
		_ = os.Remove(oldFile)
		_ = os.Remove(newFile)
	}()

	err = SafeRename(oldFile, newFile)
	if err != nil {
		t.Fatalf("SafeRename failed: %v", err)
	}

	if _, err := os.Stat(newFile); err != nil {
		t.Errorf("Expected new file to exist, got error: %v", err)
	}
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Errorf("Expected old file to be gone, got error: %v", err)
	}
}

func TestSafeOpen(t *testing.T) {
	filename := "open_test"
	data := []byte("open test content")

	err := os.WriteFile(/*nolint:gosec*/ filename, data, 0600)
	if err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}
	_ = os.Remove(filename)

	f, err := SafeOpen(filename)
	if err != nil {
		t.Fatalf("SafeOpen failed: %v", err)
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("f.Stat failed: %v", err)
	}
	if info.Name() != filename {
		t.Errorf("Expected filename %s, got %s", filename, info.Name())
	}
}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "EmptyInput",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "NilInput",
			input:    nil,
			expected: nil,
		},
		{
			name:     "NoDuplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "WithDuplicates",
			input:    []string{"c", "a", "b", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqueStrings(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d", len(tt.expected), len(got))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("At index %d: expected %s, got %s", i, tt.expected[i], got[i])
				}
			}
		})
	}
}
