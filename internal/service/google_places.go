package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GooglePlacesClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey:  apiKey,
		client:  &http.Client{},
		baseURL: "https://places.googleapis.com",
	}
}

func (c *GooglePlacesClient) SetBaseURL(url string) {
	c.baseURL = url
}

type PlacesMetrics struct {
	Rating      float64
	ReviewCount int
}

type textSearchRequest struct {
	TextQuery string `json:"textQuery"`
}

type placeResponse struct {
	Places []struct {
		Rating          float64 `json:"rating"`
		UserRatingCount int     `json:"userRatingCount"`
	} `json:"places"`
}

func (c *GooglePlacesClient) FetchMetrics(ctx context.Context, title, city string) (PlacesMetrics, error) {
	if c.apiKey == "" {
		return PlacesMetrics{}, fmt.Errorf("Google Places API key is empty")
	}

	query := title
	if city != "" {
		query += ", " + city
	}

	reqBody := textSearchRequest{TextQuery: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return PlacesMetrics{}, err
	}

	url := c.baseURL + "/v1/places:searchText"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return PlacesMetrics{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.rating,places.userRatingCount")

	resp, err := c.client.Do(req)
	if err != nil {
		return PlacesMetrics{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return PlacesMetrics{}, fmt.Errorf("Places API request failed with status %d", resp.StatusCode)
	}

	var apiResp placeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return PlacesMetrics{}, err
	}

	if len(apiResp.Places) == 0 {
		return PlacesMetrics{}, fmt.Errorf("no places found matching query")
	}

	return PlacesMetrics{
		Rating:      apiResp.Places[0].Rating,
		ReviewCount: apiResp.Places[0].UserRatingCount,
	}, nil
}

