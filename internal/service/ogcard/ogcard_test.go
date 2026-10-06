package ogcard_test

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/service/ogcard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Generate_MagicBytesAndDimensions(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()

	data := ogcard.CardData{
		Title:       "Aria Suya Kitchen",
		City:        "Arlington",
		Category:    "Restaurant",
		Rating:      4.8,
		ReviewCount: 124,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}

	pngBytes, err := svc.Generate(data)
	require.NoError(t, err)
	require.NotEmpty(t, pngBytes)

	// Magic bytes check (\x89PNG\r\n\x1a\n)
	expectedMagic := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	require.True(t, len(pngBytes) >= 8)
	assert.Equal(t, expectedMagic, pngBytes[:8])

	// Decode dimensions
	cfg, format, err := image.DecodeConfig(bytes.NewReader(pngBytes))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 1200, cfg.Width)
	assert.Equal(t, 630, cfg.Height)
}

func TestService_GenerateFallback_MagicBytesAndDimensions(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()

	pngBytes := svc.GenerateFallback()
	require.NotEmpty(t, pngBytes)

	expectedMagic := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	require.True(t, len(pngBytes) >= 8)
	assert.Equal(t, expectedMagic, pngBytes[:8])

	cfg, format, err := image.DecodeConfig(bytes.NewReader(pngBytes))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 1200, cfg.Width)
	assert.Equal(t, 630, cfg.Height)
}

func TestService_Generate_WarmPerformance(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()
	data := ogcard.CardData{
		Title:       "Aria Suya Kitchen",
		City:        "Dallas",
		Rating:      4.9,
		ReviewCount: 50,
	}

	// Warm run
	_, _ = svc.Generate(data)

	start := time.Now()
	_, err := svc.Generate(data)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Less(t, duration, 500*time.Millisecond, "Warm rendering must execute well under 500ms budget")
}
