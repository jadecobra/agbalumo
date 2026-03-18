package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/repository/cached"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// SetupServer initializes the Echo server and its dependencies.
// It returns the Echo instance or an error.
func SetupServer() (*echo.Echo, error) {
	cfg := config.LoadConfig()

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

	repo, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	renderer, err := initRenderer()
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	setupRoutes(e, repo, cfg)
	setupBackgroundServices(cfg, repo)

	return e, nil
}

// setupMiddleware wires security, rate limiting, CSRF, and session middleware.
func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	// P0.1: Gzip Compression — reduces payload 60-80%
	e.Use(middleware.Gzip())

	// Security Headers (CSP, Strict-Transport-Security, etc.)
	e.Use(customMiddleware.SecureHeaders)

	// Rate Limiter (skip in test environment)
	if cfg.Env != "test" {
		rateLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
			Rate:  rate.Limit(cfg.RateLimitRate),
			Burst: cfg.RateLimitBurst,
		})
		e.Use(rateLimiter.Middleware())
	}

	// CSRF Protection
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token,form:_csrf",
		CookiePath:     "/",
		CookieName:     "_csrf",
		CookieSameSite: http.SameSiteStrictMode,
		CookieSecure:   cfg.Env == "production",
		CookieHTTPOnly: false,
	}))

	// Session Middleware
	if cfg.SessionSecret == "dev-secret-key" && cfg.Env == "production" {
		slog.Error("SESSION_SECRET must be set in production")
		os.Exit(1)
	} else if cfg.SessionSecret == "dev-secret-key" {
		slog.Warn("Using default dev session key")
	}

	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   cfg.Env == "production",
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(customMiddleware.SessionMiddleware(store))
}

// initDatabase creates the SQLite repository.
func initDatabase(dbPath string) (*sqlite.SQLiteRepository, error) {
	return sqlite.NewSQLiteRepository(dbPath)
}

// initRenderer creates the template renderer.
func initRenderer() (*ui.TemplateRenderer, error) {
	return ui.NewTemplateRenderer(
		"ui/templates/*.html",
		"ui/templates/partials/*.html",
		"ui/templates/components/*.html",
		"ui/templates/listings/*.html",
		"ui/templates/about.html",
	)
}

