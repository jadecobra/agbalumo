package moderator_test

import (
	"context"
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/moderator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeminiModerator_Integration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	mod, err := moderator.NewGeminiModerator(ctx)
	require.NoError(t, err)

	// Case 1: Permitted Content
	goodListing := domain.Listing{
		Title:       "African Grocery Store",
		Description: "Selling yam, cassava, and plantain.",
		OwnerOrigin: "Nigeria",
	}
	err = mod.CheckListing(ctx, goodListing)
	assert.NoError(t, err, "Should permit valid content")

	// Case 2: Denied Content (Simulated bad content)
	// We need something that Gemini will definitely reject based on the prompt "No hate speech, scams, or illegal content".
	// Let's use a blatant scam example.
	badListing := domain.Listing{
		Title:       "Double your money instantly",
		Description: "Send me bitcoin and I will return 2x in 1 hour. This is a guaranteed investment scam.",
		OwnerOrigin: "Unknown",
	}
	err = mod.CheckListing(ctx, badListing)
	assert.ErrorIs(t, err, moderator.ErrContentViolation, "Should deny scam content")
}
