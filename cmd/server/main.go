package main

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/jadecobra/agbalumo/internal/seeder"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	// Try loading from local .env or the scripts location
	godotenv.Load(".env")
	if err := godotenv.Load("../scripts/agbalumo/.env"); err != nil {
		log.Printf("Error loading ../scripts/agbalumo/.env: %v", err)
	}

	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.SecureHeaders)

	// Rate Limiter
	rateLimiter := customMiddleware.NewRateLimiter(customMiddleware.RateLimitConfig{
		Rate:  20,
		Burst: 40,
	})
	e.Use(rateLimiter.Middleware())

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
	e.Use(customMiddleware.SessionMiddleware(store))

	// Database Initialization
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "agbalumo.db"
	}
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Template Renderer
	renderer, err := ui.NewTemplateRenderer("ui/templates/*.html", "ui/templates/partials/*.html")
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}
	e.Renderer = renderer

	// Handlers
	listingHandler := handler.NewListingHandler(repo)
	adminHandler := handler.NewAdminHandler(repo)
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

	// Admin Routes
	adminGroup := e.Group("/admin")
	// Use the OptionalAuth middleware first to populate the user context, then AdminMiddleware
	adminGroup.Use(authHandler.OptionalAuth)

	// Admin Claim/Login Routes (Protected by Auth, but not Admin Role)
	adminGroup.GET("/login", adminHandler.HandleLoginView, authHandler.RequireAuth)
	adminGroup.POST("/login", adminHandler.HandleLoginAction, authHandler.RequireAuth)
	
	adminGroup.Use(adminHandler.AdminMiddleware)

	adminGroup.GET("", adminHandler.HandleDashboard)
	adminGroup.POST("/listings/:id/approve", adminHandler.HandleApprove)
	adminGroup.POST("/listings/:id/reject", adminHandler.HandleReject)

	// Environment Configuration
	env := os.Getenv("AGBALUMO_ENV")

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

	// Server Config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// In production (Fly.io), TLS is handled by the proxy. We just listen on PORT.
	// In dev, we might want TLS if certificates exist, OR just HTTP.
	if env == "production" {
		log.Printf("Starting Server in PRODUCTION mode on :%s", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Fatal(err)
		}
	} else {
		// Development Mode
		certFile := "certs/cert.pem"
		keyFile := "certs/key.pem"

		if _, err := os.Stat(certFile); err == nil {
			if _, err := os.Stat(keyFile); err == nil {
				log.Println("Starting Secure Server on :8443 (HTTPS)")
				if err := e.StartTLS(":8443", certFile, keyFile); err != nil {
					e.Logger.Fatal(err)
				}
				return
			}
		}

		log.Printf("Starting Server in DEV mode on :%s (HTTP)", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
