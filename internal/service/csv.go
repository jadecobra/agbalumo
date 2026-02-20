package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
)

type CSVService struct{}

func NewCSVService() *CSVService {
	return &CSVService{}
}

// ParseAndImport reads a CSV stream and converts rows into Listings, saving them to the repo.
func (s *CSVService) ParseAndImport(ctx context.Context, r io.Reader, repo domain.ListingSaver) (*domain.BulkUploadResult, error) {
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

func parseCategory(typeStr string) (domain.Category, error) {
	cat := domain.Category(typeStr)
	switch cat {
	case domain.Business, domain.Service, domain.Product, domain.Event, domain.Job, domain.Request, domain.Food:
		return cat, nil
	}

	titleCas := strings.Title(strings.ToLower(typeStr))
	cat = domain.Category(titleCas)
	switch cat {
	case domain.Business, domain.Service, domain.Product, domain.Event, domain.Job, domain.Request, domain.Food:
		return cat, nil
	}

	return domain.Business, nil
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

	cat, err := parseCategory(get("type"))
	if err != nil {
		return nil, err
	}

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
		Address:          get("address"),
		HoursOfOperation: get("hours"),
		CreatedAt:        time.Now(),
	}, nil
}
