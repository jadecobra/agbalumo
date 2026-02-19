package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupServer initializes the Echo server and its dependencies.
// It returns the Echo instance or an error.
func SetupServer() (*echo.Echo, error) {
	env := os.Getenv("AGBALUMO_ENV")

	e := echo.New()

	setupMiddleware(e, env)

	repo, err := initDatabase()
	if err != nil {
		return nil, err
	}

	renderer, err := initRenderer()
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	setupRoutes(e, repo)
	setupBackgroundServices(env, repo)

	return e, nil
}

// setupMiddleware wires security, rate limiting, CSRF, and session middleware.
func setupMiddleware(e *echo.Echo, env string) {
	// Security Headers (CSP, Strict-Transport-Security, etc.)
	e.Use(customMiddleware.SecureHeaders)

	// Rate Limiter
	rateLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  20,
		Burst: 40,
	})
	e.Use(rateLimiter.Middleware())

	// CSRF Protection
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token,form:_csrf",
		CookiePath:     "/",
		CookieName:     "_csrf",
		CookieSameSite: http.SameSiteStrictMode,
		CookieSecure:   env == "production",
		CookieHTTPOnly: false,
	}))

	// Session Middleware
	sessionKey := os.Getenv("SESSION_SECRET")
	if sessionKey == "" {
		if os.Getenv("AGBALUMO_ENV") == "production" {
			log.Fatal("SESSION_SECRET must be set in production")
		}
		sessionKey = "dev-secret-key"
		log.Println("[WARN] Using default dev session key")
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   env == "production",
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(customMiddleware.SessionMiddleware(store))
}

// initDatabase creates the SQLite repository.
func initDatabase() (*sqlite.SQLiteRepository, error) {
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "agbalumo.db"
	}
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
func setupRoutes(e *echo.Echo, repo *sqlite.SQLiteRepository) {
	// Handlers
	listingHandler := handler.NewListingHandler(repo, nil)
	csvService := service.NewCSVService()
	adminHandler := handler.NewAdminHandler(repo, csvService)
	authHandler := handler.NewAuthHandler(repo, nil)

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
func setupBackgroundServices(env string, repo *sqlite.SQLiteRepository) {
	ctx := context.Background()

	if env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	} else {
		log.Println("Production environment detected. Skipping automatic data seeding.")
	}

	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)
}
