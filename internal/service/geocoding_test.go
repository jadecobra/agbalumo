package service

import (
	"context"
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
		responseStatus int
		responseBody   string
		expectedCity   string
		expectedError  bool
	}{
		{
			name:           "Success",
			address:        "1600 Amphitheatre Parkway, Mountain View, CA",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "Mountain View", "types": ["locality", "political"]},
						{"long_name": "Santa Clara County", "types": ["administrative_area_level_2", "political"]}
					]
				}]
			}`,
			expectedCity: "Mountain View",
		},
		{
			name:           "Sublocality Fallback",
			address:        "Some Sublocality",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "Sublocality Info", "types": ["sublocality", "political"]}
					]
				}]
			}`,
			expectedCity: "Sublocality Info",
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
			responseBody:   `{"status": "ZERO_RESULTS", "results": []}`,
			expectedCity:   "",
		},
		{
			name:           "Postal Town Fallback",
			address:        "UK Address",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "London Town", "types": ["postal_town"]}
					]
				}]
			}`,
			expectedCity: "London Town",
		},
		{
			name:           "Neighborhood Fallback",
			address:        "NY Address",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "Brooklyn", "types": ["neighborhood"]}
					]
				}]
			}`,
			expectedCity: "Brooklyn",
		},
		{
			name:           "Admin Area Level 2 Fallback",
			address:        "County Level",
			apiKey:         "valid-key",
			responseStatus: http.StatusOK,
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "Some County", "types": ["administrative_area_level_2"]}
					]
				}]
			}`,
			expectedCity: "Some County",
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
			responseBody: `{
				"status": "OK",
				"results": [{
					"address_components": [
						{"long_name": "USA", "types": ["country"]}
					]
				}]
			}`,
			expectedCity: "",
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
