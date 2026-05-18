package seeder

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

var origins []string

func init() {
	for k := range domain.ValidOrigins {
		origins = append(origins, k)
	}
}

// GenerateStressListings generates a specified number of randomized, valid domain.Listing entities.
func GenerateStressListings(count int) []domain.Listing {
	if count <= 0 {
		return nil
	}

	listings := make([]domain.Listing, count)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	if numWorkers > count {
		numWorkers = count
	}

	chunkSize := count / numWorkers

	categories := []domain.Category{
		domain.Business,
		domain.Service,
		domain.Product,
		domain.Job,
		domain.Request,
		domain.Food,
		domain.Event,
	}

	cities := []string{"Dallas", "Arlington", "Plano", "Fort Worth", "Irving", "Garland", "Richardson", "McKinney", "Mesquite"}

	now := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := workerID * chunkSize
			end := start + chunkSize
			if workerID == numWorkers-1 {
				end = count
			}

			for i := start; i < end; i++ {
				listings[i] = generateSingleStressListing(i, workerID, now, categories, origins, cities)
			}
		}(w)
	}

	wg.Wait()
	return listings
}

func generateSingleStressListing(i, workerID int, now time.Time, categories []domain.Category, origins, cities []string) domain.Listing {
	// #nosec G404 - weak random is acceptable for non-security-critical stress testing
	cat := categories[rand.IntN(len(categories))]
	// #nosec G404 - weak random is acceptable for non-security-critical stress testing
	origin := origins[rand.IntN(len(origins))]
	// #nosec G404 - weak random is acceptable for non-security-critical stress testing
	city := cities[rand.IntN(len(cities))]

	var lat, lng float64
	switch city {
	case "Dallas":
		lat, lng = 32.7767, -96.7970
	case "Arlington":
		lat, lng = 32.7357, -97.1081
	case "Plano":
		lat, lng = 33.0198, -96.6989
	case "Fort Worth":
		lat, lng = 32.7555, -97.3308
	case "Irving":
		lat, lng = 32.8140, -96.9489
	case "Garland":
		lat, lng = 32.9126, -96.6389
	case "Richardson":
		lat, lng = 32.9483, -96.7299
	case "McKinney":
		lat, lng = 33.1972, -96.6398
	case "Mesquite":
		lat, lng = 32.7668, -96.5992
	default:
		lat, lng = 32.7767, -96.7970
	}
	// Add small random offset (roughly up to ~3.5 miles spread)
	// #nosec G404 - weak random is acceptable for non-security-critical stress testing
	lat += (rand.Float64() - 0.5) * 0.1
	// #nosec G404 - weak random is acceptable for non-security-critical stress testing
	lng += (rand.Float64() - 0.5) * 0.1

	var title, desc string
	if cat == domain.Food {
		dishes := []string{"Suya", "Jollof Rice", "Egusi & Pounded Yam", "Pepper Soup", "Mama Put Specialty", "Naija Feast"}
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		dish := dishes[rand.IntN(len(dishes))]
		title = fmt.Sprintf("Nigerian %s %d", dish, i)
		desc = fmt.Sprintf("Authentic Nigerian food cooked with passion. Try our legendary %s. Best in %s!", dish, city)
	} else {
		title = fmt.Sprintf("Stress Test Listing %d", i)
		desc = "This is an automated listing generated for stress testing purposes."
	}

	l := domain.Listing{
		ID: fmt.Sprintf("lst_%d_%d_%d", now.UnixNano(), workerID, i),
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		OwnerID:      fmt.Sprintf("usr_%d", rand.IntN(1000)),
		OwnerOrigin:  origin,
		Type:         cat,
		Title:        title,
		Description:  desc,
		City:         city,
		State:        "TX",
		Country:      "USA",
		ContactEmail: fmt.Sprintf("test%d@example.com", i),
		IsActive:     true,
		Status:       domain.ListingStatusApproved,
		Latitude:     lat,
		Longitude:    lng,
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		CreatedAt: now.Add(-time.Duration(rand.IntN(30*24)) * time.Hour), // Randomly created in last 30 days
	}

	switch cat {
	case domain.Business, domain.Food:
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		l.Address = fmt.Sprintf("%d Test Ave, %s, TX", rand.IntN(9999), city)
		l.HoursOfOperation = "9 AM - 5 PM"
	case domain.Service:
		l.HoursOfOperation = "24/7"
	case domain.Job:
		l.Company = "Acme Corp"
		l.Skills = "Go, React, SQL"
		l.PayRange = "$50k - $100k"
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		l.JobStartDate = now.Add(time.Duration(rand.IntN(30*24)+1) * time.Hour) // Starts in next 30 days
		l.JobApplyURL = "https://example.com/apply"
		l.Address = "123 Job St"
	case domain.Event:
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		l.EventStart = now.Add(time.Duration(rand.IntN(5*24)+1) * time.Hour)
		l.EventEnd = l.EventStart.Add(4 * time.Hour)
	case domain.Request:
		// #nosec G404 - weak random is acceptable for non-security-critical stress testing
		l.Deadline = now.Add(time.Duration(rand.IntN(30*24)+1) * time.Hour)
	}

	return l
}
