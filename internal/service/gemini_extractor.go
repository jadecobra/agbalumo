package service

import (
	"context"
	"net/http"
)

type GeminiHoursExtractor struct {
	client *http.Client
	apiKey string
}


func NewGeminiHoursExtractor(apiKey string, client *http.Client) *GeminiHoursExtractor {
	if client == nil {
		client = http.DefaultClient
	}
	return &GeminiHoursExtractor{
		apiKey: apiKey,
		client: client,
	}
}

func (e *GeminiHoursExtractor) ExtractHours(ctx context.Context, rawHours string) (string, error) {
	return "", nil
}
