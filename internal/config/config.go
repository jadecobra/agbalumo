package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jadecobra/agbalumo/internal/domain"
)

const EnvBaseURL = "BASE_URL"

type Config struct {
	Env                  string
	DatabaseURL          string
	SessionSecret        string
	AdminCode            string
	DevAuthEmail         string
	UploadDir            string
	GoogleMapsAPIKey     string
	RateLimitRate        int
	RateLimitBurst       int
	SlowQueryThresholdMs int
	HasGoogleAuth        bool
	MockAuth             bool
}

func LoadConfig() *Config {
	env := getEnv(domain.EnvKeyAppEnv, domain.EnvDevelopment)

	uploadDir := getEnv(domain.EnvKeyUploadDir, domain.DefaultUploadDir)
	if !filepath.IsAbs(uploadDir) {
		if cwd, err := os.Getwd(); err == nil {
			uploadDir = filepath.Join(cwd, uploadDir)
		}
	}

	clientID := os.Getenv(domain.EnvKeyGoogleClientID)
	clientSecret := os.Getenv(domain.EnvKeyGoogleClientSecret)
	hasGoogleAuth := clientID != "" && clientSecret != ""
	MockAuth := os.Getenv(domain.EnvKeyMockAuth) == "true"

	return &Config{
		Env:                  env,
		DatabaseURL:          getEnv(domain.EnvKeyDatabaseURL, domain.DefaultDatabaseURL),
		SessionSecret:        getEnv(domain.EnvKeySessionSecret, "dev-secret-key"),
		AdminCode:            getAdminCode(env),
		DevAuthEmail:         getEnv(domain.EnvKeyDevAuthEmail, "dev@agbalumo.com"),
		RateLimitRate:        getEnvAsInt(domain.EnvKeyRateLimitRate, 20),
		RateLimitBurst:       getEnvAsInt(domain.EnvKeyRateLimitBurst, 40),
		UploadDir:            uploadDir,
		GoogleMapsAPIKey:     getEnv(domain.EnvKeyGoogleMapsAPIKey, ""),
		HasGoogleAuth:        hasGoogleAuth || MockAuth,
		MockAuth:             MockAuth,
		SlowQueryThresholdMs: getEnvAsInt(domain.EnvKeySlowQueryThreshold, 50),
	}
}

func getAdminCode(env string) string {
	code := os.Getenv(domain.EnvKeyAdminCode)

	if env == domain.EnvProduction && code == "" {
		slog.Error(domain.EnvKeyAdminCode + " environment variable is required in production")
		os.Exit(1)
	}

	if code == "" {
		slog.Warn("Using default admin code - set " + domain.EnvKeyAdminCode + " for production")
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
