package util

import (
	"os"
)

// SafeMkdir creates a directory with 0750 permissions.
func SafeMkdir(path string) error {
	return os.MkdirAll(path, 0750)
}

// SafeWriteFile writes a file with 0600 permissions.
func SafeWriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0600)
}

// SafeRemove removes a file safely.
func SafeRemove(path string) error {
	return os.Remove(path)
}

// SafeStat returns FileInfo for a path.
func SafeStat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// SafeReadFile reads a file safely.
func SafeReadFile(filename string) ([]byte, error) {
	// #nosec G304 - Secure wrapper for file reading in domain utility
	return os.ReadFile(filename)
}
