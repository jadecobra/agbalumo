package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	Rate  rate.Limit
	Burst int
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	config   RateLimitConfig
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		config:   config,
	}

	// Background cleanup
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.purge()
	}
}

func (rl *RateLimiter) purge() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.config.Rate, rl.config.Burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getVisitor(ip)

			if !limiter.Allow() {
				return c.String(http.StatusTooManyRequests, "Too Many Requests")
			}

			return next(c)
		}
	}
}
