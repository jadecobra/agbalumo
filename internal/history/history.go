package history

import (
	"fmt"
	"time"

	"github.com/jadecobra/agbalumo/internal/util"
)

// SquadDecision represents a milestone decision made by the agent squad.
type SquadDecision struct {
	FeatureName      string
	SystemsArchitect string
	ProductOwner     string
	SDET             string
	BackendEngineer  string
	DecisionSummary  string
}

// DefaultStorageDir is where history files are stored.
var DefaultStorageDir = "internal/history"

// Store persists a SquadDecision to a Markdown file and returns the path.
func Store(decision SquadDecision) (string, error) {
	dir := DefaultStorageDir
	if err := util.SafeMkdir(dir); err != nil {
		return "", fmt.Errorf("failed to create history directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s/%s_%s.md", dir, timestamp, decision.FeatureName)

	content := fmt.Sprintf(`---
FeatureName: %s
SystemsArchitect: %s
ProductOwner: %s
SDET: %s
BackendEngineer: %s
Date: %s
---

# [SUMMARY]
%s
`,
		decision.FeatureName,
		decision.SystemsArchitect,
		decision.ProductOwner,
		decision.SDET,
		decision.BackendEngineer,
		time.Now().Format(time.RFC3339),
		decision.DecisionSummary,
	)

	if err := util.SafeWriteFile(filename, []byte(content)); err != nil {
		return "", fmt.Errorf("failed to write decision file: %w", err)
	}

	return filename, nil
}
