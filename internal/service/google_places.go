package service

import (
	"context"
	"net/http"
)

type GooglePlacesClient struct {
	client *http.Client
	apiKey string
}

func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type PlacesMetrics struct {
	Rating      float64
	ReviewCount int
}

func (c *GooglePlacesClient) FetchMetrics(ctx context.Context, title, city string) (PlacesMetrics, error) {
	return PlacesMetrics{Rating: 0.0, ReviewCount: 0}, nil
}
