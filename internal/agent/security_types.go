package agent

import (
	"github.com/jadecobra/agbalumo/internal/agent/security/sectypes"
)

// SecurityViolation represents a potential security issue found in the code.
type SecurityViolation = sectypes.SecurityViolation

func deduplicateViolations(violations []SecurityViolation) []SecurityViolation {
	return sectypes.DeduplicateViolations(violations)
}
