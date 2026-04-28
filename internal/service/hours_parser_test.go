package service

import (
	"testing"
	"time"
)

func TestComputeIsOpen(t *testing.T) {
	// Let's define a fixed date for deterministic testing
	// 2026-04-27 is a Monday.
	baseTime := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		currentTime     time.Time
		name            string
		hoursText       string
		structuredHours string
		want            bool
	}{
		{
			name:            "Mon-Fri 9am-10pm - Open on Monday",
			hoursText:       "Mon-Fri 9am-10pm",
			structuredHours: "",
			currentTime:     time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC), // Monday 10:00 AM
			want:            true,
		},

		{
			name:        "Mon-Fri 9am-10pm - Closed on Monday night",
			hoursText:   "Mon-Fri 9am-10pm",
			currentTime: time.Date(2026, 4, 27, 23, 0, 0, 0, time.UTC), // Monday 11:00 PM
			want:        false,
		},
		{
			name:        "Daily 11am-9pm - Open on Saturday",
			hoursText:   "Daily 11am-9pm",
			currentTime: time.Date(2026, 5, 2, 14, 0, 0, 0, time.UTC), // Saturday 2:00 PM
			want:        true,
		},
		{
			name:        "Mon-Sun 10:00 - 22:00 - Open on Sunday",
			hoursText:   "Mon-Sun 10:00 - 22:00",
			currentTime: time.Date(2026, 5, 3, 20, 0, 0, 0, time.UTC), // Sunday 8:00 PM
			want:        true,
		},
		{
			name:        "Mon-Sat 9:00 AM - 9:00 PM - Closed on Sunday",
			hoursText:   "Mon-Sat 9:00 AM - 9:00 PM",
			currentTime: time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC), // Sunday 12:00 PM
			want:        false,
		},
		{
			name:        "Open 24 Hours - Open any time",
			hoursText:   "Open 24 Hours",
			currentTime: baseTime,
			want:        true,
		},
		{
			name:            "Ambiguous hours regex fail - Structured Success (Sat)",
			hoursText:       "Mon-Fri 9am-10pm, Sat 10am-2pm",
			structuredHours: `{"sat": ["10:00-14:00"]}`,
			currentTime:     time.Date(2026, 5, 2, 11, 0, 0, 0, time.UTC), // Saturday 11:00 AM
			want:            true,
		},
		{
			name:            "Ambiguous hours regex fail - Structured Success (Closed)",
			hoursText:       "Mon-Fri 9am-10pm, Sat 10am-2pm",
			structuredHours: `{"sat": ["10:00-14:00"]}`,
			currentTime:     time.Date(2026, 5, 2, 15, 0, 0, 0, time.UTC), // Saturday 3:00 PM
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeIsOpen(tt.hoursText, tt.structuredHours, tt.currentTime)
			if got != tt.want {
				t.Errorf("ComputeIsOpen(%q, %q, %v) = %v, want %v", tt.hoursText, tt.structuredHours, tt.currentTime, got, tt.want)
			}
		})
	}

}
