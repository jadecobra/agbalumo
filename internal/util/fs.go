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
