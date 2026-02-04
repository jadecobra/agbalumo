package main

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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
	e.Use(customMiddleware.RateLimitMiddleware)

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
	repo, err := sqlite.NewSQLiteRepository("agbalumo.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// ...
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
	e.GET("/", listingHandler.HandleHome)
	// ...

	// Seed Data (if empty)
	ctx := context.Background()
	var existing []domain.Listing
	existing, _ = repo.FindAll(ctx, "", "", true)
	if len(existing) == 0 {
		seedData(ctx, repo)
	}

	// Static files (CSS, JS, Images)
	e.Static("/static", "ui/static")

	// Template Renderer
	// Parse both templates and partials
	tmpl := template.New("").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
		"add": func(i, j int) int { return i + j },
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	})
	// Base to allow appending
	template.Must(tmpl.Parse("{{define \"base\"}}{{end}}"))
	// Note: ParseGlob might error if no files match, so be careful.
	// For simplicity, let's parse specific globs.
	template.Must(tmpl.ParseGlob("ui/templates/*.html"))
	template.Must(tmpl.ParseGlob("ui/templates/partials/*.html"))

	renderer := &TemplateRenderer{
		templates: tmpl,
	}
	e.Renderer = renderer

	// Routes
	e.GET("/", listingHandler.HandleHome)
	e.GET("/listings/fragment", listingHandler.HandleFragment)
	e.GET("/listings/:id", listingHandler.HandleDetail)
	e.POST("/listings", listingHandler.HandleCreate)
	// Edit Routes
	e.Use(authHandler.OptionalAuth) // Ensure user is available for check
	e.GET("/listings/:id/edit", listingHandler.HandleEdit, authHandler.RequireAuth)
	e.PUT("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth)
	e.POST("/listings/:id", listingHandler.HandleUpdate, authHandler.RequireAuth) // Fallback support for POST? HTMX sends PUT if requested.

	// Admin Routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(adminHandler.AuthMiddleware)

	e.GET("/admin/login", adminHandler.HandleLoginView)
	e.POST("/admin/login", adminHandler.HandleLoginAction)

	adminGroup.GET("", adminHandler.HandleDashboard)
	adminGroup.DELETE("/listings/:id", adminHandler.HandleDelete)

	// Start server
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}

func seedData(ctx context.Context, repo domain.ListingRepository) {
	listings := []domain.Listing{
		{
			ID:              "1",
			Title:           "Lagos Spot Kitchen",
			OwnerOrigin:     "Nigeria",
			Type:            domain.Business,
			Neighborhood:    "North Dallas",
			Description:     "Authentic Naija jollof and suya spots in the heart of Dallas.",
			ContactEmail:    "info@lagosspot.com",
			ContactWhatsApp: "+12145550100",
			CreatedAt:       time.Now(),
			IsActive:        true,
		},
		{
			ID:              "2",
			Title:           "Kofi's Legal Aid",
			OwnerOrigin:     "Ghana",
			Type:            domain.Service,
			Neighborhood:    "Downtown",
			Description:     "Immigration and small business legal consultation.",
			ContactEmail:    "kofi@legalaid.com",
			ContactWhatsApp: "+12145550200",
			CreatedAt:       time.Now(),
			IsActive:        true,
		},
	}
	for _, l := range listings {
		if err := repo.Save(ctx, l); err != nil {
			log.Printf("Failed to seed listing: %v", err)
		}
	}
}