// setupRoutes registers all HTTP routes.
func setupRoutes(e *echo.Echo, repo *sqlite.SQLiteRepository, cfg *config.Config) {
	// Handlers
	// P2.3: Wrap repo with cached store for GetCounts (60s TTL)
	cachedRepo := cached.NewCachedListingStore(repo, 60*time.Second)
	geocodingSvc := service.NewGoogleGeocodingService(cfg.GoogleMapsAPIKey)
	listingSvc := service.NewListingService(
		domain.ListingStore(cachedRepo),
		domain.CategoryStore(cachedRepo),
		domain.ClaimRequestStore(cachedRepo),
	)
	listingHandler := handler.NewListingHandler(
		domain.ListingStore(cachedRepo),
		domain.CategoryStore(cachedRepo),
		listingSvc,
		nil,
		geocodingSvc,
		cfg,
	)
	listingHandler.GoogleMapsAPIKey = cfg.GoogleMapsAPIKey
	csvService := service.NewCSVService()
	csvService.Geocoding = geocodingSvc
	adminHandler := handler.NewAdminHandler(
		domain.AdminStore(repo),
		domain.FeedbackStore(repo),
		domain.AnalyticsStore(repo),
		domain.CategoryStore(repo),
		domain.UserStore(repo),
		domain.ListingStore(repo),
		domain.ClaimRequestStore(repo),
		csvService,
		cfg,
	)
	authHandler := auth.NewAuthHandler(domain.UserStore(repo), nil, cfg)
	authMw := auth.NewAuthMiddleware(domain.UserStore(repo))

	// Auth Routes
	e.GET("/auth/dev", authHandler.DevLogin)
	e.GET("/auth/logout", authHandler.Logout)
	e.GET("/auth/google/login", authHandler.GoogleLogin)
	e.GET("/auth/google/callback", authHandler.GoogleCallback)

	// Health Check (before auth middleware — bypasses CSRF, rate limiting, sessions)
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Global Auth Middleware
	e.Use(authMw.OptionalAuth)

	// Static files with cache-control (P0.2)
	staticCacheMiddleware := staticCacheHeaders()
	e.Group("/static", staticCacheMiddleware).Static("/", "ui/static")

	// Serve uploaded images at /static/uploads
	// This is needed because in production, uploads go to a different directory (e.g., /data/uploads)
	e.Group("/static/uploads", staticCacheMiddleware).Static("/", cfg.UploadDir)

	// Public Routes
	e.GET("/", listingHandler.HandleHome)
	e.GET("/about", listingHandler.HandleAbout)
	e.GET("/listings/fragment", listingHandler.HandleFragment)
	e.GET("/listings/:id", listingHandler.HandleDetail)

	// Authenticated Routes
	e.POST("/listings", listingHandler.HandleCreate, authMw.RequireAuth)
	e.GET("/listings/:id/edit", listingHandler.HandleEdit, authMw.RequireAuth)
	e.PUT("/listings/:id", listingHandler.HandleUpdate, authMw.RequireAuth)
	e.POST("/listings/:id", listingHandler.HandleUpdate, authMw.RequireAuth)
	e.DELETE("/listings/:id", listingHandler.HandleDelete, authMw.RequireAuth)
	e.GET("/profile", listingHandler.HandleProfile, authMw.RequireAuth)
	e.POST("/listings/:id/claim", listingHandler.HandleClaim, authMw.RequireAuth)

	// Feedback
	feedbackHandler := handler.NewFeedbackHandler(repo)
	e.GET("/feedback/modal", feedbackHandler.HandleModal, authMw.RequireAuth)
	e.POST("/feedback", feedbackHandler.HandleSubmit, authMw.RequireAuth)

	// Admin Routes
	// Strict rate limiter for sensitive admin login endpoint (5 req/min, burst 10)
	adminAuthLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  rate.Limit(5),
		Burst: 10,
	})

	adminGroup := e.Group("/admin")
	adminGroup.Use(authMw.OptionalAuth)
	adminGroup.GET("/login", adminHandler.HandleLoginView)

	// Admin login POST with strict rate limiting
	adminLoginGroup := adminGroup.Group("/login")
	adminLoginGroup.Use(adminAuthLimiter.Middleware())
	adminLoginGroup.POST("", adminHandler.HandleLoginAction)
	adminGroup.Use(adminHandler.AdminMiddleware)
	adminGroup.GET("", adminHandler.HandleDashboard)
	adminGroup.GET("/users", adminHandler.HandleUsers)
	adminGroup.GET("/listings", adminHandler.HandleAllListings)
	adminGroup.POST("/claims/:id/approve", adminHandler.HandleApproveClaim)
	adminGroup.POST("/claims/:id/reject", adminHandler.HandleRejectClaim)
	adminGroup.POST("/listings/bulk", adminHandler.HandleBulkAction)
	adminGroup.GET("/listings/delete-confirm", adminHandler.HandleAdminDeleteView)
	adminGroup.POST("/listings/delete", adminHandler.HandleAdminDeleteAction)
	adminGroup.POST("/listings/:id/featured", adminHandler.HandleToggleFeatured)
	adminGroup.POST("/upload", adminHandler.HandleBulkUpload)
	adminGroup.GET("/listings/export", adminHandler.HandleExportListings)
	adminGroup.POST("/categories", adminHandler.HandleAddCategory)
}

// staticCacheHeaders returns middleware that sets Cache-Control headers for static assets.
// Immutable assets (CSS, JS, fonts, images) get a 1-year cache.
func staticCacheHeaders() echo.MiddlewareFunc {
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

// setupBackgroundServices starts seeding and background tickers.
func setupBackgroundServices(cfg *config.Config, repo *sqlite.SQLiteRepository) {
	ctx := context.Background()

	// Always ensure core categories are seeded/upserted from config
	// This lets developers update categories.json and deploy without manual DB intervention.
	if err := seeder.EnsureCategoriesSeeded(ctx, repo, "config/categories.json"); err != nil {
		slog.Error("Failed to seed categories", "error", err)
	}

	if cfg.Env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	} else {
		slog.Info("Production environment detected. Skipping automatic data seeding.")
	}

	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)
}

// dummy change
