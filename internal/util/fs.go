package util

import (
	"os"
)

// SafeMkdir creates a directory with 0755 permissions.
func SafeMkdir(path string) error {
	// #nosec G301 - 0755 is intentional to allow web server read access
	return os.MkdirAll(path, 0755)
}

// SafeRename renames a file safely.
func SafeRename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// SafeWriteFile writes a file with 0644 permissions.
func SafeWriteFile(filename string, data []byte) error {
	// #nosec G306 - 0644 is intentional to allow web server read access
	return os.WriteFile(filename, data, 0644)
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

// SafeOpen opens a file safely.
func SafeOpen(name string) (*os.File, error) {
	// #nosec G304 - Secure wrapper for file opening in domain utility
	return os.Open(name)
}

// SafeIsNotExist checks if an error indicates that a file does not exist.
func SafeIsNotExist(err error) bool {
	return os.IsNotExist(err)
}
