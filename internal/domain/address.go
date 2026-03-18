package domain

import "strings"

// ExtractCityFromAddress tries to guess the city from an address string.
// Assumes formats like "Street, City, Country" or "City, Country".
func ExtractCityFromAddress(address string) string {
	parts := strings.Split(address, ",")
	// If 3 parts, middle is usually city
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[1])
	}
	// If 2 parts, first is usually city
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0])
	}
	// Fallback to the whole string if it's short, or just return it
	return strings.TrimSpace(address)
}
