package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jadecobra/agbalumo/internal/maintenance"
	"github.com/spf13/cobra"
)

type rawEvidenceViolation struct {
	Pattern string `json:"pattern"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

var designEvidenceCmd = &cobra.Command{
	Use:   "design-evidence",
	Short: "Output raw design and accessibility violations as pattern IDs (no rubric coupling)",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Retrieve violations from internal maintenance checkers
		designViolations, err := maintenance.CheckDesignStandards("ui/templates")
		if err != nil {
			return fmt.Errorf("failed design standards check: %w", err)
		}

		tokenViolations, err := maintenance.CheckTokenStrictness("ui/templates")
		if err != nil {
			return fmt.Errorf("failed token strictness check: %w", err)
		}

		keyViolations, err := maintenance.CheckTemplateKeyGaps("ui/templates")
		if err != nil {
			return fmt.Errorf("failed template key gaps check: %w", err)
		}

		// Merge all violations
		allViolations := append(designViolations, tokenViolations...)
		allViolations = append(allViolations, keyViolations...)

		// Map to raw evidence violations
		evidence := make([]rawEvidenceViolation, 0, len(allViolations))
		for _, v := range allViolations {
			evidence = append(evidence, rawEvidenceViolation{
				Pattern: mapViolationToPatternID(v),
				File:    v.File,
				Line:    v.Line,
			})
		}

		if jsonOutput {
			report := map[string]interface{}{
				"violations": evidence,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(report)
		}

		// Text format: raw facts only
		for _, e := range evidence {
			fmt.Printf("%s:%d: %s\n", e.File, e.Line, e.Pattern)
		}

		return nil
	},
}

func mapViolationToPatternID(v maintenance.DesignViolation) string {
	r := v.Reason
	switch {
	case strings.Contains(r, "rounding"):
		return "forbidden-rounding"
	case strings.Contains(r, "Hardcoded hex"):
		return "hardcoded-hex"
	case strings.Contains(r, "Font size below 10px"):
		return "font-size-below-min"
	case strings.Contains(r, "opacity below 70%"):
		return "low-contrast-opacity"
	case strings.Contains(r, "dark background bypasses"):
		return "hardcoded-modal-bg"
	case strings.Contains(r, "Ad-hoc color"):
		return "ad-hoc-color"
	case strings.Contains(r, "select-none"):
		return "text-selection-usability"
	case strings.Contains(r, "Blend-mode"):
		return "blend-mode-gradient"
	case strings.Contains(r, "Inline event handler"):
		return "inline-handler"
	case strings.Contains(r, "Inline style"):
		return "inline-style"
	case strings.Contains(r, "missing data-testid"):
		return "missing-data-testid"
	case strings.Contains(r, "Empty section"):
		return "empty-section"
	case strings.Contains(r, "Multiple close actions"):
		return "multiple-close-actions"
	case strings.Contains(r, "missing aria-label"):
		return "icon-button-no-aria-label"
	case strings.Contains(r, "missing 'for' attribute"):
		return "label-no-associated-control"
	case strings.Contains(r, "missing 'id' attribute"):
		return "input-no-id"
	case strings.Contains(r, "missing 'alt' attribute"):
		return "img-no-alt"
	case strings.Contains(r, "missing hx-indicator"):
		return "htmx-no-indicator"
	case strings.Contains(r, "arbitrary Tailwind") || strings.Contains(r, "Tailwind value"):
		return "arbitrary-tailwind-token"
	case strings.Contains(r, "density too high"):
		return "uppercase-density-too-high"
	case strings.Contains(r, "references variable") && strings.Contains(r, "not passed in dict"):
		return "template-key-gap-missing-required"
	case strings.Contains(r, "accepts properties") && strings.Contains(r, "extraneous property"):
		return "template-key-gap-extraneous"
	default:
		return "unknown-violation"
	}
}

func init() {
	designEvidenceCmd.Flags().Bool("json", false, "Output violations as JSON")
}
