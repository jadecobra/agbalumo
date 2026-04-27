package service

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ComputeIsOpen determines if a listing is currently open based on its unstructured hours text.
func ComputeIsOpen(hoursText string, currentTime time.Time) bool {
	hoursText = strings.ToLower(strings.TrimSpace(hoursText))
	if hoursText == "" {
		return false
	}

	if strings.Contains(hoursText, "open 24 hours") {
		return true
	}

	// 1. Check for explicit closed days
	dayNames := []string{"", "mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	currentWeekday := int(currentTime.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7 // Map Sunday from 0 to 7
	}

	// Examples: "sun closed", "closed sundays", "closed on sunday"
	for i := 1; i <= 7; i++ {
		if i == currentWeekday {
			closedPattern := regexp.MustCompile(`(?i)` + dayNames[i] + `.*\bclosed\b|\bclosed.*\b` + dayNames[i])
			if closedPattern.MatchString(hoursText) {
				return false
			}
		}
	}

	// 2. Check day ranges
	// Example: "mon-fri", "mon - sun", "daily"
	dayRangePattern := regexp.MustCompile(`(?i)(mon|tue|wed|thu|fri|sat|sun)\s*(?:-|to)\s*(mon|tue|wed|thu|fri|sat|sun)`)
	rangeMatch := dayRangePattern.FindStringSubmatch(hoursText)

	dayMatch := false
	if len(rangeMatch) == 3 {
		startDay := getDayNumber(rangeMatch[1])
		endDay := getDayNumber(rangeMatch[2])

		if startDay <= endDay {
			if currentWeekday >= startDay && currentWeekday <= endDay {
				dayMatch = true
			}
		} else {
			// Wraps around, e.g., Fri-Tue
			if currentWeekday >= startDay || currentWeekday <= endDay {
				dayMatch = true
			}
		}
	} else if strings.Contains(hoursText, "daily") || strings.Contains(hoursText, "mon-sun") || strings.Contains(hoursText, "mon - sun") {
		dayMatch = true
	} else {
		// If no explicit range, assume it applies to the current day unless excluded
		dayMatch = true
	}

	if !dayMatch {
		return false
	}

	// 3. Extract times
	// Matches patterns like: 9am, 10:30pm, 22:00, 9:00 am
	timePattern := regexp.MustCompile(`(\d{1,2})(?::(\d{2}))?\s*(am|pm)?`)
	matches := timePattern.FindAllStringSubmatch(hoursText, -1)

	// We need at least an opening and closing time (2 matches)
	if len(matches) < 2 {
		return false
	}

	// Let's assume the first match is opening, and the last match (or second) is closing
	openMatch := matches[0]
	closeMatch := matches[1]

	openMinutes, ok1 := parseTimeToMinutes(openMatch[1], openMatch[2], openMatch[3])
	closeMinutes, ok2 := parseTimeToMinutes(closeMatch[1], closeMatch[2], closeMatch[3])

	if !ok1 || !ok2 {
		return false
	}

	// Heuristic: If opening lacks AM/PM but closing has PM
	if openMatch[3] == "" && closeMatch[3] == "pm" {
		openHour, _ := strconv.Atoi(openMatch[1])
		closeHour, _ := strconv.Atoi(closeMatch[1])
		if openHour > closeHour {
			// e.g., 9-5pm -> 9 AM to 5 PM
			openMinutes += 0 // It's already AM
		} else if openHour < closeHour {
			// e.g., 1-5pm -> 1 PM to 5 PM
			if openHour < 12 {
				openMinutes += 12 * 60 // Make it PM
			}
		}
	}

	currentMinutes := currentTime.Hour()*60 + currentTime.Minute()

	if closeMinutes < openMinutes {
		// Overnight hours, e.g., 4pm - 2am
		// Open if current time is >= openMinutes OR current time is < closeMinutes
		return currentMinutes >= openMinutes || currentMinutes < closeMinutes
	}

	return currentMinutes >= openMinutes && currentMinutes < closeMinutes
}

func getDayNumber(day string) int {
	switch strings.ToLower(day) {
	case "mon":
		return 1
	case "tue":
		return 2
	case "wed":
		return 3
	case "thu":
		return 4
	case "fri":
		return 5
	case "sat":
		return 6
	case "sun":
		return 7
	}
	return 0
}

func parseTimeToMinutes(hourStr, minStr, ampmStr string) (int, bool) {
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return 0, false
	}
	min := 0
	if minStr != "" {
		min, _ = strconv.Atoi(minStr)
	}

	ampm := strings.ToLower(ampmStr)
	if ampm == "pm" && hour < 12 {
		hour += 12
	} else if ampm == "am" && hour == 12 {
		hour = 0
	}

	return hour*60 + min, true
}
