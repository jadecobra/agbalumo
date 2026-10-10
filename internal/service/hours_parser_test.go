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

func TestAdjustImplicitPM(t *testing.T) {
	tests := []struct {
		name         string
		openHourStr  string
		closeHourStr string
		openMinutes  int
		want         int
	}{
		{
			name:         "openHour < closeHour and openHour < 12 adjusts to PM",
			openHourStr:  "1",
			closeHourStr: "5",
			openMinutes:  60,
			want:         780,
		},
		{
			name:         "openHour >= closeHour remains unchanged",
			openHourStr:  "9",
			closeHourStr: "5",
			openMinutes:  540,
			want:         540,
		},
		{
			name:         "equal hours remain unchanged",
			openHourStr:  "5",
			closeHourStr: "5",
			openMinutes:  300,
			want:         300,
		},
		{
			name:         "openHour >= 12 remains unchanged",
			openHourStr:  "12",
			closeHourStr: "5",
			openMinutes:  720,
			want:         720,
		},
		{
			name:         "invalid openHourStr returns original minutes",
			openHourStr:  "invalid",
			closeHourStr: "5",
			openMinutes:  60,
			want:         60,
		},
		{
			name:         "invalid closeHourStr returns original minutes",
			openHourStr:  "1",
			closeHourStr: "invalid",
			openMinutes:  60,
			want:         60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adjustImplicitPM(tt.openHourStr, tt.closeHourStr, tt.openMinutes)
			if got != tt.want {
				t.Errorf("adjustImplicitPM(%q, %q, %d) = %d, want %d", tt.openHourStr, tt.closeHourStr, tt.openMinutes, got, tt.want)
			}
		})
	}
}

func TestExtractOpenCloseMinutes(t *testing.T) {
	tests := []struct {
		name      string
		matches   [][]string
		wantOpen  int
		wantClose int
		wantOk    bool
	}{
		{
			name: "1-5pm parses with implicit PM on opening",
			matches: [][]string{
				{"1", "1", "", ""},
				{"5pm", "5", "", "pm"},
			},
			wantOpen:  780,
			wantClose: 1020,
			wantOk:    true,
		},
		{
			name: "9-5pm parses with AM opening and PM closing",
			matches: [][]string{
				{"9", "9", "", ""},
				{"5pm", "5", "", "pm"},
			},
			wantOpen:  540,
			wantClose: 1020,
			wantOk:    true,
		},
		{
			name: "invalid opening time returns false",
			matches: [][]string{
				{"abc", "abc", "", ""},
				{"5pm", "5", "", "pm"},
			},
			wantOpen:  0,
			wantClose: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, closeMin, ok := extractOpenCloseMinutes(tt.matches)
			if ok != tt.wantOk || open != tt.wantOpen || closeMin != tt.wantClose {
				t.Errorf("extractOpenCloseMinutes(%v) = (%d, %d, %v), want (%d, %d, %v)",
					tt.matches, open, closeMin, ok, tt.wantOpen, tt.wantClose, tt.wantOk)
			}
		})
	}
}

func TestComputeIsOpen_ImplicitPM(t *testing.T) {
	mondayAfternoon := time.Date(2026, 4, 27, 14, 0, 0, 0, time.UTC)
	mondayMorning := time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	mondayEarly := time.Date(2026, 4, 27, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		currentTime time.Time
		name        string
		hoursText   string
		want        bool
	}{
		{
			name:        "1-5pm open during afternoon",
			hoursText:   "Mon-Fri 1-5pm",
			currentTime: mondayAfternoon,
			want:        true,
		},
		{
			name:        "1-5pm closed during morning",
			hoursText:   "Mon-Fri 1-5pm",
			currentTime: mondayMorning,
			want:        false,
		},
		{
			name:        "9-5pm open during morning",
			hoursText:   "Mon-Fri 9-5pm",
			currentTime: mondayMorning,
			want:        true,
		},
		{
			name:        "9-5pm closed before 9am",
			hoursText:   "Mon-Fri 9-5pm",
			currentTime: mondayEarly,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeIsOpen(tt.hoursText, "", tt.currentTime)
			if got != tt.want {
				t.Errorf("ComputeIsOpen(%q, \"\", %v) = %v, want %v", tt.hoursText, tt.currentTime, got, tt.want)
			}
		})
	}
}
