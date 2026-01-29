package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/moderator"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer implements echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database Initialization
	repo, err := sqlite.NewSQLiteRepository("agbalumo.db")
	if err != nil {
		log.Fatal(err)
	}

	// Seed Data (if empty)
	ctx := context.Background()
	existing, _ := repo.FindAll(ctx, "")
	if len(existing) == 0 {
		seedData(ctx, repo)
	}

	// Static files (CSS, JS, Images)
	e.Static("/static", "ui/static")

	// Template Renderer
	// Parse both templates and partials
	tmpl := template.Must(template.New("").Parse("{{define \"base\"}}{{end}}")) // Base to allow appending
	// Note: ParseGlob might error if no files match, so be careful.
	// For simplicity, let's parse specific globs.
	template.Must(tmpl.ParseGlob("ui/templates/*.html"))
	template.Must(tmpl.ParseGlob("ui/templates/partials/*.html"))

	renderer := &TemplateRenderer{
		templates: tmpl,
	}
	e.Renderer = renderer

	// Routes
	e.GET("/", func(c echo.Context) error {
		listings, err := repo.FindAll(c.Request().Context(), "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Listings": listings,
		})
	})

	e.GET("/listings/fragment", func(c echo.Context) error {
		filterType := c.QueryParam("type")
		listings, err := repo.FindAll(c.Request().Context(), filterType)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Render just the fragment loop
		// We need to loop over listings and render "listing_card.html" for each.
		// Since our "listing_card" partial expects a single listing, we iterate here.
		// Alternatively, pass the slice to a new "listing_list.html" partial.
		// For simplicity/performance: use a dedicated partial for the list or iterate in code.
		// Let's create a partial for the list to keep templates clean, OR iterate here writing to response.
		// Echo's Render doesn't support multiple calls easily for streaming.
		// Best practice: Create a partial "partials/listing_list.html" that iterates.
		return c.Render(http.StatusOK, "listing_list.html", map[string]interface{}{
			"Listings": listings,
		})
	})

	// Endpoint for detail modal
	e.GET("/listings/:id", func(c echo.Context) error {
		id := c.Param("id")
		listing, err := repo.FindByID(c.Request().Context(), id)
		if err != nil {
			return c.String(http.StatusNotFound, "Listing not found")
		}

		return c.Render(http.StatusOK, "modal_detail.html", listing)
	})

	e.POST("/listings", func(c echo.Context) error {
		// New: Handle Multipart Form manually to support file upload + strict logic
		var l domain.Listing
		l.ID = uuid.New().String()
		now := time.Now()
		l.CreatedAt = now
		l.IsActive = true

		// Bind Form Values
		l.Title = c.FormValue("title")
		l.Type = domain.Category(c.FormValue("type"))
		l.OwnerOrigin = c.FormValue("owner_origin")
		l.Description = c.FormValue("description")
		l.Neighborhood = c.FormValue("neighborhood") // Add if in UI, or keep default?
		l.ContactEmail = c.FormValue("contact_email")
		l.ContactPhone = c.FormValue("contact_phone")
		l.ContactWhatsApp = c.FormValue("contact_whatsapp")
		l.WebsiteURL = c.FormValue("website_url")

		// Handle Image Upload
		file, err := c.FormFile("image")
		if err == nil {
			src, err := file.Open()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Image Upload Error")
			}
			defer src.Close()

			// Simple file save (Production would use object storage)
			// Ensure directory exists: ui/static/uploads
			// Create unique name
			ext := ".jpg" // Default or parse file.Filename (security risk to trust directly)
			// MVP: Simplistic extension check or just use uuid
			filename := l.ID + ext
			dstPath := "ui/static/uploads/" + filename
			dst, err := os.Create(dstPath)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Image Save Error")
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return c.String(http.StatusInternalServerError, "Image Copy Error")
			}
			l.ImageURL = "/static/uploads/" + filename
		}

		// Handle Deadline Logic
		if l.Type == domain.Request {
			dateStr := c.FormValue("deadline_date")
			if dateStr != "" {
				parsedTime, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return c.String(http.StatusBadRequest, "Invalid Date Format")
				}
				// Set time to end of day? Or same time as now?
				// content of Parse is 00:00 UTC usually.
				l.Deadline = parsedTime

				// Validate < 90 days from Now
				limit := now.Add(90 * 24 * time.Hour)
				if l.Deadline.After(limit) {
					return c.String(http.StatusBadRequest, "Validation Error: Deadline cannot be more than 90 days away")
				}

				// Validate > Now (not in past) - Optional but good UX
				if l.Deadline.Before(now.Add(-24 * time.Hour)) { // Allow today
					return c.String(http.StatusBadRequest, "Validation Error: Deadline cannot be in the past")
				}

			} else {
				l.Deadline = now.Add(90*24*time.Hour - time.Minute) // Default < 90 days
			}
		}

		// Validation
		if err := l.Validate(); err != nil {
			return c.String(http.StatusBadRequest, "Validation Error: "+err.Error())
		}

		// Moderation
		mod, err := moderator.NewGeminiModerator(c.Request().Context())
		if err != nil {
			log.Println("Moderator Init Error:", err)
		} else {
			if err := mod.CheckListing(c.Request().Context(), l); err != nil {
				return c.String(http.StatusBadRequest, "Content Moderation Failed: "+err.Error())
			}
		}

		// Save
		if err := repo.Save(c.Request().Context(), l); err != nil {
			return c.String(http.StatusInternalServerError, "Database Error: "+err.Error())
		}

		// Return just the new card fragment
		// Note: The template expects a struct that has .ImageURL populated correctly.
		return c.Render(http.StatusOK, "listing_card.html", l)
	})

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
