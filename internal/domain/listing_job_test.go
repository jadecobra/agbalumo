package domain

import (
	"testing"
	"time"
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

	tests := []struct {
		name    string
		listing Listing
		wantErr bool
	}{
		{
			name:    "Valid Job",
			listing: validJob(),
			wantErr: false,
		},
		{
			name: "Missing Skills",
			listing: func() Listing {
				l := validJob()
				l.Skills = ""
				return l
			}(),
			wantErr: true,
		},
		{
			name: "Missing Start Date",
			listing: func() Listing {
				l := validJob()
				l.JobStartDate = time.Time{}
				return l
			}(),
			wantErr: true,
		},
		{
			name: "Start Date in Past",
			listing: func() Listing {
				l := validJob()
				l.JobStartDate = time.Now().Add(-24 * time.Hour)
				return l
			}(),
			wantErr: true,
		},
		{
			name: "Missing Company",
			listing: func() Listing {
				l := validJob()
				l.Company = ""
				return l
			}(),
			wantErr: true,
		},
		{
			name: "Missing Location",
			listing: func() Listing {
				l := validJob()
				l.City = ""
				return l
			}(),
			wantErr: true,
		},
		{
			name: "Missing Apply URL",
			listing: func() Listing {
				l := validJob()
				l.JobApplyURL = ""
				return l
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.listing.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Listing.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
