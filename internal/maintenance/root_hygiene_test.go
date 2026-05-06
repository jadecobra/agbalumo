package maintenance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyRootHygiene(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some whitelisted files
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0600))
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "internal"), 0750))

	t.Run("Clean root", func(t *testing.T) {
		err := VerifyRootHygiene(tmpDir)
		assert.NoError(t, err)
	})

	t.Run("Cluttered root with file", func(t *testing.T) {
		clutterFile := filepath.Join(tmpDir, "listings.db")
		require.NoError(t, os.WriteFile(clutterFile, []byte("junk"), 0600))
		defer func() { _ = os.Remove(clutterFile) }()

		err := VerifyRootHygiene(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "File: listings.db")
	})

	t.Run("Cluttered root with directory", func(t *testing.T) {
		clutterDir := filepath.Join(tmpDir, "unexpected_dir")
		require.NoError(t, os.Mkdir(clutterDir, 0750))
		defer func() { _ = os.RemoveAll(clutterDir) }()

		err := VerifyRootHygiene(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Directory: unexpected_dir/")
	})

	t.Run("Hidden files are ignored", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(""), 0600))
		err := VerifyRootHygiene(tmpDir)
		assert.NoError(t, err)
	})
}
