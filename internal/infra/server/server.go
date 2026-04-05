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
	"github.com/jadecobra/agbalumo/internal/infra/env"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/module/feedback"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// Setup initializes the Echo server and its dependencies.
// It returns the Echo instance, a cleanup function to close resources, and any error encountered.
func Setup(cfg *config.Config) (*echo.Echo, func(), error) {
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
		return nil, nil, err
	}
	repo.SetSlowQueryThreshold(time.Duration(cfg.SlowQueryThresholdMs) * time.Millisecond)

	listingSvc := listing.NewListingService(repo, repo, repo)
	csvSvc := service.NewCSVService()
	geocodingSvc := service.NewGoogleGeocodingService(cfg.GoogleMapsAPIKey)
	csvSvc.Geocoding = geocodingSvc
	imageSvc := service.NewLocalImageService(cfg.UploadDir)

	app := &env.AppEnv{
		DB:           repo,
		Cfg:          cfg,
		Logger:       slog.Default(),
		CSVService:   csvSvc,
		GeocodingSvc: geocodingSvc,
		ImageSvc:     imageSvc,
		ListingSvc:   listingSvc,
		CatCache:     &env.CategoryCache{},
	}

	renderer, err := ui.NewTemplateRenderer(
		"ui/templates/*.html",
		"ui/templates/partials/*.html",
		"ui/templates/components/*.html",
		"ui/templates/listings/*.html",
		"ui/templates/about.html",
	)
	if err != nil {
		return nil, nil, err
	}
	e.Renderer = renderer

	setupRoutes(e, app)

	bgCtx, cancelBg := context.WithCancel(context.Background())
	setupBackgroundServices(bgCtx, cfg, repo)

	cleanup := func() {
		slog.Info("Executing server cleanup...")
		cancelBg()
		if err := repo.Close(); err != nil {
			slog.Error("Failed to close repository", "error", err)
		}
	}

	return e, cleanup, nil
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

func setupRoutes(e *echo.Echo, app *env.AppEnv) {
	repo := app.DB

	// Modules now use AppEnv for all dependencies.
	listingHandler := listing.NewListingHandler(app)
	adminHandler := admin.NewAdminHandler(app)
	authHandler := auth.NewAuthHandler(app)

	authMw := auth.NewAuthMiddleware(domain.UserStore(repo))
	fbHandler := feedback.NewFeedbackHandler(app)
	pageHandler := common.NewPageHandler(app)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	staticCacheMiddleware := StaticCacheHeaders()
	e.Group("/static", staticCacheMiddleware).Static("/", "ui/static")
	e.Group("/static/uploads", staticCacheMiddleware).Static("/", app.Cfg.UploadDir)

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

func setupBackgroundServices(ctx context.Context, cfg *config.Config, repo *sqlite.SQLiteRepository) {
	if err := seeder.EnsureCategoriesSeeded(ctx, repo, "config/categories.json"); err != nil {
		slog.Error("Failed to seed categories", "error", err)
	}

	if cfg.Env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	}

	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)
}
