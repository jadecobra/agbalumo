package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func setupHandler(rl *RateLimiter) echo.HandlerFunc {
	return rl.Middleware()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
}

func TestRateLimiter(t *testing.T) {
	t.Parallel()
	e := echo.New()

	config := RateLimitConfig{
		Rate:  10,
		Burst: 20,
	}
	rl := NewRateLimiter(config)
	h := setupHandler(rl)

	t.Run("allow requests within limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, h(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("block requests exceeding limit", func(t *testing.T) {
		// Create a strict limiter for this subtest
		strictRl := NewRateLimiter(RateLimitConfig{
			Rate:  1,
			Burst: 1,
		})
		strictH := setupHandler(strictRl)

		// First request allowed
		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		rec1 := httptest.NewRecorder()
		c1 := e.NewContext(req1, rec1)
		_ = strictH(c1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		// Second request immediately after should be blocked
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)
		_ = strictH(c2)
		assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
	})
}

func TestRateLimiter_CustomConfig(t *testing.T) {
	t.Parallel()
	e := echo.New()
	config := RateLimitConfig{
		Rate:  rate.Limit(5),
		Burst: 10,
	}
	rl := NewRateLimiter(config)
	assert.Equal(t, rate.Limit(5), rl.config.Rate)
	assert.Equal(t, 10, rl.config.Burst)

	h := setupHandler(rl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, h(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimiter_Concurrency(t *testing.T) {
	t.Parallel()
	e := echo.New()
	rl := NewRateLimiter(RateLimitConfig{
		Rate:  1000,
		Burst: 1000,
	})
	h := setupHandler(rl)

	const numRequests = 100
	done := make(chan bool)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = h(c)
			done <- true
		}()
	}

	for i := 0; i < numRequests; i++ {
		<-done
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	t.Parallel()
	e := echo.New()
	rl := NewRateLimiter(RateLimitConfig{
		Rate:  10,
		Burst: 20,
	})

	h := setupHandler(rl)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "1.2.3.4"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = h(c)
	}

	// Artificially age the visitor
	rl.mu.Lock()
	if v, ok := rl.visitors["1.2.3.4"]; ok {
		v.lastSeen = time.Now().Add(-4 * time.Minute)
	}
	rl.mu.Unlock()

	// Call purge manually
	rl.purge()

	// Verify visitor was removed
	rl.mu.Lock()
	_, exists := rl.visitors["1.2.3.4"]
	rl.mu.Unlock()
	assert.False(t, exists)
}
