package maintenance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LessonsViolation represents a violation of lessons conformance.
type LessonsViolation struct {
	File    string
	Message string
	Line    int
}

// CheckLessonsConformance verifies the integrity of the strict lessons.
func CheckLessonsConformance(rootDir string) ([]LessonsViolation, error) {
	standardsPath := filepath.Join(rootDir, ".agents", "workflows", "coding-standards.md")

	// If the file doesn't exist, return no violations (graceful fallback)
	if _, err := os.Stat(standardsPath); os.IsNotExist(err) {
		return nil, nil
	}

	// #nosec G304 - rootDir is trusted within the context of the maintenance tool chain
	file, err := os.Open(standardsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open coding-standards.md: %w", err)
	}
	defer func() { _ = file.Close() }()

	activeCount, violations, err := parseLessonsFile(file, standardsPath)
	if err != nil {
		return nil, err
	}

	// Enforce 20-lesson ceiling
	if activeCount > 20 {
		violations = append(violations, LessonsViolation{
			File:    standardsPath,
			Line:    1,
			Message: fmt.Sprintf("active prose strict lessons count (%d) exceeds strict lesson ceiling of 20", activeCount),
		})
	}

	return violations, nil
}

// parseLessonsFile parses the lines of coding-standards.md and returns the count of active strict lessons and any violations.
func parseLessonsFile(file *os.File, filePath string) (int, []LessonsViolation, error) {
	var violations []LessonsViolation
	scanner := bufio.NewScanner(file)
	lineNum := 0
	inStrictLessons := false
	activeLessonsCount := 0
	listItemRegex := regexp.MustCompile(`^\s*([\*\-\+]|\d+\.)\s+(.*)$`)

	for scanner.Scan() {
		lineNum++
		trimmed := strings.TrimSpace(scanner.Text())
		var isLesson bool
		var violation *LessonsViolation

		inStrictLessons, isLesson, violation = processLine(trimmed, inStrictLessons, filePath, lineNum, listItemRegex)
		if isLesson {
			activeLessonsCount++
			if violation != nil {
				violations = append(violations, *violation)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, nil, fmt.Errorf("error scanning coding-standards.md: %w", err)
	}

	return activeLessonsCount, violations, nil
}

// processLine processes a single line and returns the updated section state, whether it's a lesson, and any corresponding violation.
func processLine(trimmed string, inStrictLessons bool, filePath string, lineNum int, regex *regexp.Regexp) (bool, bool, *LessonsViolation) {
	if strings.HasPrefix(trimmed, "#") {
		return updateSectionState(trimmed, inStrictLessons), false, nil
	}
	if inStrictLessons && trimmed != "" {
		isLesson, v := checkLessonsListItem(trimmed, filePath, lineNum, regex)
		return inStrictLessons, isLesson, v
	}
	return inStrictLessons, false, nil
}

// checkLessonsListItem checks if a list item is a valid strict lesson with a trigger tag.
func checkLessonsListItem(trimmed string, filePath string, lineNum int, regex *regexp.Regexp) (bool, *LessonsViolation) {
	if !regex.MatchString(trimmed) {
		return false, nil
	}
	if !strings.Contains(trimmed, "[TRIGGER:") {
		return true, &LessonsViolation{
			File:    filePath,
			Line:    lineNum,
			Message: "strict lesson missing [TRIGGER: ...] annotation",
		}
	}
	return true, nil
}

// updateSectionState handles entry and exit of the Strict Lessons section.
func updateSectionState(trimmed string, inStrictLessons bool) bool {
	headingText := strings.TrimLeft(trimmed, "# ")
	if strings.Contains(strings.ToLower(headingText), "strict lessons") {
		return true
	} else if inStrictLessons {
		// Only exit if the heading is a level 1 or level 2 heading (starts with "# " or "## ")
		// sub-headings like "### " (level 3 or deeper) do not exit the section
		if strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ") {
			return false
		}
	}
	return inStrictLessons
}
