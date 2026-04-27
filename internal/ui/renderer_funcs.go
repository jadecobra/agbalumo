package ui

import (
	"encoding/json"
	"errors"
	"html/template"
	"strings"
	"time"
)

func seq(start, end int) []int {
	var s []int
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	d := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		d[key] = values[i+1]
	}
	return d, nil
}

func toJson(v interface{}) (template.JS, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	// #nosec G203 - Intentional template escape for trusted content
	return template.JS(b), nil
}

func isNew(createdAt time.Time) bool {
	if createdAt.IsZero() {
		return false
	}
	return time.Since(createdAt) < 7*24*time.Hour
}

func safeHTML(s string) template.HTML {
	// #nosec G203 - Intentional template escape for trusted content
	return template.HTML(s)
}

func safeHTMLAttr(s string) template.HTMLAttr {
	// #nosec G203 - Intentional template escape for trusted content
	return template.HTMLAttr(s)
}

func safeJS(s string) template.JS {
	// #nosec G203 - Intentional template escape for trusted content
	return template.JS(s)
}

func displayCity(city, address string) string {
	if city != "" {
		return city
	}
	if address == "" {
		return ""
	}
	// Fallback: extract city from address (e.g., "123 Main St, City, ST 12345")
	parts := strings.Split(address, ",")
	if len(parts) >= 2 {
		// Most common format is [Street], [City], [State Zip]
		// We take the second part if available
		return strings.TrimSpace(parts[1])
	}
	// If no comma exists, we don't know if it's a city or a street.
	// To avoid showing "123 Test St" in the city slot, we return empty.
	return ""
}

func fallbackImageURL(imageURL, websiteURL string) string {
	return ""
}

