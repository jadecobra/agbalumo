package maintenance

import (
	"fmt"
	"os"
	"os/exec"
)

// DumpSQLiteSchema uses sqlite3 to dump the schema of the provided database path.
func DumpSQLiteSchema(dbPath string) (string, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("database file does not exist: %s", dbPath)
	}

	// #nosec G204
	cmd := exec.Command("sqlite3", dbPath, ".schema")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to dump schema: %w (output: %s)", err, string(out))
	}

	return string(out), nil
}

