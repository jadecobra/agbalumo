package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventValidation(t *testing.T) {
	t.Parallel()
	now := time.Now()

	baseListing := func() Listing {
		return Listing{
			ID:           "test-event",
			OwnerOrigin:  "Nigeria",
			Type:         Event,
			Title:        "Test Event",
			ContactEmail: "event@example.com",
			Address:      "123 Valid St",
			City:         "Lagos",
			CreatedAt:    now,
			IsActive:     true,
			EventStart:   now.Add(24 * time.Hour),
			EventEnd:     now.Add(26 * time.Hour),
		}
	}

	runEventTest := func(name string, mutate func(*Listing), wantErr string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := baseListing()
			mutate(&l)
			err := l.Validate()
			if wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	runEventTest("Valid Case", func(l *Listing) {}, "")
	runEventTest("Event Missing Start Date", func(l *Listing) { l.EventStart = time.Time{} }, "event start time is required")
	runEventTest("Event Missing End Date", func(l *Listing) { l.EventEnd = time.Time{} }, "event end time is required")
	runEventTest("Event End Before Start", func(l *Listing) {
		l.EventStart = now.Add(26 * time.Hour)
		l.EventEnd = now.Add(24 * time.Hour)
	}, "event end time cannot be before start time")
	runEventTest("Non-Event Ignores Dates", func(l *Listing) {
		l.Type = Business
		l.EventStart = time.Time{}
		l.EventEnd = time.Time{}
	}, "")
}
