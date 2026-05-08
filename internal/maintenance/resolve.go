package maintenance

import (
	"strings"
)

type ResolvedMatch struct {
	Name        string
	Path        string // For skills
	Command     string // For commands
	Description string
	Type        string // "skill" or "command"
}

func ResolveIntent(rootDir, intent string) ([]ResolvedMatch, error) {
	manifest, err := LoadVerifyManifest(rootDir)
	if err != nil {
		return nil, err
	}

	intent = strings.ToLower(strings.TrimSpace(intent))
	var matches []ResolvedMatch

	// Check Skills
	for _, s := range manifest.Skills {
		if intentMatches(intent, s.Trigger) {
			matches = append(matches, ResolvedMatch{
				Name: s.Name,
				Path: s.Path,
				Type: "skill",
			})
		}
	}

	// Check Commands
	for _, c := range manifest.Commands {
		if intentMatches(intent, c.Trigger) {
			matches = append(matches, ResolvedMatch{
				Name:        c.Name,
				Description: c.Description,
				Type:        "command",
			})
		}
	}

	return matches, nil
}

func intentMatches(intent, triggerStr string) bool {
	triggers := strings.Split(triggerStr, ",")
	iNorm := strings.ReplaceAll(strings.ToLower(intent), "_", " ")

	for _, t := range triggers {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		tNorm := strings.ReplaceAll(t, "_", " ")

		if strings.Contains(iNorm, tNorm) || strings.Contains(tNorm, iNorm) {
			return true
		}
		if wordMatches(iNorm, tNorm) {
			return true
		}
	}
	return false
}

func wordMatches(intentNorm, triggerNorm string) bool {
	triggerWords := strings.Fields(triggerNorm)
	if len(triggerWords) == 0 {
		return false
	}

	intentWords := strings.Fields(intentNorm)
	intentWordMap := make(map[string]bool)
	for _, w := range intentWords {
		intentWordMap[w] = true
	}

	matchCount := 0
	for _, tw := range triggerWords {
		if intentWordMap[tw] {
			matchCount++
		}
	}

	// Requirement: strictly more than half of the trigger's words appear in the intent
	// This ensures "config change" (2 words) matches "ci_config_change" (2/3 matches)
	// but "config change" does NOT match "ui_change" (1/2 matches)
	return matchCount*2 > len(triggerWords)
}
