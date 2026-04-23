package maintenance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpInvariants(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "invariants-test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create .env
	envContent := "BASE_URL=https://localhost:8443\n"
	err = os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(envContent), 0600)
	require.NoError(t, err)

	// Create .agents/coverage.json
	err = os.MkdirAll(filepath.Join(tmpDir, ".agents"), 0750)
	require.NoError(t, err)
	coverageContent := `{"default": 85.5}`
	err = os.WriteFile(filepath.Join(tmpDir, ".agents", "coverage.json"), []byte(coverageContent), 0600)
	require.NoError(t, err)

	// Create internal/module/admin/listings.go
	err = os.MkdirAll(filepath.Join(tmpDir, "internal/module/admin"), 0750)
	require.NoError(t, err)
	listingsContent := `package admin
const maxFeatured = 5
func someFunc() {
    if len(featured) >= 5 {
    }
}`
	err = os.WriteFile(filepath.Join(tmpDir, "internal/module/admin/listings.go"), []byte(listingsContent), 0600)
	require.NoError(t, err)

	// Run DumpInvariants
	err = DumpInvariants(tmpDir)
	require.NoError(t, err)

	// Verify output
	outputPath := filepath.Join(tmpDir, ".agents", "invariants.json")
	require.FileExists(t, outputPath)

	data, err := os.ReadFile(outputPath) //nolint:gosec // test code
	require.NoError(t, err)

	var inv Invariants
	err = json.Unmarshal(data, &inv)
	require.NoError(t, err)

	assert.Equal(t, "https", inv.Protocol)
	assert.Equal(t, "8443", inv.Port)
	assert.Equal(t, 85.5, inv.DefaultCoverage)
	assert.Equal(t, "sqlite", inv.DBEngine)
	assert.Equal(t, "script-src 'self'", inv.CSPPolicy)
	assert.Equal(t, 5, inv.MaxFeaturedListings)
}

func TestDumpInvariants_Defaults(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "invariants-test-defaults")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Run without files
	err = DumpInvariants(tmpDir)
	require.NoError(t, err)

	outputPath := filepath.Join(tmpDir, ".agents", "invariants.json")
	require.FileExists(t, outputPath)

	data, err := os.ReadFile(outputPath) //nolint:gosec // test code
	require.NoError(t, err)

	var inv Invariants
	err = json.Unmarshal(data, &inv)
	require.NoError(t, err)

	assert.Equal(t, "https", inv.Protocol)
	assert.Equal(t, "8443", inv.Port)
	assert.Equal(t, 0.0, inv.DefaultCoverage)
	assert.Equal(t, "sqlite", inv.DBEngine)
	assert.Equal(t, "script-src 'self'", inv.CSPPolicy)
	assert.Equal(t, 0, inv.MaxFeaturedListings)
}
