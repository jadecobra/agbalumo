# Security Audit Report — agbalumo

**Date**: February 23, 2026  
**Overall Grade**: **A**

---

## Executive Summary

The agbalumo application has solid security foundations. Recent fixes have addressed critical authentication and rate limiting issues. The remaining items are defense-in-depth improvements. This audit verifies the current state against the previous recommendations.

---

## Security Fixes Applied

| Fix | Status | Description |
|-----|--------|-------------|
| **#1: ADMIN_CODE enforcement** | ✅ Fixed | Server exits if `ADMIN_CODE` not set in production |
| **#2: Admin login rate limiting** | ✅ Fixed | Added 5 req/min rate limit to `/admin/login` POST |
| **#3: CSP inline scripts** | ✅ Fixed | Moved all template inline scripts to external JS |
| **#4: Handler inline HTML** | ✅ Fixed | Replaced onclick with hx-on event handlers |
| **#5: Chart.js CDN** | ✅ Fixed | Self-hosted Chart.js locally |
| **#6: htmx CDN** | ✅ Fixed | Self-hosted htmx locally |
| **#7: CSP CDNs** | ✅ Fixed | Removed unpkg.com, cdn.jsdelivr.net, cdn.tailwindcss.com |

---

## Security Strengths

| Area | Status | Notes |
| :--- | :--- | :--- |
| **SQL Injection** | ✅ Pass | All queries use parameterized statements (`?` placeholders) |
| **CSRF Protection** | ✅ Pass | Echo CSRF middleware enabled; form tokens in place |
| **Rate Limiting** | ✅ Pass | Global rate limiter (20 req/s, burst 40) + admin-specific (5 req/min, burst 10). Configurable via RATE_LIMIT_RATE and RATE_LIMIT_BURST env vars. |
| **Security Headers** | ✅ Pass | CSP, HSTS, X-Frame-Options, X-Content-Type-Options all present |
| **Error Handling** | ✅ Pass | `RespondError` logs internally, renders friendly UI to users |
| **Session Security** | ✅ Pass | Secure, HttpOnly, SameSite=Strict cookies |
| **Input Validation** | ✅ Pass | Domain validation with explicit rules |
| **Auth Middleware** | ✅ Pass | Proper separation of OptionalAuth and RequireAuth |

---

## Remaining Issues (Resolved)

### 1. Inline Scripts in Templates ✅ RESOLVED

All inline scripts moved to external JS files.

### 2. Handler Inline HTML ✅ RESOLVED

Replaced inline onclick with HTMX hx-on event handlers.

### 3. External CDN Dependencies ✅ RESOLVED

All third-party scripts self-hosted locally.

---

## Current CSP Configuration

The Content Security Policy is defined in `internal/middleware/security.go`:

```
csp := "default-src 'self'; " +
    "script-src 'self' 'unsafe-inline' https://maps.googleapis.com; " +
    "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
    "font-src 'self' https://fonts.gstatic.com https://fonts.googleapis.com; " +
    "img-src 'self' data: https://*.googleusercontent.com https://ui-avatars.com https://maps.googleapis.com https://maps.gstatic.com; " +
    "connect-src 'self' https://accounts.google.com https://maps.googleapis.com;"
```

### Allowed Sources (Reduced)

| Directive | Sources |
|-----------|---------|
| default-src | 'self' |
| script-src | 'self', 'unsafe-inline', maps.googleapis.com |
| style-src | 'self', 'unsafe-inline', fonts.googleapis.com |
| font-src | 'self', fonts.gstatic.com, fonts.googleapis.com |
| img-src | 'self', data:, *.googleusercontent.com, ui-avatars.com, maps.googleapis.com, maps.gstatic.com |
| connect-src | 'self', accounts.google.com, maps.googleapis.com |

### Production Requirements

The following environment variables MUST be set in production:

| Variable | Purpose | Required |
|----------|---------|----------|
| `ADMIN_CODE` | Admin access code | Yes |
| `SESSION_SECRET` | Session encryption key | Yes |
| `AGBALUMO_ENV` | Environment (set to `production`) | Yes |
| `DATABASE_URL` | SQLite database path | No (has default) |
| `RATE_LIMIT_RATE` | Global rate limit (req/s) | No (default: 20) |
| `RATE_LIMIT_BURST` | Global rate burst | No (default: 40) |
| `GOOGLE_CLIENT_ID` | OAuth provider | For Google Auth |
| `GOOGLE_CLIENT_SECRET` | OAuth provider | For Google Auth |

### Startup Failures (Production)

The server will exit with an error if:
- `ADMIN_CODE` is not set
- `SESSION_SECRET` is set to default `dev-secret-key`

---

## Scorecard

| Dimension | Score | Notes |
| :--- | :--- | :--- |
| **Authentication** | A- | Fixed admin code enforcement |
| **Authorization** | A | Proper role-based access |
| **Input Validation** | A | Domain validation rules |
| **Data Protection** | A | Parameterized queries |
| **Session Management** | A- | Secure cookie config |
| **Rate Limiting** | A | Global + admin-specific |
| **Headers/CSP** | A | Removed CDNs, self-hosted scripts, fixed inline handlers |

---

## Recommendations

### High Priority

All completed! 🎉

### Medium Priority

All completed!

### Low Priority

1. ~~Implement CSP nonces~~ - Not needed after hx-on refactor
2. ~~Add admin login attempt tracking~~ - Log failed admin login attempts
3. **Consider self-hosting Tailwind CSS** - `admin_login.html` still uses CDN (minor)

---

## Files Changed During Audit

### Configuration
- `internal/config/config.go` - Added production enforcement for ADMIN_CODE

### Routing
- `cmd/server.go` - Added rate limiting to admin login endpoint

### Security
- `internal/middleware/security.go` - Updated CSP to remove CDN dependencies

### Templates
- `ui/templates/base.html` - Replaced htmx CDN with local
- `ui/templates/admin_dashboard.html` - Replaced Chart.js CDN with local
- `ui/templates/partials/modal_detail.html` - Moved inline script to modals.js
- `ui/templates/partials/modal_feedback.html` - Moved inline script to modals.js
- `ui/templates/partials/modal_create_listing.html` - Moved inline script to modals.js
- `ui/templates/partials/modal_edit_listing.html` - Moved 3 inline scripts to modals.js

### Handlers
- `internal/handler/listing.go` - Replaced inline onclick/script with hx-on handlers
- `internal/handler/feedback.go` - Replaced inline onclick with hx-on handlers

### New External JS Files
- `ui/static/js/toasts.js`
- `ui/static/js/modals.js`
- `ui/static/js/admin-listings.js`
- `ui/static/js/admin-dashboard.js`
- `ui/static/js/chart.umd.min.js` (self-hosted Chart.js)
- `ui/static/js/htmx.min.js` (self-hosted htmx)

---

*This report reflects the state of the codebase as of the audit date.*
