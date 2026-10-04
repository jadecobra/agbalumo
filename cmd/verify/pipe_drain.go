package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var pipeDrainCmd = makeSimpleCmd("pipe-drain", "Verify that test files using os.Pipe drain concurrently to prevent buffer deadlocks", func() error {
	fmt.Println("🔍 Checking for synchronous os.Pipe usage without concurrent drain...")
	violations, err := maintenance.CheckPipeDrain(".")
	if err != nil {
		return err
	}

	if len(violations) > 0 {
		fmt.Printf("❌ Found %d pipe drain violations:\n", len(violations))
		for _, v := range violations {
			fmt.Printf("  - %s:%d: %s\n", v.File, v.Line, v.Message)
		}
		return fmt.Errorf("pipe drain check failed")
	}

	fmt.Println("✅ All os.Pipe usages drain concurrently.")
	return nil
})
