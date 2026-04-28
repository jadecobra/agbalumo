package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type geminiRequest struct {
	GenerationConfig *geminiConfig   `json:"generationConfig,omitempty"`
	Contents         []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiConfig struct {
	ResponseMimeType string `json:"responseMimeType"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (e *GeminiHoursExtractor) ExtractHours(ctx context.Context, rawHours string) (string, error) {
	if rawHours == "" {
		return "{}", nil
	}

	prompt := fmt.Sprintf(`Parse the following operating hours text into a structured JSON object. 
The keys must be the lowercase 3-letter day abbreviations: "mon", "tue", "wed", "thu", "fri", "sat", "sun".
The values must be an array of time ranges in 24-hour HH:MM-HH:MM format, e.g. ["09:00-17:00", "18:00-22:00"].
If the store is closed on a given day, use an empty array.
Return ONLY valid raw JSON with no markdown wrapping.

Text: %s`, rawHours)

	reqPayload := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: &geminiConfig{
			ResponseMimeType: "application/json",
		},
	}

	jsonBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", e.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "{}", nil
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
