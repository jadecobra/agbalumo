package config

import (
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
	return &Config{
		Env:            getEnv("AGBALUMO_ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", "agbalumo.db"),
		SessionSecret:  getEnv("SESSION_SECRET", "dev-secret-key"),
		AdminCode:      getEnv("ADMIN_CODE", "agbalumo2024"),
		DevAuthEmail:   getEnv("DEV_AUTH_EMAIL", "dev@agbalumo.com"),
		RateLimitRate:  getEnvAsInt("RATE_LIMIT_RATE", 20),
		RateLimitBurst: getEnvAsInt("RATE_LIMIT_BURST", 40),
	}
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
