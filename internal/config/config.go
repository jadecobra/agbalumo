package config

import (
	"log/slog"
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	Env            string
	DatabaseURL    string
	SessionSecret  string
	AdminCode      string
	DevAuthEmail   string
	RateLimitRate  int
	RateLimitBurst int
}

// LoadConfig reads environment variables and returns a populated Config with defaults.
func LoadConfig() *Config {
	env := getEnv("AGBALUMO_ENV", "development")

	return &Config{
		Env:            env,
		DatabaseURL:    getEnv("DATABASE_URL", "agbalumo.db"),
		SessionSecret:  getEnv("SESSION_SECRET", "dev-secret-key"),
		AdminCode:      getAdminCode(env),
		DevAuthEmail:   getEnv("DEV_AUTH_EMAIL", "dev@agbalumo.com"),
		RateLimitRate:  getEnvAsInt("RATE_LIMIT_RATE", 20),
		RateLimitBurst: getEnvAsInt("RATE_LIMIT_BURST", 40),
	}
}

// getAdminCode returns the admin code, failing in production if not set.
func getAdminCode(env string) string {
	code := os.Getenv("ADMIN_CODE")

	if env == "production" && code == "" {
		slog.Error("ADMIN_CODE environment variable is required in production")
		os.Exit(1)
	}

	// In development, use a default but warn
	if code == "" {
		slog.Warn("Using default admin code - set ADMIN_CODE for production")
		code = "agbalumo2024"
	}

	return code
}

// getEnv returns the env var value or a fallback if empty.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// getEnvAsInt returns the env var as an integer, or the fallback if missing/invalid.
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
