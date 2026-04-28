package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestGeminiHoursExtractor_ExtractHours(t *testing.T) {
	tests := []struct {
		name         string
		rawHours     string
		mockResponse string
		wantJSON     string
		expectErr    bool
	}{
		{
			name:     "valid simple hours",
			rawHours: "Mon-Fri 9am-5pm",
			mockResponse: `{
				"candidates": [{
					"content": {
						"parts": [{
							"text": "{\"mon\": [\"09:00-17:00\"], \"tue\": [\"09:00-17:00\"], \"wed\": [\"09:00-17:00\"], \"thu\": [\"09:00-17:00\"], \"fri\": [\"09:00-17:00\"]}"
						}]
					}
				}]
			}`,
			wantJSON:  `{"mon": ["09:00-17:00"], "tue": ["09:00-17:00"], "wed": ["09:00-17:00"], "thu": ["09:00-17:00"], "fri": ["09:00-17:00"]}`,
			expectErr: false,
		},
		{
			name:     "complex split hours",
			rawHours: "Mon-Fri 9-5, Sat 10-2",
			mockResponse: `{
				"candidates": [{
					"content": {
						"parts": [{
							"text": "{\"mon\": [\"09:00-17:00\"], \"sat\": [\"10:00-14:00\"]}"
						}]
					}
				}]
			}`,
			wantJSON:  `{"mon": ["09:00-17:00"], "sat": ["10:00-14:00"]}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &http.Client{
				Transport: &MockTransport{
					RoundTripFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(tt.mockResponse)),
							Header:     make(http.Header),
						}, nil
					},
				},
			}

			extractor := NewGeminiHoursExtractor("test-api-key", mockClient)
			got, err := extractor.ExtractHours(context.Background(), tt.rawHours)

			if (err != nil) != tt.expectErr {
				t.Fatalf("ExtractHours() error = %v, expectErr %v", err, tt.expectErr)
			}

			if !tt.expectErr && got != tt.wantJSON {
				t.Errorf("ExtractHours() = %q, want %q", got, tt.wantJSON)
			}
		})
	}
}
