package seeder

import (
	"fmt"
	"github.com/jadecobra/agbalumo/internal/domain"
	"math/rand/v2"
	"time"
)

var origins []string

func init() {
	for k := range domain.ValidOrigins {
		origins = append(origins, k)
	}
}

// GenerateStressListings generates a specified number of randomized, valid domain.Listing entities.
func GenerateStressListings(count int) []domain.Listing {
	var listings []domain.Listing

	categories := []domain.Category{
		domain.Business,
		domain.Service,
		domain.Product,
		domain.Job,
		domain.Request,
		domain.Food,
		domain.Event,
	}

	cities := []string{"Lagos", "Abuja", "Accra", "Dakar", "Nairobi", "Johannesburg", "London", "New York", "Toronto"}

	now := time.Now()

	for i := 0; i < count; i++ {
		cat := categories[rand.IntN(len(categories))]
		origin := origins[rand.IntN(len(origins))]
		city := cities[rand.IntN(len(cities))]

		l := domain.Listing{
			ID:              fmt.Sprintf("lst_%d_%d", now.UnixNano(), i),
			OwnerID:         fmt.Sprintf("usr_%d", rand.IntN(1000)),
			OwnerOrigin:     origin,
			Type:            cat,
			Title:           fmt.Sprintf("Stress Test Listing %d", i),
			Description:     "This is an automated listing generated for stress testing purposes.",
			City:            city,
			ContactEmail:    fmt.Sprintf("test%d@example.com", i),
			IsActive:        true,
			Status:          domain.ListingStatusApproved,
			CreatedAt:       now.Add(-time.Duration(rand.IntN(30*24)) * time.Hour), // Randomly created in last 30 days
		}

		switch cat {
		case domain.Business, domain.Food:
			l.Address = fmt.Sprintf("%d Test Ave, %s", rand.IntN(9999), city)
			l.HoursOfOperation = "9 AM - 5 PM"
		case domain.Service:
			l.HoursOfOperation = "24/7"
		case domain.Job:
			l.Company = "Acme Corp"
			l.Skills = "Go, React, SQL"
			l.PayRange = "$50k - $100k"
			l.JobStartDate = now.Add(time.Duration(rand.IntN(30*24)) * time.Hour) // Starts in next 30 days
			l.JobApplyURL = "https://example.com/apply"
			l.Address = "123 Job St"
		case domain.Event:
			l.EventStart = now.Add(time.Duration(rand.IntN(5*24)) * time.Hour)
			l.EventEnd = l.EventStart.Add(4 * time.Hour)
		case domain.Request:
			l.Deadline = l.CreatedAt.Add(30 * 24 * time.Hour) // 30 days after creation
		}

		listings = append(listings, l)
	}

	return listings
}
