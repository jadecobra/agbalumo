package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

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
			"EventStart", "EventEnd", "Deadline", "EnrichmentAttemptedAt",
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
	attemptedAtStr := ""
	if l.EnrichmentAttemptedAt != nil {
		attemptedAtStr = l.EnrichmentAttemptedAt.Format(time.RFC3339)
	}

	return []string{
		l.ID, l.Title, string(l.Type), l.Description, l.City, l.Address,
		l.OwnerOrigin, l.ContactEmail, l.ContactPhone, l.ContactWhatsApp,
		l.WebsiteURL, l.CreatedAt.Format(time.RFC3339), string(l.Status),
		fmt.Sprintf("%v", l.IsActive), fmt.Sprintf("%v", l.Featured),
		l.Company, l.PayRange, l.Skills, l.JobApplyURL,
		l.JobStartDate.Format(time.RFC3339), l.EventStart.Format(time.RFC3339),
		l.EventEnd.Format(time.RFC3339), l.Deadline.Format(time.RFC3339),
		attemptedAtStr,
	}
}
