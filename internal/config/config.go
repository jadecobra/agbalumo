package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Env              string
	DatabaseURL      string
	SessionSecret    string
	AdminCode        string
	DevAuthEmail     string
	RateLimitRate    int
	RateLimitBurst   int
	UploadDir        string
	GoogleMapsAPIKey     string
	HasGoogleAuth        bool
	MockAuth             bool
	SlowQueryThresholdMs int
}

func LoadConfig() *Config {
	env := getEnv("AGBALUMO_ENV", "development")

	uploadDir := getEnv("UPLOAD_DIR", "ui/static/uploads")
	if !filepath.IsAbs(uploadDir) {
		if cwd, err := os.Getwd(); err == nil {
			uploadDir = filepath.Join(cwd, uploadDir)
		}
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	hasGoogleAuth := clientID != "" && clientSecret != ""

	return &Config{
		Env:              env,
		DatabaseURL:      getEnv("DATABASE_URL", ".tester/data/agbalumo.db"),
		SessionSecret:    getEnv("SESSION_SECRET", "dev-secret-key"),
		AdminCode:        getAdminCode(env),
		DevAuthEmail:     getEnv("DEV_AUTH_EMAIL", "dev@agbalumo.com"),
		RateLimitRate:    getEnvAsInt("RATE_LIMIT_RATE", 20),
		RateLimitBurst:   getEnvAsInt("RATE_LIMIT_BURST", 40),
		UploadDir:            uploadDir,
		GoogleMapsAPIKey:     getEnv("GOOGLE_MAPS_API_KEY", ""),
		HasGoogleAuth:        hasGoogleAuth || os.Getenv("MOCK_AUTH") == "true",
		MockAuth:             os.Getenv("MOCK_AUTH") == "true",
		SlowQueryThresholdMs: getEnvAsInt("SLOW_QUERY_THRESHOLD_MS", 50),
	}
}

func getAdminCode(env string) string {
	code := os.Getenv("ADMIN_CODE")

	if env == "production" && code == "" {
		slog.Error("ADMIN_CODE environment variable is required in production")
		os.Exit(1)
	}

	if code == "" {
		slog.Warn("Using default admin code - set ADMIN_CODE for production")
		code = "agbalumo2024"
	}

	return code
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}
