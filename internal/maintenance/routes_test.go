package maintenance

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateRouteMap(t *testing.T) {
	// server.Setup expects to find ui/templates relative to CWD.
	// Since tests run in the package directory, we move up to the project root.
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	// Assuming we are in internal/maintenance, we move up two levels
	if errChdir := os.Chdir("../../"); errChdir != nil {
		t.Fatalf("failed to change directory to root: %v", errChdir)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	routeMap, err := GenerateRouteMap()
	if err != nil {
		t.Fatalf("GenerateRouteMap() failed: %v", err)
	}

	if routeMap == "" {
		t.Error("GenerateRouteMap() returned empty string")
	}

	expectedRoutes := []string{
		"GET", "/healthz",
		"GET", "/about",
	}

	for _, expected := range expectedRoutes {
		if !strings.Contains(routeMap, expected) {
			t.Errorf("Expected substring %q not found in output.\nOutput:\n%s", expected, routeMap)
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
