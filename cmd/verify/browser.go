package main

import (
	"fmt"
)

var browserCmd = makeSimpleCmd("browser", "Execute Playwright end-to-end UI verification tests", func() error {
	fmt.Println("🎭 Running Playwright E2E tests...")
	return runCmd("npm", "run", "test:e2e")
})
