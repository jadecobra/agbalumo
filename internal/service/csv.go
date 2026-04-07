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

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	headerMap, err := s.validateHeaders(headers)
	if err != nil {
		return nil, err
	}

	return s.processRecords(ctx, reader, headerMap, repo), nil
}

func (s *CSVService) validateHeaders(headers []string) (map[string]int, error) {
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"title", "type", "description"}
	for _, req := range required {
		if _, ok := headerMap[req]; !ok {
			return nil, fmt.Errorf("missing required header: %s", req)
		}
	}
	return headerMap, nil
}

func (s *CSVService) processRecords(ctx context.Context, reader *csv.Reader, headerMap map[string]int, repo domain.ListingStore) *domain.BulkUploadResult {
	result := &domain.BulkUploadResult{}
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		result.TotalProcessed++

		if err != nil {
			s.recordFailure(result, lineNum, fmt.Errorf("failed to read row: %v", err))
			continue
		}

		if err := s.processRow(ctx, record, headerMap, repo, result, lineNum); err != nil {
			s.recordFailure(result, lineNum, err)
		} else {
			result.SuccessCount++
		}
	}
	return result
}

func (s *CSVService) processRow(ctx context.Context, record []string, headerMap map[string]int, repo domain.ListingStore, result *domain.BulkUploadResult, lineNum int) error {
	listing, err := s.parseRow(record, headerMap)
	if err != nil {
		return err
	}

	isDup, err := s.isDuplicate(ctx, repo, listing)
	if err != nil {
		return fmt.Errorf("duplicate check failed: %v", err)
	}
	if isDup {
		return fmt.Errorf("duplicate listing detected (title and >2 fields match)")
	}

	listing.OwnerID = ""
	listing.IsActive = true
	listing.Status = domain.ListingStatusApproved

	if err := repo.Save(ctx, *listing); err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	return nil
}

func (s *CSVService) recordFailure(result *domain.BulkUploadResult, lineNum int, err error) {
	result.FailureCount++
	result.Errors = append(result.Errors, fmt.Sprintf("Line %d: %v", lineNum, err))
}

func (s *CSVService) isDuplicate(ctx context.Context, repo domain.ListingStore, listing *domain.Listing) (bool, error) {
	existingListings, err := repo.FindByTitle(ctx, listing.Title)
	if err != nil {
		return false, err
	}

	for _, ex := range existingListings {
		if s.matchFieldsCount(ex, listing) > 2 {
			return true, nil
		}
	}
	return false, nil
}

func (s *CSVService) matchFieldsCount(ex domain.Listing, listing *domain.Listing) int {
	matches := 0
	check := func(s1, s2 string) {
		if s1 == s2 && s1 != "" {
			matches++
		}
	}

	if ex.Type == listing.Type {
		matches++
	}
	check(ex.Description, listing.Description)
	check(ex.OwnerOrigin, listing.OwnerOrigin)
	check(ex.ContactEmail, listing.ContactEmail)
	check(ex.ContactPhone, listing.ContactPhone)
	check(ex.ContactWhatsApp, listing.ContactWhatsApp)
	check(ex.Address, listing.Address)
	check(ex.City, listing.City)
	return matches
}

// GenerateCSV converts a slice of Listings into a CSV stream.
func (s *CSVService) GenerateCSV(ctx context.Context, listings []domain.Listing) (io.Reader, error) {
	pr, pw := io.Pipe()

	go func() {
		defer func() {
			_ = pw.Close()
		}()
		writer := csv.NewWriter(pw)
		defer writer.Flush()

		headers := []string{
			"ID", "Title", "Type", "Description", "City", "Address",
			"Origin", "Email", "Phone", "WhatsApp",
			"WebsiteURL", "CreatedAt", "Status", "IsActive", "Featured",
			"Company", "PayRange", "Skills", "JobApplyURL", "JobStartDate",
			"EventStart", "EventEnd", "Deadline",
		}
		if err := writer.Write(headers); err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		for _, l := range listings {
			row := s.listingToCSVRow(l)
			if err := writer.Write(row); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr, nil
}

func (s *CSVService) listingToCSVRow(l domain.Listing) []string {
	return []string{
		l.ID, l.Title, string(l.Type), l.Description, l.City, l.Address,
		l.OwnerOrigin, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp,
		l.WebsiteURL, l.CreatedAt.Format(time.RFC3339), string(l.Status),
		fmt.Sprintf("%v", l.IsActive), fmt.Sprintf("%v", l.Featured),
		l.Company, l.PayRange, l.Skills, l.JobApplyURL,
		l.JobStartDate.Format(time.RFC3339), l.EventStart.Format(time.RFC3339),
		l.EventEnd.Format(time.RFC3339), l.Deadline.Format(time.RFC3339),
	}
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
	if titleCas == "" {
		return domain.Business
	}
	return domain.Category(titleCas)
}

func validateParsedRow(title, desc, email, phone, whatsapp, website string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if desc == "" {
		return fmt.Errorf("description is required")
	}
	if email == "" && phone == "" && whatsapp == "" && website == "" {
		return fmt.Errorf("at least one contact method (email, phone, whatsapp, or website) is required")
	}
	return nil
}

func resolveCity(s *CSVService, city, address string) string {
	if city != "" || address == "" || s.Geocoding == nil {
		return city
	}
	if foundCity, err := s.Geocoding.GetCity(context.Background(), address); err == nil && foundCity != "" {
		return foundCity
	}
	return city
}

func (s *CSVService) parseRow(record []string, headerMap map[string]int) (*domain.Listing, error) {
	get := func(col string) string {
		if idx, ok := headerMap[col]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	title, desc := get("title"), get("description")
	email, phone, whatsapp, website := get("email"), get("phone"), get("whatsapp"), get("website")

	if err := validateParsedRow(title, desc, email, phone, whatsapp, website); err != nil {
		return nil, err
	}

	origin := get("origin")
	if origin == "" {
		origin = "Nigeria"
	}

	city := resolveCity(s, get("city"), get("address"))

	return &domain.Listing{
		ID: uuid.New().String(), Title: title, Type: parseCategory(get("type")),
		Description: desc, OwnerOrigin: origin, ContactEmail: email, WebsiteURL: website,
		ContactPhone: phone, ContactWhatsApp: whatsapp, Address: get("address"), City: city,
		HoursOfOperation: get("hours"), CreatedAt: time.Now(),
	}, nil
}
