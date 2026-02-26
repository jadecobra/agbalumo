# Middleware

## OVERVIEW
HTTP middleware chain: auth, sessions, rate limiting, security headers.

## WHERE TO LOOK

| Task | File |
|------|------|
| Protect routes | `auth.go` - `RequireAuth()`, `OptionalAuth()` |
| Session handling | `session.go` - `SessionMiddleware()`, `GetSession(c)` |
| Request throttling | `ratelimit.go` - `NewRateLimiter()`, `rl.Middleware()` |
| Security headers | `security.go` - `SecureHeaders()` |

## CONVENTIONS

- Middleware returns `echo.MiddlewareFunc`: `func(next echo.HandlerFunc) echo.HandlerFunc`
- Chain via `e.Use(middlewareFn)` in router setup
- Rate limiter: IP-based, in-memory, with background purge goroutine
- Sessions: gorilla/sessions, store configured externally, key `"auth_session"`
- Auth: user injected into context via `c.Set("User", user)`, retrieved via `c.Get("User")`

## ANTI-PATTERNS

- Do NOT hardcode secrets in `security.go` CSP - use environment variables
- Do NOT block on rate limiter cleanup - runs in background goroutine
- Do NOT skip error handling in `RequireAuth` redirects
- Do NOT use rate limit for login endpoints without separate configuration
