package maintenance

import (
	"strings"
	"testing"
)

func TestGenerateRouteMap(t *testing.T) {
	routeMap, err := GenerateRouteMap()
	if err != nil {
		t.Fatalf("GenerateRouteMap() failed: %v", err)
	}

	if routeMap == "" {
		t.Error("GenerateRouteMap() returned empty string")
	}

	expectedRoutes := []string{
		"GET  /healthz",
		"GET  /about",
	}

	for _, expected := range expectedRoutes {
		if !strings.Contains(routeMap, expected) {
			t.Errorf("Expected route %q not found in output", expected)
		}
	}

	if !strings.Contains(routeMap, "HandlerName") && !strings.Contains(routeMap, "func") {
		// Echo handler names might vary but should at least contain something descriptive or function references
		// Actually, let's just check for the arrow separator
		if !strings.Contains(routeMap, " → ") {
			t.Error("Route map should contain arrow separator ' → '")
		}
	}
}
