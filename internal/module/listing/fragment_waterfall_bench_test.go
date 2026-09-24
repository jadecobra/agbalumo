package listing_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func BenchmarkHandleFragment_QueryWaterfall(b *testing.B) {
	t := &testing.T{}
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.SeedStandardData(t)
	testutil.SeedAdaDallasData(t, env.App.DB)

	h := listing.NewListingHandler(env.App)

	b.Run("Page1_WithSearch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c, _ := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?q=Nigerian&type=Food&page=1", nil)
			_ = h.HandleFragment(c)
		}
	})

	b.Run("Page1_NoSearch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c, _ := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?type=Food&page=1", nil)
			_ = h.HandleFragment(c)
		}
	})

	b.Run("Page2_NoFeatured", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c, _ := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?type=Food&page=2", nil)
			_ = h.HandleFragment(c)
		}
	})
}

func TestHandleFragment_LatencyBaseline(t *testing.T) {
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	env.SeedStandardData(t)
	testutil.SeedAdaDallasData(t, env.App.DB)

	h := listing.NewListingHandler(env.App)

	const iterations = 100
	var totalNs int64
	for i := 0; i < iterations; i++ {
		c, rec := testutil.SetupModuleContext(http.MethodGet, "/listings/fragment?q=Nigerian&type=Food&page=1", nil)
		start := time.Now()
		_ = h.HandleFragment(c)
		totalNs += time.Since(start).Nanoseconds()
		if rec.Code != 0 && rec.Code != 200 {
			t.Fatalf("unexpected status %d on iteration %d", rec.Code, i)
		}
	}
	avgMs := float64(totalNs) / float64(iterations) / 1e6
	t.Logf("HandleFragment avg latency over %d iterations: %.3f ms", iterations, avgMs)
}
