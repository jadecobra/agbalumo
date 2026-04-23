package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// GoogleGeocodingService implements domain.GeocodingService using Google Maps API.
type GoogleGeocodingService struct {
	APIKey  string // #nosec G117 - Service field for API key injection
	BaseURL string
}

func NewGoogleGeocodingService(apiKey string) *GoogleGeocodingService {
	return &GoogleGeocodingService{
		APIKey:  apiKey,
		BaseURL: "https://maps.googleapis.com/maps/api/geocode/json",
	}
}

type geocodingResponse struct {
	Status  string `json:"status"`
	Results []struct {
		AddressComponents []struct {
			LongName string   `json:"long_name"`
			Types    []string `json:"types"`
		} `json:"address_components"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func (s *GoogleGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	res, err := s.fetchGeocode(ctx, address)
	if err != nil {
		return "", err
	}
	if len(res.Results) == 0 {
		return "", nil
	}
	return s.extractCity(res.Results[0].AddressComponents), nil
}

func (s *GoogleGeocodingService) Geocode(ctx context.Context, address string) (float64, float64, error) {
	res, err := s.fetchGeocode(ctx, address)
	if err != nil {
		return 0, 0, err
	}
	if len(res.Results) == 0 {
		return 0, 0, nil
	}
	loc := res.Results[0].Geometry.Location
	return loc.Lat, loc.Lng, nil
}

func (s *GoogleGeocodingService) fetchGeocode(ctx context.Context, address string) (*geocodingResponse, error) {
	if s.APIKey == "" {
		// Development Fallback: Hardcoded coordinates for common cities to allow demo without API key.
		lowCity := strings.ToLower(address)
		if strings.Contains(lowCity, "dallas") {
			resp := &geocodingResponse{Status: "OK"}
			resp.Results = make([]struct {
				AddressComponents []struct {
					LongName string   `json:"long_name"`
					Types    []string `json:"types"`
				} `json:"address_components"`
				Geometry struct {
					Location struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"location"`
				} `json:"geometry"`
			}, 1)
			resp.Results[0].Geometry.Location.Lat = 32.7767
			resp.Results[0].Geometry.Location.Lng = -96.7970
			return resp, nil
		}
		return nil, fmt.Errorf("google maps api key is not configured")
	}

	apiURL, err := s.buildURL(address)
	if err != nil {
		return nil, err
	}

	body, err := s.fetch(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var res geocodingResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Status != "OK" {
		if res.Status == "ZERO_RESULTS" {
			return &res, nil
		}
		return nil, fmt.Errorf("geocoding api error status: %s", res.Status)
	}

	return &res, nil
}

func (s *GoogleGeocodingService) buildURL(address string) (string, error) {
	baseURL, err := url.Parse(s.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid geocoding base url: %w", err)
	}

	q := baseURL.Query()
	q.Set(domain.FieldAddress, address)
	q.Set("key", s.APIKey)
	baseURL.RawQuery = q.Encode()
	return baseURL.String(), nil
}

func (s *GoogleGeocodingService) fetch(ctx context.Context, apiURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// #nosec G107 G704 - SSRF check: BaseURL is constant, address is encoded.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding api returned status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}



func (s *GoogleGeocodingService) extractCity(components []struct {
	LongName string   `json:"long_name"`
	Types    []string `json:"types"`
}) string {
	var city string
	for _, component := range components {
		types := component.Types
		if contains(types, "locality") {
			return component.LongName
		}
		if contains(types, "sublocality_level_1") || contains(types, "sublocality") {
			city = component.LongName
		} else if city == "" && (contains(types, "postal_town") || contains(types, "administrative_area_level_2") || contains(types, "neighborhood")) {
			city = component.LongName
		}
	}
	return city
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
