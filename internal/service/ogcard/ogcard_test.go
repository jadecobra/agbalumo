package ogcard_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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
	outDir := "/Users/johnnyblase/.gemini/antigravity/brain/46806a74-a1cc-4cac-a9ae-ea74f611ac32"
	svc := ogcard.NewService()

	asafoData := ogcard.CardData{
		Title:       "Asafo Market",
		City:        "Grand Prairie",
		Category:    "Food",
		Rating:      4.4,
		ReviewCount: 439,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	asafoBytes, err := svc.Generate(asafoData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_asafo.png"), asafoBytes, 0600)

	lolaData := ogcard.CardData{
		Title:       "Lola's Restaurant & Lounge",
		City:        "Irving",
		Category:    "Food",
		Rating:      4.1,
		ReviewCount: 579,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	lolaBytes, err := svc.Generate(lolaData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_lolas.png"), lolaBytes, 0600)

	adomData := ogcard.CardData{
		Title:       "Adom African Market",
		City:        "Arlington",
		Category:    "Food",
		Rating:      0.0,
		ReviewCount: 0,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	adomBytes, err := svc.Generate(adomData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_adom.png"), adomBytes, 0600)

	twoLineData := ogcard.CardData{
		Title:       "Mama Put Authentic Nigerian Kitchen & Suya Spot",
		City:        "Dallas",
		Category:    "Restaurant",
		Rating:      4.9,
		ReviewCount: 312,
		Tagline:     "Find African food in under 60 seconds.",
		UpdatedAt:   time.Now(),
	}
	twoLineBytes, err := svc.Generate(twoLineData)
	require.NoError(t, err)
	_ = os.WriteFile(filepath.Join(outDir, "og_card_twoline.png"), twoLineBytes, 0600)

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

func TestService_Generate_LongNameStaysInsideCard(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()

	titles := []string{
		"Mama Put Authentic Nigerian Kitchen & Suya Spot",
		"Extremely Long African Market And Restaurant Lounge With Extensive Name That Spans Far",
		"SupercalifragilisticexpialidociousLongNameWithoutAnySpacesInBetweenWordsHere",
	}

	for _, title := range titles {
		data := ogcard.CardData{
			Title:       title,
			City:        "Grand Prairie",
			Category:    "Food",
			Rating:      4.5,
			ReviewCount: 120,
			UpdatedAt:   time.Now(),
		}

		pngBytes, err := svc.Generate(data)
		require.NoError(t, err)
		require.NotEmpty(t, pngBytes)

		img, err := png.Decode(bytes.NewReader(pngBytes))
		require.NoError(t, err)

		bounds := img.Bounds()
		assert.Equal(t, 1200, bounds.Dx())
		assert.Equal(t, 630, bounds.Dy())

		bg := color.RGBA{R: 0x1A, G: 0x12, B: 0x0E, A: 0xFF}
		for y := 150; y < 510; y++ {
			for x := 1165; x < 1195; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				expectedR, expectedG, expectedB, expectedA := bg.RGBA()
				assert.Equal(t, expectedR, r, "pixel at (%d, %d) must not bleed past right margin for title %q", x, y, title)
				assert.Equal(t, expectedG, g)
				assert.Equal(t, expectedB, b)
				assert.Equal(t, expectedA, a)
			}
		}
	}
}

func TestService_Generate_NoRating(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()

	data := ogcard.CardData{
		Title:       "Adom African Market",
		City:        "Arlington",
		Category:    "Food",
		Rating:      0.0,
		ReviewCount: 0,
		UpdatedAt:   time.Now(),
	}

	pngBytes, err := svc.Generate(data)
	require.NoError(t, err)
	require.NotEmpty(t, pngBytes)

	cfg, format, err := image.DecodeConfig(bytes.NewReader(pngBytes))
	require.NoError(t, err)
	assert.Equal(t, "png", format)
	assert.Equal(t, 1200, cfg.Width)
	assert.Equal(t, 630, cfg.Height)
}

func TestService_Generate_FoodLabelAndFooterAlignment(t *testing.T) {
	t.Parallel()

	svc := ogcard.NewService()

	data := ogcard.CardData{
		Title:       "Asafo Market",
		City:        "Grand Prairie",
		Category:    "Food",
		Rating:      4.4,
		ReviewCount: 439,
		UpdatedAt:   time.Now(),
	}

	pngBytes, err := svc.Generate(data)
	require.NoError(t, err)

	img, err := png.Decode(bytes.NewReader(pngBytes))
	require.NoError(t, err)

	cardBg := color.RGBA{R: 0x24, G: 0x1B, B: 0x16, A: 0xFF}

	for y := 95; y <= 108; y++ {
		for x := 1128; x <= 1150; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			expR, expG, expB, expA := cardBg.RGBA()
			assert.Equal(t, expR, r, "badge pixel at (%d, %d) must not extend past 75px margin", x, y)
			assert.Equal(t, expG, g)
			assert.Equal(t, expB, b)
			assert.Equal(t, expA, a)
		}
	}

	for y := 550; y <= 565; y++ {
		for x := 1128; x <= 1150; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			expR, expG, expB, expA := cardBg.RGBA()
			assert.Equal(t, expR, r, "footer pixel at (%d, %d) must not extend past 75px margin", x, y)
			assert.Equal(t, expG, g)
			assert.Equal(t, expB, b)
			assert.Equal(t, expA, a)
		}
	}
}
