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

// ArchivePassedCategories moves passed features to an archive file if the threshold is met.
func ArchivePassedCategories(progressPath, archivePath string, threshold int) error {
	// 1. Read progress.json
	data, err := util.SafeReadFile(progressPath)
	if err != nil {
		if util.SafeIsNotExist(err) {
			return nil
		}
		return err
	}

	var tracker ProgressTracker
	if err := json.Unmarshal(data, &tracker); err != nil {
		return err
	}

	// 2. Check if we need to archive
	if len(tracker.Features) <= threshold {
		return nil
	}

	// 3. Separate passed and pending features
	var passed []Feature
	var pending []Feature

	for _, f := range tracker.Features {
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
		_ = json.Unmarshal(archiveData, &archivedTracker)
	}
	archivedTracker.Features = append(archivedTracker.Features, passed...)

	newArchiveData, err := json.MarshalIndent(archivedTracker, "", "  ")
	if err != nil {
		return err
	}
	if err := util.SafeWriteFile(archivePath, newArchiveData); err != nil {
		return err
	}

	// 5. Update progress.json
	tracker.Features = pending
	newProgressData, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		return err
	}
	if err := util.SafeWriteFile(progressPath, newProgressData); err != nil {
		return err
	}

	fmt.Printf("📦 Archived %d passed features to %s\n", len(passed), archivePath)
	return nil
}
