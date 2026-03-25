package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHoursOfOperationField(t *testing.T) {
	// This test ensures the field exists and can be set.
	l := Listing{
		ID:               "test-hours",
		OwnerOrigin:      "Togo",
		Type:             Business,
		Title:            "Hours Test",
		ContactEmail:     "hours@example.com",
		Address:          "Main St",
		HoursOfOperation: "Mon-Fri 9-5",
		CreatedAt:        time.Now(),
	}

	assert.Equal(t, "Mon-Fri 9-5", l.HoursOfOperation)
}
