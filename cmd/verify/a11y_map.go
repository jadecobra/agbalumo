package main

import (
	"fmt"

	"github.com/jadecobra/agbalumo/internal/maintenance"
)

var a11yMapCmd = makeSimpleCmd("a11y-map", "Parse latest Playwright a11y test-results and map violations to template files", func() error {
	violations, err := maintenance.MapA11yViolations(".")
	if err != nil {
		return err
	}

	if len(violations) == 0 {
		fmt.Println("✅ No a11y violations found in latest test-results.")
		return nil
	}

	fmt.Println("| Violation | Impact | File | Line | Fix |")
	fmt.Println("|-----------|--------|------|------|-----|")
	for _, v := range violations {
		fmt.Printf("| %s | %s | %s | %d | %s |\n",
			v.ViolationID,
			v.Impact,
			v.TemplateFile,
			v.Line,
			v.FixSuggestion,
		)
	}

	return nil
})

func init() {
	// a11yMapCmd is registered in main.go init() to keep it central
}
