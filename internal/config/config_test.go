package config_test

import (
	"os"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Clean env before testing
	keys := []string{"AGBALUMO_ENV", "DATABASE_URL", "SESSION_SECRET", "ADMIN_CODE", "DEV_AUTH_EMAIL", "RATE_LIMIT_RATE", "RATE_LIMIT_BURST"}
	for _, k := range keys {
		os.Unsetenv(k)
	}

	t.Run("defaults", func(t *testing.T) {
		cfg := config.LoadConfig()
		require.Equal(t, "development", cfg.Env)
		require.Equal(t, "agbalumo.db", cfg.DatabaseURL)
		require.Equal(t, "dev-secret-key", cfg.SessionSecret)
		require.Equal(t, "agbalumo2024", cfg.AdminCode)
		require.Equal(t, "dev@agbalumo.com", cfg.DevAuthEmail)
		require.Equal(t, 20, cfg.RateLimitRate)
		require.Equal(t, 40, cfg.RateLimitBurst)
	})

	t.Run("overrides", func(t *testing.T) {
		os.Setenv("AGBALUMO_ENV", "production")
		os.Setenv("DATABASE_URL", "prod.db")
		os.Setenv("SESSION_SECRET", "super-secret")
		os.Setenv("ADMIN_CODE", "admin123")
		os.Setenv("DEV_AUTH_EMAIL", "test@example.com")
		os.Setenv("RATE_LIMIT_RATE", "50")
		os.Setenv("RATE_LIMIT_BURST", "100")
		defer func() {
			for _, k := range keys {
				os.Unsetenv(k)
			}
		}()

		cfg := config.LoadConfig()
		require.Equal(t, "production", cfg.Env)
		require.Equal(t, "prod.db", cfg.DatabaseURL)
		require.Equal(t, "super-secret", cfg.SessionSecret)
		require.Equal(t, "admin123", cfg.AdminCode)
		require.Equal(t, "test@example.com", cfg.DevAuthEmail)
		require.Equal(t, 50, cfg.RateLimitRate)
		require.Equal(t, 100, cfg.RateLimitBurst)
	})

	t.Run("invalid int fallback", func(t *testing.T) {
		os.Setenv("RATE_LIMIT_RATE", "abc")
		defer os.Unsetenv("RATE_LIMIT_RATE")

		cfg := config.LoadConfig()
		require.Equal(t, 20, cfg.RateLimitRate) // Default fallback
	})
}
