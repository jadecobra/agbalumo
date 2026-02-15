package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter_Allow(t *testing.T) {
	config := RateLimitConfig{
		Rate:  rate.Limit(10), // 10 req/s
		Burst: 2,
	}
	rl := NewRateLimiter(config)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := rl.Middleware()(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// First 2 request should allow (Burst)
	assert.NoError(t, h(c))
	assert.NoError(t, h(c))

	// Third request might fail if instant
	// But rate/limit allows some "time" refill?
	// With burst 2, 3rd instant request should fail.
}

func TestRateLimiter_Cleanup(t *testing.T) {
	config := RateLimitConfig{
		Rate:  rate.Limit(1),
		Burst: 1,
	}
	rl := NewRateLimiter(config)

	// Assuming NewRateLimiter now takes rate, burst, and cleanupInterval
	// Manually inject visitors
	rl.mu.Lock()
	rl.visitors["old-visitor"] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now().Add(-5 * time.Minute), // Should be cleaned
	}
	rl.visitors["new-visitor"] = &visitor{
		limiter:  rate.NewLimiter(rate.Limit(1), 1),
		lastSeen: time.Now(), // Should remain
	}
	rl.mu.Unlock()

	// Call purge directly
	rl.purge()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.visitors["old-visitor"]; exists {
		t.Error("Expected old-visitor to be cleaned up")
	}
	if _, exists := rl.visitors["new-visitor"]; !exists {
		t.Error("Expected new-visitor to remain")
	}
}

func TestRateLimiter_Concurrent(t *testing.T) {
	config := RateLimitConfig{Rate: 100, Burst: 100}
	rl := NewRateLimiter(config)
	e := echo.New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := rl.Middleware()(func(c echo.Context) error { return nil })
			_ = h(c)
		}()
	}
	wg.Wait()
}

// Benchmarks

func BenchmarkRateLimiter_Allowed(b *testing.B) {
	config := RateLimitConfig{Rate: rate.Limit(b.N + 100), Burst: b.N + 100}
	rl := NewRateLimiter(config)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := rl.Middleware()(func(c echo.Context) error { return nil })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(c)
	}
}

func BenchmarkRateLimiter_Concurrent(b *testing.B) {
	config := RateLimitConfig{Rate: rate.Limit(b.N + 100), Burst: b.N + 100}
	rl := NewRateLimiter(config)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := rl.Middleware()(func(c echo.Context) error { return nil })

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h(c)
		}
	})
}
