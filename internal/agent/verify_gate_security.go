package agent

import (
	"fmt"
)

func VerifySecurityStaticGate(paths ...string) bool {
	fmt.Println("Running AST-based structural security checks...")
	var allViolations []SecurityViolation
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, p := range paths {
		violations, err := VerifySecurityStatic(p)
		if err != nil {
			fmt.Printf("❌ Error running security checks on %s: %v\n", p, err)
			return false
		}
		allViolations = append(allViolations, violations...)
	}

	if len(allViolations) == 0 {
		fmt.Println("✅ Gate PASS: no structural security violations found.")
		return true
	}

	limit := 5
	fmt.Printf("❌ Gate FAIL: %d security violations detected (showing first %d).\n", len(allViolations), limit)
	for i, v := range allViolations {
		if i >= limit {
			break
		}
		fmt.Printf("  [%s] %s:%d:%d: %s\n", v.Type, v.File, v.Line, v.Column, v.Message)
	}
	return false
}
