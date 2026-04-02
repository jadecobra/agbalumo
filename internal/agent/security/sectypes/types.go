package sectypes

import (
	"fmt"
)

// SecurityViolation represents a potential security issue found in the code.
type SecurityViolation struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
}

func (v SecurityViolation) String() string {
	return fmt.Sprintf("%s:%d:%d: [%s] %s", v.File, v.Line, v.Column, v.Type, v.Message)
}

func DeduplicateViolations(violations []SecurityViolation) []SecurityViolation {
	seen := make(map[string]SecurityViolation)
	var lines []string
	for _, v := range violations {
		key := fmt.Sprintf("%s:%d", v.File, v.Line)
		existing, ok := seen[key]
		if !ok {
			seen[key] = v
			lines = append(lines, key)
		} else if existing.Type == "Entropy" && v.Type == "Secret" {
			// Prioritize "Secret" over "Entropy" for the same line
			seen[key] = v
		}
	}
	var unique []SecurityViolation
	for _, key := range lines {
		unique = append(unique, seen[key])
	}
	return unique
}
