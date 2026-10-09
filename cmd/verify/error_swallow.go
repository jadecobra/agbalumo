package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var errorSwallowCmd = makeSimpleCmd("error-swallow", "Verify that error return values are not swallowed with blank identifier in internal/ packages", func() error {
	fmt.Println("🔍 Checking for swallowed error returns in internal/repository and internal/service...")
	violations, err := maintenance.CheckErrorSwallow(".")
	if err != nil {
		return err
	}

	if len(violations) > 0 {
		fmt.Printf("❌ Found %d swallowed error violations:\n", len(violations))
		for _, v := range violations {
			fmt.Printf("  - %s:%d: %s\n", v.File, v.Line, v.Message)
		}
		return fmt.Errorf("error swallow check failed")
	}

	fmt.Println("✅ No swallowed error returns found.")
	return nil
})
