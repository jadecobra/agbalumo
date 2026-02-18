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
	// Environment Configuration is assumed to be loaded (from .env or env vars)
	env := os.Getenv("AGBALUMO_ENV")

	// Initialize Echo instance
	e := echo.New()

	// Custom CSP Middleware
	// Security Middleware (CSP, Strict-Transport-Security, etc.)
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
		// Fallback for dev, or panic in prod
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
		Secure:   env == "production", // Secure only in prod (or if using TLS in dev)
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(customMiddleware.SessionMiddleware(store))

	// Database Initialization
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "agbalumo.db"
	}
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		return nil, err
	}

	// Template Renderer
	// Include all necessary directories
	renderer, err := ui.NewTemplateRenderer(
		"ui/templates/*.html",
		"ui/templates/partials/*.html",
		"ui/templates/listings/*.html", // Added listings directory
		"ui/templates/about.html",      // Explicitly add about page
	)
	if err != nil {
		return nil, err
	}
	e.Renderer = renderer

	// Handlers
	// Use default (nil) ImageService which defaults to LocalImageService
	listingHandler := handler.NewListingHandler(repo, nil)
	csvService := service.NewCSVService()
	adminHandler := handler.NewAdminHandler(repo, csvService)
	authHandler := handler.NewAuthHandler(repo, nil) // New Auth Handler with default provider

	// Auth Routes
	e.GET("/auth/dev", authHandler.DevLogin)
	e.GET("/auth/logout", authHandler.Logout)
	e.GET("/auth/google/login", authHandler.GoogleLogin)
	e.GET("/auth/google/callback", authHandler.GoogleCallback)

	// Routes
	e.Use(authHandler.OptionalAuth) // Inject user if logged in

	// Static files (CSS, JS, Images)
	e.Static("/static", "ui/static")

	e.GET("/", listingHandler.HandleHome)
	e.GET("/about", listingHandler.HandleAbout)
	e.GET("/listings/fragment", listingHandler.HandleFragment)
	e.GET("/listings/:id", listingHandler.HandleDetail)
	e.POST("/listings", listingHandler.HandleCreate, authHandler.RequireAuth)

	// Edit Routes
	// e.Use(authHandler.OptionalAuth) // Already applied globally above
	e.GET("/listings/:id/edit", listingHandler.HandleEdit, authHandler.RequireAuth)
	e.PUT("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth)
	e.POST("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth) // Fallback support for POST
	e.DELETE("/listings/:id", listingHandler.HandleDelete, authHandler.RequireAuth)

	// Profile
	e.GET("/profile", listingHandler.HandleProfile, authHandler.RequireAuth)

	// Feedback
	feedbackHandler := handler.NewFeedbackHandler(repo)
	e.GET("/feedback/modal", feedbackHandler.HandleModal, authHandler.RequireAuth)
	e.POST("/feedback", feedbackHandler.HandleSubmit, authHandler.RequireAuth)

	// Claim Route
	e.POST("/listings/:id/claim", listingHandler.HandleClaim, authHandler.RequireAuth)

	// Admin Routes
	adminGroup := e.Group("/admin")
	// Use the OptionalAuth middleware first to populate the user context, then AdminMiddleware
	adminGroup.Use(authHandler.OptionalAuth)

	// Admin Claim/Login Routes (Protected by Auth, but not Admin Role)
	adminGroup.GET("/login", adminHandler.HandleLoginView, authHandler.RequireAuth)
	adminGroup.POST("/login", adminHandler.HandleLoginAction, authHandler.RequireAuth)

	adminGroup.Use(adminHandler.AdminMiddleware)

	adminGroup.GET("", adminHandler.HandleDashboard)
	adminGroup.GET("/users", adminHandler.HandleUsers)
	adminGroup.POST("/listings/:id/approve", adminHandler.HandleApprove)
	adminGroup.POST("/listings/:id/reject", adminHandler.HandleReject)
	adminGroup.POST("/upload", adminHandler.HandleBulkUpload)

	// Seed Data (if empty) AND not in production
	ctx := context.Background()
	if env != "production" {
		seeder.EnsureSeeded(ctx, repo)
	} else {
		log.Println("Production environment detected. Skipping automatic data seeding.")
	}

	// Background Services
	bgService := service.NewBackgroundService(repo)
	go bgService.StartTicker(ctx)

	return e, nil
}
