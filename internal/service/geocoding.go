package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GoogleGeocodingService implements domain.GeocodingService using Google Maps API.
type GoogleGeocodingService struct {
	APIKey  string
	BaseURL string
}

func NewGoogleGeocodingService(apiKey string) *GoogleGeocodingService {
	return &GoogleGeocodingService{
		APIKey:  apiKey,
		BaseURL: "https://maps.googleapis.com/maps/api/geocode/json",
	}
}

type geocodingResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName string   `json:"long_name"`
			Types    []string `json:"types"`
		} `json:"address_components"`
	} `json:"results"`
	Status string `json:"status"`
}

func (s *GoogleGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	if s.APIKey == "" {
		return "", fmt.Errorf("google maps api key is not configured")
	}

	u := fmt.Sprintf("%s?address=%s&key=%s",
		s.BaseURL, url.QueryEscape(address), s.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("geocoding api returned status: %d", resp.StatusCode)
	}

	var res geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Status != "OK" {
		if res.Status == "ZERO_RESULTS" {
			return "", nil
		}
		return "", fmt.Errorf("geocoding api error status: %s", res.Status)
	}

	if len(res.Results) == 0 {
		return "", nil
	}

	// Logic similar to app.js
	var city string
	for _, component := range res.Results[0].AddressComponents {
		types := component.Types
		if contains(types, "locality") {
			city = component.LongName
			break
		} else if contains(types, "sublocality_level_1") || contains(types, "sublocality") {
			city = component.LongName
		} else if city == "" && (contains(types, "postal_town") || contains(types, "administrative_area_level_2") || contains(types, "neighborhood")) {
			city = component.LongName
		}
	}

	return city, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
