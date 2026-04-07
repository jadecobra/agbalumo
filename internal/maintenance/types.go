package maintenance

import (
	"regexp"
	"strings"
)

// Route represents an HTTP route with method and path.
type Route struct {
	Method string
	Path   string
}

// NormalizePath ensures consistent path formatting for comparison.
func NormalizePath(path string) string {
	p := strings.TrimSpace(path)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = strings.TrimSuffix(p, "/")

	// Treat {id} (OpenAPI) and :id (Echo) as equivalent for comparison
	re := regexp.MustCompile(`\{([^}]+)\}`)
	p = re.ReplaceAllString(p, ":$1")

	if p == "" {
		return "/"
	}
	return p
}

// NewRoute creates a normalized Route.
func NewRoute(method, path string) Route {
	return Route{
		Method: method,
		Path:   NormalizePath(path),
	}
}
