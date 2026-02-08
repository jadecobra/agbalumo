package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventValidation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		eventStart time.Time
		eventEnd   time.Time
		typeStr    Category
		wantErr    string // Substring match
	}{
		{
			name:       "Event with valid dates",
			eventStart: now.Add(24 * time.Hour),
			eventEnd:   now.Add(26 * time.Hour),
			typeStr:    Event,
			wantErr:    "",
		},
		{
			name:       "Event missing start date",
			eventStart: time.Time{}, // Zero
			eventEnd:   now.Add(26 * time.Hour),
			typeStr:    Event,
			wantErr:    "event start time is required",
		},
		{
			name:       "Event missing end date",
			eventStart: now.Add(24 * time.Hour),
			eventEnd:   time.Time{}, // Zero
			typeStr:    Event,
			wantErr:    "event end time is required",
		},
		{
			name:       "Event end before start",
			eventStart: now.Add(26 * time.Hour),
			eventEnd:   now.Add(24 * time.Hour),
			typeStr:    Event,
			wantErr:    "event end time cannot be before start time",
		},
		{
			name:       "Non-Event ignores dates",
			eventStart: time.Time{},
			eventEnd:   time.Time{},
			typeStr:    Business,
			wantErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:           "test-event",
				OwnerOrigin:  "Nigeria",
				Type:         tt.typeStr,
				Title:        "Test Event",
				ContactEmail: "event@example.com",
				CreatedAt:    now,
				IsActive:     true,
				// These fields don't exist yet, so this code won't compile initially (RED phase if we consider compilation failure as part of it, 
				// but to run the test tool we usually need it to compile. 
				// However, strictly speaking, if the struct fields are missing, I can't even write this test without compilation errors.
				// I will add the fields to the struct in the next step, but for `go test` to not completely barf on the file, 
				// I usually add the fields first OR rely on the fact that I'm about to add them. 
				// But true RED in Go usually implies writing the test that fails to compile or fails assertion.
				// Since I can't run a test that doesn't compile, I will assume the fields exist for the purpose of the test logic, 
				// but I will add the empty fields to the struct in the next step or just the fields now.
				// Actually, to make it compile so I can see it fail (logic wise) or just compile failure IS the test failure.
				// I will write the test assuming fields exist.
				EventStart: tt.eventStart,
				EventEnd:   tt.eventEnd,
			}
			
			err := l.Validate()
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
