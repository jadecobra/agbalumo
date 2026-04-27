package service

import (
	"time"
)

// ComputeIsOpen determines if a listing is currently open based on its unstructured hours text.
func ComputeIsOpen(hoursText string, currentTime time.Time) bool {
	return false
}
