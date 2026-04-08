package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		csp := "default-src 'self'; " +
			"script-src 'self' https://maps.googleapis.com; " +
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
		c.Response().Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		return next(c)
	}
}
// CanonicalPath ensures the request path has a leading slash, preventing
// potential authorization bypasses via non-canonical paths (reference GHSA-p77j-4mvh-x3m3).
func CanonicalPath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		// HTTP/2 :path must be non-empty and start with /.
		// Rejecting on the outer layer mimics the gRPC fix (codes.Unimplemented).
		if path == "" || path[0] != '/' {
			return echo.NewHTTPError(http.StatusNotImplemented, "malformed path: missing leading slash")
		}
		return next(c)
	}
}
