package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
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
	// Security Headers (CSP, Strict-Transport-Security, etc.)
	e.Use(customMiddleware.SecureHeaders)

	// Rate Limiter
	rateLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  rate.Limit(cfg.RateLimitRate),
		Burst: cfg.RateLimitBurst,
	})
	e.Use(rateLimiter.Middleware())

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
		"ui/templates/listings/*.html",
		"ui/templates/about.html",
	)
}

// setupRoutes registers all HTTP routes.
func setupRoutes(e *echo.Echo, repo *sqlite.SQLiteRepository, cfg *config.Config) {
	// Handlers
	listingHandler := handler.NewListingHandler(repo, nil)
	csvService := service.NewCSVService()
	adminHandler := handler.NewAdminHandler(repo, csvService, cfg)
	authHandler := handler.NewAuthHandler(repo, nil, cfg)

	// Auth Routes
	e.GET("/auth/dev", authHandler.DevLogin)
	e.GET("/auth/logout", authHandler.Logout)
	e.GET("/auth/google/login", authHandler.GoogleLogin)
	e.GET("/auth/google/callback", authHandler.GoogleCallback)

	// Global Auth Middleware
	e.Use(authHandler.OptionalAuth)

	// Static files
	e.Static("/static", "ui/static")

	// Public Routes
	e.GET("/", listingHandler.HandleHome)
	e.GET("/about", listingHandler.HandleAbout)
	e.GET("/listings/fragment", listingHandler.HandleFragment)
	e.GET("/listings/:id", listingHandler.HandleDetail)

	// Authenticated Routes
	e.POST("/listings", listingHandler.HandleCreate, authHandler.RequireAuth)
	e.GET("/listings/:id/edit", listingHandler.HandleEdit, authHandler.RequireAuth)
	e.PUT("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth)
	e.POST("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth)
	e.DELETE("/listings/:id", listingHandler.HandleDelete, authHandler.RequireAuth)
	e.GET("/profile", listingHandler.HandleProfile, authHandler.RequireAuth)
	e.POST("/listings/:id/claim", listingHandler.HandleClaim, authHandler.RequireAuth)

	// Feedback
	feedbackHandler := handler.NewFeedbackHandler(repo)
	e.GET("/feedback/modal", feedbackHandler.HandleModal, authHandler.RequireAuth)
	e.POST("/feedback", feedbackHandler.HandleSubmit, authHandler.RequireAuth)

	// Admin Routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(authHandler.OptionalAuth)
	adminGroup.GET("/login", adminHandler.HandleLoginView, authHandler.RequireAuth)
	adminGroup.POST("/login", adminHandler.HandleLoginAction, authHandler.RequireAuth)
	adminGroup.Use(adminHandler.AdminMiddleware)
	adminGroup.GET("", adminHandler.HandleDashboard)
	adminGroup.GET("/users", adminHandler.HandleUsers)
	adminGroup.POST("/listings/:id/approve", adminHandler.HandleApprove)
	adminGroup.POST("/listings/:id/reject", adminHandler.HandleReject)
	adminGroup.POST("/upload", adminHandler.HandleBulkUpload)
}

// setupBackgroundServices starts seeding and background tickers.
func setupBackgroundServices(cfg *config.Config, repo *sqlite.SQLiteRepository) {
	ctx := context.Background()

	if cfg.Env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	} else {
		slog.Info("Production environment detected. Skipping automatic data seeding.")
	}

	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)
}
