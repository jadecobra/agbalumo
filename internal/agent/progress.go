package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jadecobra/agbalumo/internal/util"
)

type Feature struct {
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Passes      bool     `json:"passes"`
	Steps       []string `json:"steps"`
}

type ProgressTracker struct {
	Features []Feature `json:"features"`
}

// HasPending checks if any steps are NOT completed
func HasPending(steps []string) bool {
	for _, step := range steps {
		if !strings.Contains(step, "(Completed)") {
			return true
		}
	}
	return false
}

// ParseMarkdownTracker parses a Markdown task list into a ProgressTracker.
func ParseMarkdownTracker(content string) (ProgressTracker, error) {
	lines := strings.Split(content, "\n")
	var tracker ProgressTracker
	var currentFeature *Feature

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ") {
			// New Category
			if currentFeature != nil {
				currentFeature.Passes = !HasPending(currentFeature.Steps)
				tracker.Features = append(tracker.Features, *currentFeature)
			}
			currentFeature = &Feature{
				Category: strings.TrimPrefix(line, "# "),
			}
		} else if strings.HasPrefix(line, "- [") {
			// Step
			if currentFeature == nil {
				continue
			}
			step := strings.TrimSpace(line[5:])
			// [x] or [ ]
			isDone := strings.HasPrefix(line, "- [x]")
			if isDone {
				if !strings.Contains(step, "(Completed)") {
					step += " (Completed)"
				}
			}
			currentFeature.Steps = append(currentFeature.Steps, step)
		} else if currentFeature != nil && currentFeature.Description == "" {
			// Description (if not a heading or step)
			currentFeature.Description = line
		}
	}

	if currentFeature != nil {
		currentFeature.Passes = !HasPending(currentFeature.Steps)
		tracker.Features = append(tracker.Features, *currentFeature)
	}

	return tracker, nil
}

// ToMarkdown serializes a ProgressTracker to Markdown.
func ToMarkdown(tracker ProgressTracker) string {
	var sb strings.Builder
	for _, f := range tracker.Features {
		fmt.Fprintf(&sb, "# %s\n", f.Category)
		if f.Description != "" {
			fmt.Fprintf(&sb, "%s\n", f.Description)
		}
		for _, step := range f.Steps {
			status := "[ ]"
			if strings.Contains(step, "(Completed)") {
				status = "[x]"
			}
			cleanStep := strings.TrimSuffix(step, " (Completed)")
			fmt.Fprintf(&sb, "- %s %s\n", status, cleanStep)
		}
	}
	return sb.String()
}

// ArchivePassedCategories moves passed features to an archive file if the threshold is met.
func ArchivePassedCategories(progressPath, archivePath string, threshold int) error {
	// 1. Read progress file (support both .json and .md during migration)
	data, err := util.SafeReadFile(progressPath)
	if err != nil {
		if util.SafeIsNotExist(err) {
			return nil
		}
		return err
	}

	var pTracker ProgressTracker
	isMarkdown := strings.HasSuffix(progressPath, ".md")

	if isMarkdown {
		pTracker, err = ParseMarkdownTracker(string(data))
		if err != nil {
			return err
		}
	} else {
		if err = json.Unmarshal(data, &pTracker); err != nil {
			return err
		}
	}

	// 2. Check if we need to archive
	if len(pTracker.Features) <= threshold {
		return nil
	}

	// 3. Separate passed and pending features
	var passed []Feature
	var pending []Feature

	for _, f := range pTracker.Features {
		if f.Passes {
			passed = append(passed, f)
		} else {
			pending = append(pending, f)
		}
	}

	// If we have nothing to archive, just return
	if len(passed) == 0 {
		return nil
	}

	// 4. Update archive
	var archivedTracker ProgressTracker
	archiveData, err := util.SafeReadFile(archivePath)
	if err == nil {
		if strings.HasSuffix(archivePath, ".md") {
			archivedTracker, _ = ParseMarkdownTracker(string(archiveData))
		} else {
			_ = json.Unmarshal(archiveData, &archivedTracker)
		}
	}
	archivedTracker.Features = append(archivedTracker.Features, passed...)

	var newArchiveData []byte
	if strings.HasSuffix(archivePath, ".md") {
		newArchiveData = []byte(ToMarkdown(archivedTracker))
	} else {
		newArchiveData, err = json.MarshalIndent(archivedTracker, "", "  ")
		if err != nil {
			return err
		}
	}
	if err = util.SafeWriteFile(archivePath, newArchiveData); err != nil {
		return err
	}

	// 5. Update progress file
	pTracker.Features = pending
	var newProgressData []byte
	if isMarkdown {
		newProgressData = []byte(ToMarkdown(pTracker))
	} else {
		newProgressData, err = json.MarshalIndent(pTracker, "", "  ")
		if err != nil {
			return err
		}
	}
	if err := util.SafeWriteFile(progressPath, newProgressData); err != nil {
		return err
	}

	fmt.Printf("📦 Archived %d passed features to %s\n", len(passed), archivePath)
	return nil
}
