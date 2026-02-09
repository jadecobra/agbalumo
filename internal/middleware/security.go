package middleware

import (
	"github.com/labstack/echo/v4"
)

func SecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Content Security Policy
		// Allow:
		// - 'self'
		// - Google Fonts (fonts.googleapis.com, fonts.gstatic.com)
		// - Tailwind CDN (cdn.tailwindcss.com)
		// - HTMX (unpkg.com)
		// - Google Auth (accounts.google.com)
		// - Google Maps (maps.googleapis.com, maps.gstatic.com)
		// - Inline scripts (unsafe-inline) - Required for current setup (HTMX/Tailwind config in HTML)
		//   TODO: Move inline scripts to files to enable stricter CSP.
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com https://unpkg.com https://maps.googleapis.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"font-src 'self' https://fonts.gstatic.com https://fonts.googleapis.com; " +
			"img-src 'self' data: https://*.googleusercontent.com https://ui-avatars.com https://maps.googleapis.com https://maps.gstatic.com; " +
			"connect-src 'self' https://accounts.google.com https://maps.googleapis.com;"

		c.Response().Header().Set("Content-Security-Policy", csp)
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return next(c)
	}
}
