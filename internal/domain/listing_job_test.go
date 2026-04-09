package domain

import (
	"testing"
	"time"
)

func TestListing_Validate_Job(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		listing Listing
		wantErr bool
	}{
		{
			name: "Valid Job",
			listing: Listing{
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
			},
			wantErr: false,
		},
		{
			name: "Missing Skills",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Ghana",
				Description:  "Write code",
				City:         "Accra",
				JobStartDate: time.Now().Add(24 * time.Hour),
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
		{
			name: "Missing Start Date",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Kenya",
				Description:  "Write code",
				City:         "Nairobi",
				Skills:       "Go",
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
		{
			name: "Start Date in Past",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Nigeria",
				Description:  "Write code",
				City:         "Lagos",
				Skills:       "Go",
				JobStartDate: time.Now().Add(-24 * time.Hour),
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
		{
			name: "Missing Company",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Nigeria",
				Description:  "Write code",
				City:         "Lagos",
				Skills:       "Go",
				JobStartDate: time.Now().Add(24 * time.Hour),
				JobApplyURL:  "https://example.com/apply",
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
		{
			name: "Missing Location",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Nigeria",
				Description:  "Write code",
				Company:      "TechCorp",
				PayRange:     "$100k - $150k",
				Skills:       "Go",
				JobStartDate: time.Now().Add(24 * time.Hour),
				JobApplyURL:  "https://example.com/apply",
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
		{
			name: "Missing Apply URL",
			listing: Listing{
				Title:        "Software Engineer",
				Type:         Job,
				OwnerOrigin:  "Nigeria",
				Description:  "Write code",
				Company:      "TechCorp",
				PayRange:     "$100k - $150k",
				Skills:       "Go",
				JobStartDate: time.Now().Add(24 * time.Hour),
				City:         "Lagos",
				ContactEmail: "hr@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.listing.Validate()
			if err != nil {
				t.Logf("Helper log: error was %v", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Listing.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
