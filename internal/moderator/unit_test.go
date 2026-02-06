package moderator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/moderator"
)

// MockGenAIModel helps verify moderator logic without API calls
type MockGenAIModel struct {
	Response *genai.GenerateContentResponse
	Err      error
}

func (m *MockGenAIModel) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	return m.Response, m.Err
}

func TestCheckListing_Unit(t *testing.T) {
	ctx := context.Background()

	t.Run("Permits Content", func(t *testing.T) {
		// Mock Response saying "PERMIT"
		mockResp := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{
							genai.Text("PERMIT"),
						},
					},
				},
			},
		}

		mod := moderator.NewWithModel(nil, &MockGenAIModel{Response: mockResp})

		listing := domain.Listing{Title: "Good Listing"}
		err := mod.CheckListing(ctx, listing)
		if err != nil {
			t.Errorf("Expected nil error for PERMIT, got %v", err)
		}
	})

	t.Run("Denies Content", func(t *testing.T) {
		// Mock Response saying "DENY"
		mockResp := &genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{
				{
					Content: &genai.Content{
						Parts: []genai.Part{
							genai.Text("DENY"),
						},
					},
				},
			},
		}

		mod := moderator.NewWithModel(nil, &MockGenAIModel{Response: mockResp})

		listing := domain.Listing{Title: "Bad Listing"}
		err := mod.CheckListing(ctx, listing)
		if !errors.Is(err, moderator.ErrContentViolation) {
			t.Errorf("Expected ErrContentViolation for DENY, got %v", err)
		}
	})

	t.Run("API Error Handling", func(t *testing.T) {
		mod := moderator.NewWithModel(nil, &MockGenAIModel{Err: errors.New("API failure")})

		listing := domain.Listing{Title: "Listing"}
		// Logic handles error by logging and returning nil (failing open)
		err := mod.CheckListing(ctx, listing)
		if err != nil {
			t.Errorf("Expected fail-open behavior (nil error), got %v", err)
		}
	})
}

func TestNewGeminiModerator_EnvVar(t *testing.T) {
	t.Run("Missing API Key", func(t *testing.T) {
		t.Setenv("GEMINI_API_KEY", "")
		ctx := context.Background()
		mod, err := moderator.NewGeminiModerator(ctx)
		if err != nil {
			t.Fatalf("Expected nil error for missing key, got %v", err)
		}
		if mod != nil {
			t.Error("Expected nil moderator for missing key")
		}
	})

	t.Run("With API Key", func(t *testing.T) {
		t.Setenv("GEMINI_API_KEY", "dummy-key")
		ctx := context.Background()
		mod, err := moderator.NewGeminiModerator(ctx)
		if err != nil {
			t.Fatalf("Expected check to pass (or at least try client creation), got %v", err)
		}
		if mod == nil {
			t.Error("Expected moderator instance")
		}
	})
}
