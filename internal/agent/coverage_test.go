package agent

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCoverageProfile(t *testing.T) {
	profile := `mode: set
github.com/jadecobra/agbalumo/internal/handler/auth.go:12.72,14.2 1 1
github.com/jadecobra/agbalumo/internal/handler/auth.go:16.89,18.2 2 0
github.com/jadecobra/agbalumo/internal/handler/auth.go:19.49,21.2 1 1
github.com/jadecobra/agbalumo/internal/domain/user.go:10.10,12.2 4 1
github.com/jadecobra/agbalumo/internal/domain/user.go:15.10,18.2 6 1
`
	
	coverageByPkg, err := ParseCoverageProfile(strings.NewReader(profile))
	require.NoError(t, err)

	assert.Equal(t, 2, len(coverageByPkg))

	// handler package: total = 4 stmts, covered = 2 (1+1). 2/4 = 50%
	handlerPkg := "github.com/jadecobra/agbalumo/internal/handler"
	assert.Contains(t, coverageByPkg, handlerPkg)
	assert.Equal(t, 50.0, coverageByPkg[handlerPkg])

	// domain package: total = 10 stmts, covered = 10. 10/10 = 100%
	domainPkg := "github.com/jadecobra/agbalumo/internal/domain"
	assert.Contains(t, coverageByPkg, domainPkg)
	assert.Equal(t, 100.0, coverageByPkg[domainPkg])
}

func TestEnforceCoverageThresholds(t *testing.T) {
	coverage := map[string]float64{
		"github.com/jadecobra/agbalumo/internal/handler": 85.0,
		"github.com/jadecobra/agbalumo/internal/domain":  100.0,
		"github.com/jadecobra/agbalumo/internal/service": 70.0,
	}

	thresholdJSON := `{
		"github.com/jadecobra/agbalumo/internal/handler": 90.0,
		"github.com/jadecobra/agbalumo/internal/domain": 100.0,
		"default": 80.0
	}`

	thresholds, err := ParseThresholds([]byte(thresholdJSON))
	require.NoError(t, err)

	violations := EnforceCoverage(coverage, thresholds)

	assert.Equal(t, 2, len(violations))

	// Handler failed explicit threshold (85 < 90)
	assert.Equal(t, "github.com/jadecobra/agbalumo/internal/handler: coverage 85.0% is below explicit threshold of 90.0%", violations[0])

	// Service failed default threshold (70 < 80)
	assert.Equal(t, "github.com/jadecobra/agbalumo/internal/service: coverage 70.0% is below default threshold of 80.0%", violations[1])
}
