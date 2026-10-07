package ogcard_test

import (
	"bytes"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
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

func TestService_ExportArtifacts(t *testing.T) {
	outDir := "/Users/johnnyblase/.gemini/antigravity/brain/b9436bc0-55c6-4582-aa73-286f63a9c27c"
	svc := ogcard.NewService()

	sampleData := ogcard.CardData{
		Title:       "Aria Suya Kitchen",
		City:        "Arlington, TX",
		Category:    "West African Restaurant",
		Rating:      4.8,
		ReviewCount: 124,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	sampleBytes, err := svc.Generate(sampleData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_sample.png"), sampleBytes, 0600)

	longData := ogcard.CardData{
		Title:       "Mama Put Authentic Nigerian Kitchen & Suya Spot",
		City:        "Dallas, TX",
		Category:    "Restaurant",
		Rating:      4.9,
		ReviewCount: 312,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	longBytes, err := svc.Generate(longData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_long_title.png"), longBytes, 0600)

	fallbackBytes := svc.GenerateFallback()
	_ = os.WriteFile(filepath.Join(outDir, "og_card_fallback.png"), fallbackBytes, 0600)
}

func BenchmarkService_Generate_Cold(b *testing.B) {
	svc := ogcard.NewService()
	data := ogcard.CardData{
		Title:       "Aria Suya Kitchen",
		City:        "Dallas, TX",
		Category:    "Restaurant",
		Rating:      4.9,
		ReviewCount: 50,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.UpdatedAt = time.Unix(int64(i+1), 0)
		_, err := svc.Generate(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkService_Generate_Warm(b *testing.B) {
	svc := ogcard.NewService()
	data := ogcard.CardData{
		Title:       "Aria Suya Kitchen",
		City:        "Dallas, TX",
		Category:    "Restaurant",
		Rating:      4.9,
		ReviewCount: 50,
	}
	_, _ = svc.Generate(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Generate(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
