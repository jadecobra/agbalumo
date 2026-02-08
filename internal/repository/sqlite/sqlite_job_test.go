package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestSaveAndFindJob(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second) // Truncate for DB precision

	job := domain.Listing{
		ID:           "job-1",
		Title:        "Go Developer",
		OwnerOrigin:  "Nigeria",
		Type:         domain.Job,
		Description:  "Write Go code",
		Company:      "TechCorp",
		PayRange:     "$100k - $150k",
		Skills:       "Go, SQL, Docker",
		JobStartDate: start,
		JobApplyURL:  "https://example.com/apply",
		ContactEmail: "hr@company.com",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	// 1. Save
	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("Failed to save job: %v", err)
	}

	// 2. Find
	found, err := repo.FindByID(ctx, "job-1")
	if err != nil {
		t.Fatalf("Failed to find job: %v", err)
	}

	// 3. Verify Job Specific Fields
	if found.Company != job.Company {
		t.Errorf("Expected company '%s', got '%s'", job.Company, found.Company)
	}
	if found.PayRange != job.PayRange {
		t.Errorf("Expected pay range '%s', got '%s'", job.PayRange, found.PayRange)
	}
	if found.Skills != job.Skills {
		t.Errorf("Expected skills '%s', got '%s'", job.Skills, found.Skills)
	}
	if !found.JobStartDate.Equal(job.JobStartDate) {
		t.Errorf("Expected start date %v, got %v", job.JobStartDate, found.JobStartDate)
	}
	if found.JobApplyURL != job.JobApplyURL {
		t.Errorf("Expected apply URL '%s', got '%s'", job.JobApplyURL, found.JobApplyURL)
	}
}
