package service

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ComputeIsOpen determines if a listing is currently open based on its unstructured hours text.
func ComputeIsOpen(hoursText string, structuredHours string, currentTime time.Time) bool {
	if isOpen, evaluated := evaluateStructuredHours(structuredHours, currentTime); evaluated {
		return isOpen
	}

	hoursText = strings.ToLower(strings.TrimSpace(hoursText))
	if hoursText == "" {
		return false
	}

	if strings.Contains(hoursText, "open 24 hours") {
		return true
	}

	if isClosedToday(hoursText, currentTime) {
		return false
	}

	if !isDayMatch(hoursText, currentTime) {
		return false
	}

	return isTimeOpen(hoursText, currentTime)
}

func isClosedToday(hoursText string, currentTime time.Time) bool {
	dayNames := []string{"", "mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	currentWeekday := int(currentTime.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}

	for i := 1; i <= 7; i++ {
		if i == currentWeekday {
			closedPattern := regexp.MustCompile(`(?i)` + dayNames[i] + `.*\bclosed\b|\bclosed.*\b` + dayNames[i])
			if closedPattern.MatchString(hoursText) {
				return true
			}
		}
	}
	return false
}

func isDayMatch(hoursText string, currentTime time.Time) bool {
	currentWeekday := int(currentTime.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}

	dayRangePattern := regexp.MustCompile(`(?i)(mon|tue|wed|thu|fri|sat|sun)\s*(?:-|to)\s*(mon|tue|wed|thu|fri|sat|sun)`)
	rangeMatch := dayRangePattern.FindStringSubmatch(hoursText)

	if len(rangeMatch) == 3 {
		return checkDayRange(rangeMatch, currentWeekday)
	}

	return true
}

func checkDayRange(rangeMatch []string, currentWeekday int) bool {
	startDay := getDayNumber(rangeMatch[1])
	endDay := getDayNumber(rangeMatch[2])

	if startDay <= endDay {
		return currentWeekday >= startDay && currentWeekday <= endDay
	}
	return currentWeekday >= startDay || currentWeekday <= endDay
}

func isTimeOpen(hoursText string, currentTime time.Time) bool {
	timePattern := regexp.MustCompile(`(\d{1,2})(?::(\d{2}))?\s*(am|pm)?`)
	matches := timePattern.FindAllStringSubmatch(hoursText, -1)

	if len(matches) < 2 {
		return false
	}

	openMinutes, closeMinutes, ok := extractOpenCloseMinutes(matches)
	if !ok {
		return false
	}

	currentMinutes := currentTime.Hour()*60 + currentTime.Minute()

	if closeMinutes < openMinutes {
		return currentMinutes >= openMinutes || currentMinutes < closeMinutes
	}

	return currentMinutes >= openMinutes && currentMinutes < closeMinutes
}

func extractOpenCloseMinutes(matches [][]string) (int, int, bool) {
	openMatch := matches[0]
	closeMatch := matches[1]

	openMinutes, ok1 := parseTimeToMinutes(openMatch[1], openMatch[2], openMatch[3])
	closeMinutes, ok2 := parseTimeToMinutes(closeMatch[1], closeMatch[2], closeMatch[3])

	if !ok1 || !ok2 {
		return 0, 0, false
	}

	if openMatch[3] == "" && closeMatch[3] == "pm" {
		openHour, _ := strconv.Atoi(openMatch[1])
		closeHour, _ := strconv.Atoi(closeMatch[1])
		if openHour > closeHour {
			openMinutes += 0
		} else if openHour < closeHour {
			if openHour < 12 {
				openMinutes += 12 * 60
			}
		}
	}
	return openMinutes, closeMinutes, true
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

func evaluateStructuredHours(structured string, currentTime time.Time) (bool, bool) {
	if structured == "" {
		return false, false
	}

	var schedule map[string][]string
	if err := json.Unmarshal([]byte(structured), &schedule); err != nil {
		return false, false // invalid JSON, fallback to regex
	}

	dayNames := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}
	weekday := strings.ToLower(dayNames[currentTime.Weekday()])

	ranges, ok := schedule[weekday]
	if !ok {
		return false, false // day not present in JSON, fallback to regex
	}

	if len(ranges) == 0 {
		return false, true // explicitly closed today
	}

	currentMinutes := currentTime.Hour()*60 + currentTime.Minute()

	for _, r := range ranges {
		if isTimeInRange(r, currentMinutes) {
			return true, true
		}
	}

	return false, true // was evaluated and determined to be closed
}

func isTimeInRange(timeRange string, currentMinutes int) bool {
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return false
	}
	openMin, ok1 := parseTimeStr(parts[0])
	closeMin, ok2 := parseTimeStr(parts[1])
	if !ok1 || !ok2 {
		return false
	}

	if closeMin < openMin { // overlaps midnight
		return currentMinutes >= openMin || currentMinutes < closeMin
	}
	return currentMinutes >= openMin && currentMinutes < closeMin
}

func parseTimeStr(t string) (int, bool) {
	parts := strings.Split(strings.TrimSpace(t), ":")
	if len(parts) != 2 {
		return 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, false
	}
	return h*60 + m, true
}
