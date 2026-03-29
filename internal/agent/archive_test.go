package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestArchivePassedCategories(t *testing.T) {
	tmpDir := t.TempDir()
	progressPath := filepath.Join(tmpDir, "progress.json")
	archivePath := filepath.Join(tmpDir, "progress_archive.json")

	// Setup: Create a progress.json with many passed features and some pending ones
	features := []Feature{}
	for i := 1; i <= 25; i++ {
		features = append(features, Feature{
			Category:    "Passed Cat",
			Description: "Description",
			Passes:      true,
			Steps:       []string{"Step (Completed)"},
		})
	}
	features = append(features, Feature{
		Category:    "Pending Cat",
		Description: "Description",
		Passes:      false,
		Steps:       []string{"Step"},
	})

	tracker := ProgressTracker{Features: features}
	data, _ := json.MarshalIndent(tracker, "", "  ")
	_ = os.WriteFile(progressPath, data, 0644)

	// Threshold is 20 features
	err := ArchivePassedCategories(progressPath, archivePath, 20)
	if err != nil {
		t.Fatalf("ArchivePassedCategories failed: %v", err)
	}

	// Verify progress.json
	newData, _ := os.ReadFile(progressPath)
	var newTracker ProgressTracker
	_ = json.Unmarshal(newData, &newTracker)

	// Should have moved enough to get under/equal to 20
	if len(newTracker.Features) != 1 {
		t.Errorf("expected 1 feature in progress.json (pending), got %d", len(newTracker.Features))
	}

	// Verify archive exists and contains the passed features
	archiveData, _ := os.ReadFile(archivePath)
	var archivedTracker ProgressTracker
	_ = json.Unmarshal(archiveData, &archivedTracker)
	if len(archivedTracker.Features) != 25 {
		t.Errorf("expected 25 features in archive, got %d", len(archivedTracker.Features))
	}

	// Test: Archive not needed (threshold not met)
	err = ArchivePassedCategories(progressPath, archivePath, 100)
	if err != nil {
		t.Errorf("ArchivePassedCategories failed on high threshold: %v", err)
	}

	// Test: No progress file (should not error)
	err = ArchivePassedCategories("non_existent.json", archivePath, 0)
	if err != nil {
		t.Errorf("ArchivePassedCategories failed on missing file: %v", err)
	}

	// Test: No passed features (should not archive)
	_ = os.WriteFile(progressPath, []byte(`{"features":[{"category":"A"}]}`), 0644)
	err = ArchivePassedCategories(progressPath, archivePath, 0)
	if err != nil {
		t.Errorf("ArchivePassedCategories failed on no passes: %v", err)
	}
}

func TestHasPending(t *testing.T) {
	tests := []struct {
		steps    []string
		expected bool
	}{
		{[]string{"Step 1 (Completed)", "Step 2 (Completed)"}, false},
		{[]string{"Step 1 (Completed)", "Step 2"}, true},
		{[]string{"Step 1"}, true},
		{[]string{}, false},
	}

	for _, tt := range tests {
		if got := HasPending(tt.steps); got != tt.expected {
			t.Errorf("HasPending(%v) = %v; want %v", tt.steps, got, tt.expected)
		}
	}
}
