package maintenance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Invariants represents the project-wide invariants.
type Invariants struct {
	Protocol            string  `json:"protocol"`
	Port                string  `json:"port"`
	DBEngine            string  `json:"db_engine"`
	CSPPolicy           string  `json:"csp_policy"`
	DefaultCoverage     float64 `json:"default_coverage"`
	MaxFeaturedListings int     `json:"max_featured_listings"`
}

// DumpInvariants generates .agents/invariants.json from project config sources.
func DumpInvariants(rootDir string) error {
	inv := Invariants{
		DBEngine:  "sqlite",
		CSPPolicy: "script-src 'self'",
	}

	// 1. Read .env for BASE_URL
	protocol, port := parseBaseURL(rootDir)
	inv.Protocol = protocol
	inv.Port = port

	// 2. Read .agents/coverage.json for default threshold
	inv.DefaultCoverage = getCoverageThreshold(rootDir)

	// 3. Read internal/module/admin/listings.go for maxFeatured
	inv.MaxFeaturedListings = getMaxFeatured(rootDir)

	// Write to .agents/invariants.json
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal invariants: %w", err)
	}

	outputPath := filepath.Join(rootDir, ".agents", "invariants.json")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0750); err != nil {
		return fmt.Errorf("failed to create agents directory: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write invariants.json: %w", err)
	}

	fmt.Printf("✅ Generated %s\n", outputPath)
	return nil
}

const (
	defaultProtocol = "https"
	defaultPort     = "8443"
)

func parseBaseURL(rootDir string) (string, string) {
	envPath := filepath.Join(rootDir, ".env")
	file, err := os.Open(envPath) //nolint:gosec // maintenance utility
	if err != nil {
		return defaultProtocol, defaultPort
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "BASE_URL=") {
			return extractURLParts(strings.TrimPrefix(line, "BASE_URL="))
		}
	}
	return defaultProtocol, defaultPort
}

func extractURLParts(val string) (string, string) {
	u, err := url.Parse(val)
	if err != nil || u.Scheme == "" {
		return defaultProtocol, defaultPort
	}

	port := u.Port()
	if port != "" {
		return u.Scheme, port
	}

	if u.Scheme == defaultProtocol {
		return u.Scheme, "443"
	}
	return u.Scheme, "80"
}

func getCoverageThreshold(rootDir string) float64 {
	path := filepath.Join(rootDir, ".agents", "coverage.json")
	data, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return 0
	}

	var thresholds map[string]float64
	if err := json.Unmarshal(data, &thresholds); err != nil {
		return 0
	}

	if val, ok := thresholds["default"]; ok {
		return val
	}
	return 0
}

func getMaxFeatured(rootDir string) int {
	path := filepath.Join(rootDir, "internal/module/admin/listings.go")
	content, err := os.ReadFile(path) //nolint:gosec // maintenance utility
	if err != nil {
		return 0
	}

	// Try to find constant first
	reConst := regexp.MustCompile(`maxFeatured\s*=\s*(\d+)`)
	match := reConst.FindSubmatch(content)
	if len(match) > 1 {
		val, _ := strconv.Atoi(string(match[1]))
		return val
	}

	// Fallback to hardcoded logic check
	reHard := regexp.MustCompile(`len\(featured\)\s*>=\s*(\d+)`)
	matchHard := reHard.FindSubmatch(content)
	if len(matchHard) > 1 {
		val, _ := strconv.Atoi(string(matchHard[1]))
		return val
	}

	return 0
}
