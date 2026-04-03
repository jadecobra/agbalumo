package agent

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// VerifyVibeCheck parses the vibe_check.md from the task directory
// and ensures all [ ] items have been converted to [x].
func VerifyVibeCheck() bool {
	vibeFile := filepath.Join(".tester", "tasks", "vibe_check.md")

	// If it doesn't exist, we can't pass.
	// #nosec G304 - Internal harness logic reading project files
	file, err := os.Open(vibeFile)
	if err != nil {
		fmt.Printf("❌ Gate FAIL: %s not found. You must create this file from the template.\n", vibeFile)
		return false
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	allPassed := true
	foundCheckboxes := false

	// Regex to find [ ] or [/] - we only want to see [x] or [X]
	incompleteRegex := regexp.MustCompile(`\[\s?[\/\s]\s?\]`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check if it's a checkbox line
		if strings.Contains(line, "[") && strings.Contains(line, "]") {
			foundCheckboxes = true
			if incompleteRegex.MatchString(line) {
				fmt.Printf("❌ Pending Item (Line %d): %s\n", lineNum, strings.TrimSpace(line))
				allPassed = false
			}
		}
	}

	if !foundCheckboxes {
		fmt.Printf("❌ Gate FAIL: No checkboxes found in %s. Check the template.\n", vibeFile)
		return false
	}

	if allPassed {
		fmt.Printf("✅ Gate PASS: All vibe-check items manual verified by human.\n")
	} else {
		fmt.Printf("❌ Gate FAIL: One or more vibe-check items are pending user review.\n")
	}

	return allPassed
}
