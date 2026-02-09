package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobListingStrictValidation(t *testing.T) {
	// Base Valid Job Listing
	validJob := Listing{
		ID:           "job-1",
		OwnerOrigin:  "Nigeria",
		Type:         Job,
		Title:        "Senior Go Engineer",
		Description:  "Build robust backends.",
		Company:      "Tech Global",
		Skills:       "Go, SQL, TDD",
		PayRange:     "$80k - $120k",
		JobStartDate: time.Now().Add(24 * time.Hour),
		JobApplyURL:  "https://example.com/apply",
		City:         "Lagos",
		ContactEmail: "hr@techglobal.com",
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	tests := []struct {
		name    string
		mutate  func(l *Listing)
		wantErr string // Substring to match in error
	}{
		{
			name:    "Valid Job",
			mutate:  func(l *Listing) {},
			wantErr: "",
		},
		{
			name: "Missing Organization (Company)",
			mutate: func(l *Listing) {
				l.Company = ""
			},
			wantErr: "company name is required",
		},
		{
			name: "Missing Description",
			mutate: func(l *Listing) {
				l.Description = ""
			},
			wantErr: "description is required", // Note: This might need to be added to Validate()
		},
		{
			name: "Missing Skills",
			mutate: func(l *Listing) {
				l.Skills = ""
			},
			wantErr: "skills are required",
		},
		{
			name: "Missing Compensation (PayRange)",
			mutate: func(l *Listing) {
				l.PayRange = ""
			},
			wantErr: "compensation/pay range is required", // Expecting this to fail initially
		},
		{
			name: "Missing Start Date",
			mutate: func(l *Listing) {
				l.JobStartDate = time.Time{}
			},
			wantErr: "job start date is required",
		},
		{
			name: "Start Date in Past",
			mutate: func(l *Listing) {
				l.JobStartDate = time.Now().Add(-48 * time.Hour)
			},
			wantErr: "job start date cannot be in the past",
		},
		{
			name: "Missing Location (City)",
			mutate: func(l *Listing) {
				l.City = ""
				l.Address = ""
			},
			wantErr: "location (city) is required", // Expecting this to fail initially
		},
		{
			name: "Missing Apply URL",
			mutate: func(l *Listing) {
				l.JobApplyURL = ""
			},
			wantErr: "apply url is required", // Expecting this to fail initially
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := validJob // Copy
			tt.mutate(&l)
			err := l.Validate()

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.wantErr)
				}
			}
		})
	}
}
