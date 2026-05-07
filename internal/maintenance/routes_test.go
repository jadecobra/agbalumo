package maintenance

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateRouteMap(t *testing.T) {
	cleanup := setupTestRoot(t)
	defer cleanup()

	routeMap, err := GenerateRouteMap()
	if err != nil {
		t.Fatalf("GenerateRouteMap() failed: %v", err)
	}

	if routeMap == "" {
		t.Error("GenerateRouteMap() returned empty string")
	}

	verifyExpectedRoutes(t, routeMap)

	if !strings.Contains(routeMap, " → ") {
		t.Error("Route map should contain arrow separator ' → '")
	}
}

func setupTestRoot(t *testing.T) func() {
	t.Helper()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir("../../"); err != nil {
		t.Fatalf("failed to change directory to root: %v", err)
	}
	return func() {
		_ = os.Chdir(origWd)
	}
}

func verifyExpectedRoutes(t *testing.T, routeMap string) {
	t.Helper()
	expectedRoutes := []string{
		"GET", "/healthz",
		"GET", "/about",
	}

	for _, expected := range expectedRoutes {
		if !strings.Contains(routeMap, expected) {
			t.Errorf("Expected substring %q not found in output.\nOutput:\n%s", expected, routeMap)
		}
	}
}
