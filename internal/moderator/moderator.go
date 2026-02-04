package moderator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/jadecobra/agbalumo/internal/domain"
	"google.golang.org/api/option"
)

var ErrContentViolation = errors.New("content violates cultural relevancy standards")

// Moderator defines the interface for content moderation.
type Moderator interface {
	CheckListing(ctx context.Context, listing domain.Listing) error
}

// GenAIModel abstracts the content generation for testing
type GenAIModel interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

// Check if *genai.GenerativeModel implements interface (runtime check)
// Actually *genai.GenerativeModel struct does not implicitly implement our interface if the methods don't match exactly or if we don't wrap it.
// .GenerateContent takes (ctx, parts...)
// So we can use the interface directly if the signature matches.

// GeminiModerator uses Google's Gemini Pro model to moderate content.
type GeminiModerator struct {
	client *genai.Client
	model  GenAIModel
}

// NewGeminiModerator creates a new instance of GeminiModerator.
// It requires the GEMINI_API_KEY environment variable to be set.
func NewGeminiModerator(ctx context.Context) (*GeminiModerator, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("[Moderator] Warning: GEMINI_API_KEY not set. Moderation will be skipped (Stub Mode).")
		return nil, nil // Return nil, nil to indicate stub mode fallback logic if handled by caller, or just handle gracefully
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	return NewWithModel(client, model), nil
}

// NewWithModel creates a moderator with a specific model (useful for testing)
func NewWithModel(client *genai.Client, model GenAIModel) *GeminiModerator {
	return &GeminiModerator{client: client, model: model}
}

// CheckListing sends the listing content to Gemini to check for cultural relevancy.
func (m *GeminiModerator) CheckListing(ctx context.Context, listing domain.Listing) error {
	if m == nil || m.model == nil {
		log.Println("[Moderator] Stub Check: GEMINI_API_KEY was missing. Content allowed.")
		return nil
	}

	prompt := fmt.Sprintf(`
You are a content moderator for 'Agbalumo', a directory for the West African diaspora. 
Your job is to ensure listings are culturally relevant, safe, and appropriate.

Rule: The listing should be relevant to the West African community (Businesses, Services, Products, Requests).
Rule: No hate speech, scams, or illegal content.

Evaluate this listing:
Title: %s
Description: %s
Origin: %s

Respond with exactly one word: "PERMIT" or "DENY".
`, listing.Title, listing.Description, listing.OwnerOrigin)

	resp, err := m.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("[Moderator] Error calling Gemini API: %v", err)
		return nil // Fail open for now, or return err? For MVP, fail open.
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil
	}

	log.Printf("[Moderator] Inspecting %d parts", len(resp.Candidates[0].Content.Parts))

	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			decision := strings.TrimSpace(strings.ToUpper(string(txt)))
			log.Printf("[Moderator] Decision for '%s': %s", listing.Title, decision)
			if strings.Contains(decision, "DENY") {
				return ErrContentViolation
			}
		}
	}

	return nil
}
