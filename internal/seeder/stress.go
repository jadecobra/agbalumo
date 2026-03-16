package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// GenerateStressData creates 'count' randomized listings spread across categories and origins.
func GenerateStressData(ctx context.Context, repo domain.ListingStore, count int) {
	slog.Info("Starting stress data generation", "count", count)
	
	categories := []domain.Category{
		domain.Business, domain.Service, domain.Product, domain.Job, 
		domain.Request, domain.Food, domain.Event,
	}

	origins := make([]string, 0, len(domain.ValidOrigins))
	for k := range domain.ValidOrigins {
		origins = append(origins, k)
	}

	cities := []string{"Dallas", "Fort Worth", "Arlington", "Plano", "Irving", "Garland", "Frisco", "McKinney", "Grand Prairie", "Mesquite"}
	streets := []string{"Main St", "Elm St", "Commerce St", "Belt Line Rd", "Legacy Dr", "MacArthur Blvd", "Central Expy", "Preston Rd"}

	successCount := 0
	batchSize := 1000

	for i := 0; i < count; i++ {
		category := categories[rand.IntN(len(categories))]
		origin := origins[rand.IntN(len(origins))]
		city := cities[rand.IntN(len(cities))]
		address := fmt.Sprintf("%d %s, %s, TX", rand.IntN(9000)+100, streets[rand.IntN(len(streets))], city)

		title := generateRandomString(rand.IntN(20) + 10)
		description := generateRandomString(rand.IntN(500) + 50)

		l := domain.Listing{
			ID:           uuid.New().String(),
			Type:         category,
			OwnerOrigin:  origin,
			City:         city,
			Title:        "Stress " + title,
			Description:  description,
			CreatedAt:    time.Now().Add(-time.Duration(rand.IntN(365*24)) * time.Hour), // Random time in past year
			IsActive:     true,
			Status:       domain.ListingStatusApproved,
			ContactEmail: fmt.Sprintf("stress%d@example.com", i),
		}

		// Conditional Logic
		switch category {
		case domain.Business, domain.Food:
			l.Address = address
			if rand.IntN(2) == 0 {
				l.HoursOfOperation = "9 AM - 5 PM"
			}
		case domain.Job:
			l.Company = "Company " + generateRandomString(10)
			l.Skills = "Go, React, SQL"
			l.PayRange = "$50k - $100k"
			l.JobStartDate = time.Now().Add(time.Duration(rand.IntN(30)*24) * time.Hour)
			l.JobApplyURL = "https://example.com/apply"
			l.Address = address // often needed if city isn't enough, but city is set
		case domain.Event:
			l.EventStart = time.Now().Add(time.Duration(rand.IntN(30)*24) * time.Hour)
			l.EventEnd = l.EventStart.Add(2 * time.Hour)
		case domain.Request:
			l.Deadline = time.Now().Add(time.Duration(rand.IntN(60)*24) * time.Hour)
		}

		// Save the listing
		if err := repo.Save(ctx, l); err != nil {
			slog.Error("Failed to save stress listing", "id", l.ID, "error", err)
		} else {
			successCount++
		}

		if (i+1)%batchSize == 0 {
			slog.Info("Progress", "inserted", i+1, "total", count)
		}
	}

	slog.Info("Stress data generation complete", "success", successCount, "attempted", count)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "

func generateRandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}
	return sb.String()
}
