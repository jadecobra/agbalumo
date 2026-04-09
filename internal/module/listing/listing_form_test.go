package listing

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestParseDeadline(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{DeadlineDate: "2024-12-31"}
	l := &domain.Listing{Type: domain.Request}

	err := parseDeadline(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.Deadline.Year())

	req.DeadlineDate = "invalid"
	err = parseDeadline(req, l)
	assert.Error(t, err)
}

func TestParseEventDates(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{
		EventStart: "2024-12-01T10:00",
		EventEnd:   "2024-12-01T12:00",
	}
	l := &domain.Listing{Type: domain.Event}

	err := parseEventDates(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.EventStart.Year())

	req.EventStart = "invalid"
	err = parseEventDates(req, l)
	assert.Error(t, err)

	req.EventStart = "2024-12-01T10:00"
	req.EventEnd = "invalid"
	err = parseEventDates(req, l)
	assert.Error(t, err)
}

func TestParseJobStartDate(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{JobStartDate: "2024-12-01T09:00"}
	l := &domain.Listing{Type: domain.Job}

	err := parseJobStartDate(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.JobStartDate.Year())

	req.JobStartDate = "invalid"
	err = parseJobStartDate(req, l)
	assert.Error(t, err)
}
