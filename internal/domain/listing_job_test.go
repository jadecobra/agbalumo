package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListing_Validate_Job(t *testing.T) {
	t.Parallel()

	validJob := func() Listing {
		return Listing{
			Title:        "Software Engineer",
			Type:         Job,
			OwnerOrigin:  "Nigeria",
			Description:  "Write code",
			Company:      "TechCorp",
			PayRange:     "$100k - $150k",
			Skills:       "Go, SQLite",
			JobStartDate: time.Now().Add(24 * time.Hour),
			JobApplyURL:  "https://example.com/apply",
			City:         "Lagos",
			ContactEmail: "hr@example.com",
		}
	}

	runTest := func(name string, mutate func(*Listing), wantErr string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := validJob()
			mutate(&l)
			err := l.Validate()
			if wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), wantErr)
			}
		})
	}

	runTest("Valid Job", func(l *Listing) {}, "")
	runTest("Missing Skills", func(l *Listing) { l.Skills = "" }, "skills are required")
	runTest("Missing Description", func(l *Listing) { l.Description = "" }, "description is required")
	runTest("Missing Compensation", func(l *Listing) { l.PayRange = "" }, "compensation/pay range is required")
	runTest("Job Missing Start Date", func(l *Listing) { l.JobStartDate = time.Time{} }, "job start date is required")
	runTest("Job Start Date in Past", func(l *Listing) { l.JobStartDate = time.Now().Add(-48 * time.Hour) }, "job start date cannot be in the past")
	runTest("Job Missing Company", func(l *Listing) { l.Company = "" }, "company name is required")
	runTest("Job Missing City", func(l *Listing) { l.City = "" }, "city is required")
	runTest("Job Missing Apply URL", func(l *Listing) { l.JobApplyURL = "" }, "apply url is required")
}
