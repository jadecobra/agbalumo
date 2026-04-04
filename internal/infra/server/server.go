package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/module/feedback"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/repository/cached"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// Setup initializes the Echo server and its dependencies.
func Setup(cfg *config.Config) (*echo.Echo, error) {
	// Initialize Structured Logging
	var logger *slog.Logger
	if cfg.Env == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"status", v.Status,
				"URI", v.URI,
			)
			return nil
		},
	}))

	setupMiddleware(e, cfg)

	repo, err := sqlite.NewSQLiteRepository(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	repo.SetSlowQueryThreshold(time.Duration(cfg.SlowQueryThresholdMs) * time.Millisecond)

	renderer, err := ui.NewTemplateRenderer(
		"ui/templates/*.html",
		"ui/templates/partials/*.html",
		"ui/templates/components/*.html",
		"ui/templates/listings/*.html",
		"ui/templates/about.html",
	)
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	setupRoutes(e, repo, cfg)
	setupBackgroundServices(cfg, repo)

	return e, nil
}

func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.Gzip())
	e.Use(customMiddleware.SecureHeaders)

	if cfg.Env != "test" {
		rateLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
			Rate:  rate.Limit(cfg.RateLimitRate),
			Burst: cfg.RateLimitBurst,
		})
		e.Use(rateLimiter.Middleware())
	}

	baseURL := os.Getenv("BASE_URL")
	isSecure := cfg.Env == "production" || strings.HasPrefix(baseURL, "https://")

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{ //nolint:gosec // false positive for CSRF config
		TokenLookup:    "header:X-CSRF-Token,form:_csrf",
		CookiePath:     "/",
		CookieName:     "_csrf",
		CookieSameSite: http.SameSiteStrictMode,
		CookieSecure:   isSecure,
		CookieHTTPOnly: false,
	}))

	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(customMiddleware.SessionMiddleware(store))
}

func setupRoutes(e *echo.Echo, repo *sqlite.SQLiteRepository, cfg *config.Config) {
	cachedRepo := cached.NewCachedListingStore(repo, 60*time.Second)
	geocodingSvc := service.NewGoogleGeocodingService(cfg.GoogleMapsAPIKey)

	listingSvc := listing.NewListingService(
		domain.ListingStore(cachedRepo),
		domain.CategoryStore(cachedRepo),
		domain.ClaimRequestStore(cachedRepo),
	)

	listingHandler := listing.NewListingHandler(listing.ListingDependencies{
		ListingStore:     domain.ListingStore(cachedRepo),
		CategoryStore:    domain.CategoryStore(cachedRepo),
		ListingSvc:       listingSvc,
		GeocodingSvc:     geocodingSvc,
		Config:           cfg,
		GoogleMapsAPIKey: cfg.GoogleMapsAPIKey,
	})

	csvService := service.NewCSVService()
	csvService.Geocoding = geocodingSvc

	adminHandler := admin.NewAdminHandler(admin.AdminDependencies{
		AdminStore:        domain.AdminStore(repo),
		FeedbackStore:     domain.FeedbackStore(repo),
		AnalyticsStore:    domain.AnalyticsStore(repo),
		CategoryStore:     domain.CategoryStore(repo),
		UserStore:         domain.UserStore(repo),
		ListingStore:      domain.ListingStore(repo),
		ClaimRequestStore: domain.ClaimRequestStore(repo),
		CSVService:        csvService,
		Cfg:               cfg,
	})

	var googleProvider auth.GoogleProvider
	if cfg.MockAuth {
		googleProvider = &auth.MockGoogleProvider{
			Email: "test@agbalumo.com",
			Name:  "Test User",
		}
	}

	authHandler := auth.NewAuthHandler(auth.AuthDependencies{
		UserStore:      domain.UserStore(repo),
		Config:         cfg,
		GoogleProvider: googleProvider,
	})

	authMw := auth.NewAuthMiddleware(domain.UserStore(repo))
	fbHandler := feedback.NewFeedbackHandler(repo)
	pageHandler := common.NewPageHandler(domain.CategoryStore(cachedRepo), cfg)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	staticCacheMiddleware := StaticCacheHeaders()
	e.Group("/static", staticCacheMiddleware).Static("/", "ui/static")
	e.Group("/static/uploads", staticCacheMiddleware).Static("/", cfg.UploadDir)

	e.Use(authMw.OptionalAuth)

	modules := []domain.Registrar{
		authHandler,
		listingHandler,
		adminHandler,
		fbHandler,
	}
	for _, module := range modules {
		module.RegisterRoutes(e, authMw)
	}

	e.GET("/about", pageHandler.HandleAbout)
}

func StaticCacheHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") ||
				strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") ||
				strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".svg") ||
				strings.HasSuffix(path, ".woff2") || strings.HasSuffix(path, ".woff") {
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			return next(c)
		}
	}
}

func setupBackgroundServices(cfg *config.Config, repo *sqlite.SQLiteRepository) {
	ctx := context.Background()
	if err := seeder.EnsureCategoriesSeeded(ctx, repo, "config/categories.json"); err != nil {
		slog.Error("Failed to seed categories", "error", err)
	}

	if cfg.Env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	}

	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)
}
