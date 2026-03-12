package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

type CSVService struct {
	Geocoding domain.GeocodingService
}

func NewCSVService() *CSVService {
	return &CSVService{}
}

// ParseAndImport reads a CSV stream and converts rows into Listings, saving them to the repo.
func (s *CSVService) ParseAndImport(ctx context.Context, r io.Reader, repo domain.ListingStore) (*domain.BulkUploadResult, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Read Header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// map header name to index
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Validate required headers
	required := []string{"title", "type", "description"}
	for _, req := range required {
		if _, ok := headerMap[req]; !ok {
			return nil, fmt.Errorf("missing required header: %s", req)
		}
	}

	result := &domain.BulkUploadResult{}
	lineNum := 1 // Header is line 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		result.TotalProcessed++

		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Failed to read row: %v", lineNum, err))
			continue
		}

		// Parse Row
		listing, err := s.parseRow(record, headerMap)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: %v", lineNum, err))
			continue
		}

		// Check for duplicate
		existingListings, err := repo.FindByTitle(ctx, listing.Title)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Database error checking duplicates: %v", lineNum, err))
			continue
		}

		isDuplicate := false
		for _, ex := range existingListings {
			matchCount := 0
			if ex.Type == listing.Type {
				matchCount++
			}
			if ex.Description == listing.Description && listing.Description != "" {
				matchCount++
			}
			if ex.OwnerOrigin == listing.OwnerOrigin && listing.OwnerOrigin != "" {
				matchCount++
			}
			if ex.ContactEmail == listing.ContactEmail && listing.ContactEmail != "" {
				matchCount++
			}
			if ex.ContactPhone == listing.ContactPhone && listing.ContactPhone != "" {
				matchCount++
			}
			if ex.ContactWhatsApp == listing.ContactWhatsApp && listing.ContactWhatsApp != "" {
				matchCount++
			}
			if ex.Address == listing.Address && listing.Address != "" {
				matchCount++
			}
			if ex.City == listing.City && listing.City != "" {
				matchCount++
			}

			// If title matches AND more than 2 other fields match, consider it a duplicate
			if matchCount > 2 {
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Duplicate listing detected (title and >2 fields match)", lineNum))
			continue
		}

		// Set System Defaults for Uploads
		listing.OwnerID = "" // Unowned (Seeded)
		listing.IsActive = true
		listing.Status = domain.ListingStatusApproved

		// Save to Repo
		if err := repo.Save(ctx, *listing); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Database error: %v", lineNum, err))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}

// GenerateCSV converts a slice of Listings into a CSV stream.
func (s *CSVService) GenerateCSV(ctx context.Context, listings []domain.Listing) (io.Reader, error) {
	pr, pw := io.Pipe()

	go func() {
		defer func() {
			if err := pw.Close(); err != nil {
				// We can't do much here except log if we had a logger,
				// but at least we check it to satisfy the linter.
				_ = err
			}
		}()
		writer := csv.NewWriter(pw)
		defer writer.Flush()

		// Write Header
		headers := []string{
			"ID", "Title", "Type", "Description", "City", "Address",
			"Origin", "Email", "Phone", "WhatsApp",
			"WebsiteURL", "CreatedAt", "Status", "IsActive", "Featured",
			"Company", "PayRange", "Skills", "JobApplyURL", "JobStartDate",
			"EventStart", "EventEnd", "Deadline",
		}
		if err := writer.Write(headers); err != nil {
			pw.CloseWithError(err)
			return
		}

		// Write Rows
		for _, l := range listings {
			row := []string{
				l.ID,
				l.Title,
				string(l.Type),
				l.Description,
				l.City,
				l.Address,
				l.OwnerOrigin,
				l.ContactEmail,
				l.ContactPhone,
				l.ContactWhatsApp,
				l.WebsiteURL,
				l.CreatedAt.Format(time.RFC3339),
				string(l.Status),
				fmt.Sprintf("%v", l.IsActive),
				fmt.Sprintf("%v", l.Featured),
				l.Company,
				l.PayRange,
				l.Skills,
				l.JobApplyURL,
				l.JobStartDate.Format(time.RFC3339),
				l.EventStart.Format(time.RFC3339),
				l.EventEnd.Format(time.RFC3339),
				l.Deadline.Format(time.RFC3339),
			}
			if err := writer.Write(row); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr, nil
}

func parseCategory(typeStr string) domain.Category {
	cat := domain.Category(typeStr)
	switch cat {
	case domain.Business, domain.Service, domain.Product, domain.Event, domain.Job, domain.Request, domain.Food:
		return cat
	}

	r := []rune(strings.ToLower(typeStr))
	if len(r) > 0 {
		r[0] = unicode.ToUpper(r[0])
	}
	titleCas := string(r)
	cat = domain.Category(titleCas)
	switch cat {
	case domain.Business, domain.Service, domain.Product, domain.Event, domain.Job, domain.Request, domain.Food:
		return cat
	}

	return domain.Business
}

func (s *CSVService) parseRow(record []string, headerMap map[string]int) (*domain.Listing, error) {
	get := func(col string) string {
		if idx, ok := headerMap[col]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	title := get("title")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	cat := parseCategory(get("type"))

	desc := get("description")
	if desc == "" {
		return nil, fmt.Errorf("description is required")
	}

	origin := get("origin")
	if origin == "" {
		origin = "Nigeria"
	}

	email := get("email")
	phone := get("phone")
	whatsapp := get("whatsapp")

	if email == "" && phone == "" && whatsapp == "" {
		return nil, fmt.Errorf("at least one contact method (email, phone, or whatsapp) is required")
	}

	address := get("address")
	city := get("city")

	// If city is missing but address exists, try to geocode if service is available
	if city == "" && address != "" && s.Geocoding != nil {
		// Note: Using background context here as parseRow doesn't take context
		// This is a trade-off for the current structure.
		// A better refactor would pass ctx to parseRow, but let's stick to minimal changes.
		if foundCity, err := s.Geocoding.GetCity(context.Background(), address); err == nil && foundCity != "" {
			city = foundCity
		}
	}

	return &domain.Listing{
		ID:               uuid.New().String(),
		Title:            title,
		Type:             cat,
		Description:      desc,
		OwnerOrigin:      origin,
		ContactEmail:     email,
		WebsiteURL:       get("website"),
		ContactPhone:     phone,
		ContactWhatsApp:  whatsapp,
		Address:          address,
		City:             city,
		HoursOfOperation: get("hours"),
		CreatedAt:        time.Now(),
	}, nil
}
