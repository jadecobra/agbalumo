package domain

import "strings"

// ExtractCityFromAddress tries to guess the city from an address string.
// Assumes formats like "Street, City, ST Zip" or "City, Country".
func ExtractCityFromAddress(address string) string {
	parts := strings.Split(address, ",")
	if len(parts) >= 2 {
		// If 3 parts: Street, City, State ZIP -> City is parts[1]
		// If 2 parts: City, State ZIP -> City is parts[0]
		if len(parts) >= 3 {
			return strings.TrimSpace(parts[1])
		}
		return strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(address)
}

// ExtractStateFromAddress tries to pull the state code from an address.
// Handles "City, TX 75201" or "City, Texas, USA"
func ExtractStateFromAddress(address string) string {
	parts := strings.Split(address, ",")
	if len(parts) < 2 {
		return ""
	}

	// Usually the state is in the part after the city
	segment := ""
	if len(parts) >= 3 {
		segment = strings.TrimSpace(parts[2]) // Street, City, STATE ZIP
	} else {
		segment = strings.TrimSpace(parts[1]) // City, STATE (might be country too)
	}

	// Split by space to remove ZIP code if present
	subParts := strings.Fields(segment)
	if len(subParts) > 0 {
		state := subParts[0]
		// If it's a 2-letter code or looks like a state name, return it
		if len(state) == 2 || len(state) > 3 {
			return state
		}
	}

	return ""
}

// ExtractCountryFromAddress tries to pull the country from an address.
func ExtractCountryFromAddress(address string) string {
	parts := strings.Split(address, ",")
	if len(parts) == 0 {
		return CountryUSA // Default
	}
	last := strings.TrimSpace(parts[len(parts)-1])
	// If last part is a zip code like sequence, it's likely USA
	if len(last) == 5 || (len(last) > 5 && strings.Contains(last, " ")) {
		return CountryUSA
	}

	if last == "" {
		return CountryUSA
	}
	return last
}
