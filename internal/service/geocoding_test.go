package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleGeocodingService_GetCity(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		apiKey         string
		responseBody   string
		expectedCity   string
		responseStatus int
		expectedError  bool
	}{
		{
			name:           "Success",
			address:        "1600 Amphitheatre Parkway, Mountain View, CA",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("Mountain View", "locality", "political"), comp("Santa Clara County", "administrative_area_level_2", "political")),
			expectedCity:   "Mountain View",
		},
		{
			name:           "Sublocality Fallback",
			address:        "Some Sublocality",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("Sublocality Info", "sublocality", "political")),
			expectedCity:   "Sublocality Info",
		},
		{
			name:          "Empty API Key",
			address:       "Any",
			apiKey:        "",
			expectedError: true,
		},
		{
			name:           "Zero Results",
			address:        "Middle of nowhere",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("ZERO_RESULTS"),
			expectedCity:   "",
		},
		{
			name:           "Postal Town Fallback",
			address:        "UK Address",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("London Town", "postal_town")),
			expectedCity:   "London Town",
		},
		{
			name:           "Neighborhood Fallback",
			address:        "NY Address",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("Brooklyn", "neighborhood")),
			expectedCity:   "Brooklyn",
		},
		{
			name:           "Admin Area Level 2 Fallback",
			address:        "County Level",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("Some County", "administrative_area_level_2")),
			expectedCity:   "Some County",
		},
		{
			name:           "API Error Status",
			address:        "Any",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   `{"status": "REQUEST_DENIED", "results": []}`,
			expectedError:  true,
		},
		{
			name:           "HTTP Error",
			address:        "Any",
			apiKey:         "valid-key",
			responseStatus: http.StatusForbidden,
			expectedError:  true,
		},
		{
			name:           "Invalid JSON",
			address:        "Any",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   "invalid",
			expectedError:  true,
		},
		{
			name:           "OK Status Empty Results",
			address:        "Any",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   `{"status": "OK", "results": []}`,
			expectedCity:   "",
		},
		{
			name:           "No City Components",
			address:        "Any",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody:   geocodeResponse("OK", comp("USA", "country")),
			expectedCity:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			svc := NewGoogleGeocodingService(tt.apiKey)
			svc.BaseURL = server.URL

			city, err := svc.GetCity(context.Background(), tt.address)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCity, city)
			}
		})
	}
}

func geocodeResponse(status string, components ...string) string {
	comps := ""
	for i, c := range components {
		if i > 0 {
			comps += ","
		}
		comps += c
	}
	if comps == "" {
		return fmt.Sprintf(`{"status": "%s", "results": []}`, status)
	}
	return fmt.Sprintf(`{"status": "%s", "results": [{"address_components": [%s]}]}`, status, comps)
}

func comp(name string, types ...string) string {
	quotedTypes := ""
	for i, t := range types {
		if i > 0 {
			quotedTypes += ","
		}
		quotedTypes += fmt.Sprintf(`"%s"`, t)
	}
	return fmt.Sprintf(`{"long_name": "%s", "types": [%s]}`, name, quotedTypes)
}
